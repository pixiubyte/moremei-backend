package utils

import (
	"html"
	"regexp"
	"strings"
)

func RemoveHtmlTags(input string) string {
	input = html.UnescapeString(input)
	//将HTML标签全转换成小写
	input = regexp.MustCompile(`<[\S\s]+?>`).ReplaceAllStringFunc(input, strings.ToLower)
	//去除STYLE
	input = regexp.MustCompile(`<style[\S\s]+?</style>`).ReplaceAllString(input, "")
	//去除SCRIPT
	input = regexp.MustCompile(`<script[\S\s]+?</script>`).ReplaceAllString(input, "")
	//去除所有尖括号内的HTML代码
	input = regexp.MustCompile(`<[\S\s]+?>`).ReplaceAllString(input, "")
	//去除连续的换行符
	input = regexp.MustCompile(`\s{2,}`).ReplaceAllString(input, "\n")
	return strings.TrimSpace(strings.Trim(input, "\n"))
}
