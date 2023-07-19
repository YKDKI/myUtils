package myUtils

import (
	"bytes"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var (
	client *resty.Client
)

func InitHttpClient(conf *viper.Viper, log *Logger) {
	client = resty.New()
	client.SetDebug(conf.GetString("env") != "prod")
	client.SetTimeout(time.Duration(conf.GetInt("http.timeout")) * time.Second)

	if conf.GetBool("http.retry") {
		client.SetRetryCount(conf.GetInt("http.retry.count")).
			SetRetryWaitTime(time.Duration(conf.GetInt("http.retry.waittime")) * time.Microsecond).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					return r.StatusCode() == http.StatusTooManyRequests
				},
			).
			AddRetryCondition(
				func(r *resty.Response, err error) bool {
					isNetError, _, _ := isCaredNetError(err)
					return isNetError
				},
			)
	}

	log.Info("http client init done")
}

func SendBotMessage(url, msg string) (bool, error) {
	msgStruct := struct {
		Msg string
	}{
		Msg: msg,
	}
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(msgStruct).Post(url)

	if err != nil {
		return false, err
	}
	return true, nil
}

func UploadFile(url, fileName string, content []byte) (bool, error) {
	_, err := client.R().
		SetFileReader("file", fileName, bytes.NewReader(content)).
		Post(url)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Get(url string, msg url.Values, headers map[string]string) ([]byte, error) {
	r := client.R().SetQueryParamsFromValues(msg)
	if headers != nil {
		r.SetHeaders(headers)
	}
	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func Post(url, msg string, headers map[string]string) ([]byte, error) {
	r := client.R().SetBody([]byte(msg))
	if headers != nil {
		r.SetHeaders(headers)
	}
	resp, err := r.Post(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func isCaredNetError(err error) (netError bool, timeout bool, details string) {
	if err == nil {
		return false, false, ""
	}

	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			// 连接超时
			return true, true, "timeout"
		}
		// 其他网络错误
		return true, false, netErr.Error()
	}

	// *net.OpError包含更详细的错误信息
	opErr, ok := err.(*net.OpError)
	if !ok {
		return false, false, ""
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		// DNS解析错误
		return true, false, "dns error: " + t.Error()
	case *os.SyscallError:
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNRESET:
				// 对方重置连接
				return true, false, "connection reset by peer"
			case syscall.ECONNREFUSED:
				// 连接被拒绝
				return true, false, "connection refused"
			case syscall.ETIMEDOUT:
				// 连接超时
				return true, true, "timeout"
			}
		}
	}

	return false, false, ""
}
