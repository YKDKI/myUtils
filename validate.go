package myUtils

import "regexp"

func IsPhoneNumber(phone string) bool {
	last11 := phone
	if len(phone) > 11 {
		last11 = phone[len(phone)-11:]
	}
	valid, err := regexp.MatchString(`^1[3-9][0-9]{9}$`, last11)
	if err != nil {
		return false
	}
	return valid
}
