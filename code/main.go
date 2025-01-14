package main

import (
	"fmt"
	"regexp"
)

func parseFormat(str string) []string {
	//前三个位置分别存储 三个固定属性(msg,http_code,rpc_code)
	//后面的属性存储定义的format
	placeholders := make([]string, 3)
	paragraphs := make([]string, 0)

	//用正则表达式找到{}内的所有值
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatchIndex(str, -1)
	lastIndex := 0

	for _, match := range matches {
		// 提取占位符内的内容,以此得到三个属性的出现顺序
		placeholders = append(placeholders, str[match[2]:match[3]])

		// 提取占位符之前的段落格式,如果前面不存在前缀则不会添加
		if match[0] > lastIndex {
			paragraphs = append(paragraphs, str[lastIndex:match[0]])
		}
		// 更新最后索引位置
		lastIndex = match[1]
	}

	// 添加最后一个有值的段落(边界条件)
	if lastIndex < len(str) {
		paragraphs = append(paragraphs, str[lastIndex:])
	}
	return append(placeholders, paragraphs...)
}

func main() {
	s := "1f1wd{http}q:{rpc}:{msg}qwrqwt"
	w := parseFormat(s)
	fmt.Println(w)
}
