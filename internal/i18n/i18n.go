package i18n

import (
	"fmt"
	"os"
	"strings"
)

type Key string

var (
	current   map[Key]string
	isKorean  bool
)

func init() {
	lang := os.Getenv("LC_ALL")
	if lang == "" {
		lang = os.Getenv("LC_MESSAGES")
	}
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	if strings.HasPrefix(lang, "ko") {
		current = ko
		isKorean = true
	} else {
		current = en
		isKorean = false
	}
}

func T(key Key) string {
	if msg, ok := current[key]; ok {
		return msg
	}
	return string(key)
}

func Tf(key Key, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

func IsKorean() bool {
	return isKorean
}

// SetLang allows overriding for testing
func SetLang(langCode string) {
	if strings.HasPrefix(langCode, "ko") {
		current = ko
		isKorean = true
	} else {
		current = en
		isKorean = false
	}
}
