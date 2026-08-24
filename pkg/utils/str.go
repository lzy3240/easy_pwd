package utils

import (
	"encoding/base64"
	"net/url"
	"regexp"
	"strings"
)

// 检查字符串是否全部由星号组成, 如果是则返回true, 否则返回false
func IsAllAsterisk(s string) bool {
	if len(s) == 0 {
		return false // 空字符串不是"全部为星号"
	}

	for _, ch := range s {
		if ch != '*' {
			return false
		}
	}
	return true
}

// 判断字符串是否为base64编码, 如是则返回true, 否则返回false
func IsBase64(s string) bool {
	// 去除空白字符
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}

	// 尝试解码
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// IsValidUrlAny 更宽松的检查，允许没有 scheme
func IsValidUrlAny(str string) bool {
	if str == "" {
		return false
	}

	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	// 至少要有 Host 或 Path
	return u.Host != "" || u.Path != ""
}

// IsUrlRegex 使用正则表达式检查 URL 格式
func IsUrlRegex(str string) bool {
	if str == "" {
		return false
	}

	// 相对宽松的 URL 正则
	// 支持 http/https/ftp 等
	// 也支持 www.xxx.com 这种无 scheme 的写法
	pattern := `^(https?|ftp)://[^\s/$.?#].[^\s]*$|^www\.[^\s/$.?#].[^\s]*$|^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}(/.*)?$`

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}
