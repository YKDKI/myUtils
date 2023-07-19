package myUtils

import (
	"io"
	"os"
	"strings"
)

func SplitDedupliFile(filePath, splitChat string) (map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(string((content)), "\r\n", "\n"), splitChat)

	set := make(map[string]bool)
	for _, line := range lines {
		if line != "" {
			set[line] = true
		}
	}
	return set, nil
}
