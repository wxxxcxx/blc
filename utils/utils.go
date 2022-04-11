package utils

import (
	"fmt"
	"runtime"
	"strings"
)

func Pathnamify(name string) string {
	return FileName(name, "", 255)
}

func Filenamify(name string, ext string) string {
	length := 255 - len(ext) - 1
	return FileName(name, "md", length)
}

func LimitLength(s string, length int) string {
	// 0 means unlimited
	if length == 0 {
		return s
	}

	const ELLIPSES = "..."
	str := []rune(s)
	if len(str) > length {
		return string(str[:length-len(ELLIPSES)]) + ELLIPSES
	}
	return s
}
func FileName(name, ext string, length int) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’")
	name = rep.Replace(name)
	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ")
		name = rep.Replace(name)
	}
	limitedName := LimitLength(name, length)
	if ext == "" {
		return limitedName
	}
	return fmt.Sprintf("%s.%s", limitedName, ext)
}
