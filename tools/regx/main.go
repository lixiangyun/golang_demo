package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := `Hello 世界！123 Go.`

	// 查找连续的小写字母
	reg, _ := regexp.Compile(`[a-z]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["ello" "o"]

	// 查找连续的非小写字母
	reg, _ = regexp.Compile(`[^a-z]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["H" " 世界！123 G" "."]

	// 查找连续的单词字母
	reg, _ = regexp.Compile(`[\w]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "123" "Go"]

	// 查找连续的非单词字母、非空白字符
	reg, _ = regexp.Compile(`[^\w\s]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界！" "."]

	// 查找连续的大写字母
	reg, _ = regexp.Compile(`[[:upper:]]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["H" "G"]

	// 查找连续的非 ASCII 字符
	reg, _ = regexp.Compile(`[[:^ascii:]]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界！"]

	// 查找连续的标点符号
	reg, _ = regexp.Compile(`[\pP]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["！" "."]

	// 查找连续的非标点符号字符
	reg, _ = regexp.Compile(`[\PP]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界" "123 Go"]

	// 查找连续的汉字
	reg, _ = regexp.Compile(`[\p{Han}]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["世界"]

	// 查找连续的非汉字字符
	reg, _ = regexp.Compile(`[\P{Han}]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello " "！123 Go."]

	// 查找 Hello 或 Go
	reg, _ = regexp.Compile(`Hello|Go`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找行首以 H 开头，以空格结尾的字符串
	reg, _ = regexp.Compile(`^H.*\s`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界！123 "]

	// 查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	reg, _ = regexp.Compile(`(?U)^H.*\s`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello "]

	// 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	reg, _ = regexp.Compile(`(?i:^hello).*Go`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello 世界！123 Go"]

	// 查找 Go.
	reg, _ = regexp.Compile(`\QGo.\E`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Go."]

	// 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	reg, _ = regexp.Compile(`(?U)^.* `)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello "]

	// 查找以空格开头，到行尾结束，中间不包含空格字符串
	reg, _ = regexp.Compile(` [^ ]*$`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// [" Go."]

	// 查找“单词边界”之间的字符串
	reg, _ = regexp.Compile(`(?U)\b.+\b`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" " 世界！" "123" " " "Go"]

	// 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	reg, _ = regexp.Compile(`[^ ]{1,4}o`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找 Hello 或 Go
	reg, _ = regexp.Compile(`(?:Hell|G)o`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	// ["Hello" "Go"]

	// 查找 Hello 或 Go，替换为 Hellooo、Gooo
	reg, _ = regexp.Compile(`(?PHell|G)o`)
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "${n}ooo"))
	// "Hellooo 世界！123 Gooo."

	// 交换 Hello 和 Go
	reg, _ = regexp.Compile(`(Hello)(.*)(Go)`)
	fmt.Printf("%q\n", reg.ReplaceAllString(text, "$3$2$1"))
	// "Go 世界！123 Hello."

	// 特殊字符的查找
	reg, _ = regexp.Compile(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]\|]`)
	fmt.Printf("%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))
	// "----------------------"
}
