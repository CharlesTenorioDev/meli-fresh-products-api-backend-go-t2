package utils

import (
	"fmt"
	"os"
	"strings"
)

func LoadProperties(path string) error {
	bContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(bContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}

		keyAndValue := strings.Split(line, "=")
		if len(keyAndValue) != 2 {
			return ErrInvalidProperties
		}

		key := strings.TrimSpace(keyAndValue[0])
		value := strings.TrimSpace(keyAndValue[1])
		fmt.Println(key, value)
		os.Setenv(key, value)
	}

	return nil
}
