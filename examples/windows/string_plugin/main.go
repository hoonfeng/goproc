package main

import (
	"fmt"
	"strings"
	"unicode"

	"goproc/sdk"
)

// ToUpper 转换为大写
func ToUpper(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	return strings.ToUpper(str), nil
}

// ToLower 转换为小写
func ToLower(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	return strings.ToLower(str), nil
}

// Reverse 字符串反转
func Reverse(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

// Trim 去除空白字符
func Trim(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	return strings.TrimSpace(str), nil
}

// Replace 字符串替换
func Replace(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	oldStr, ok := params["old"].(string)
	if !ok {
		return nil, fmt.Errorf("参数old缺失或类型错误")
	}
	newStr, ok := params["new"].(string)
	if !ok {
		return nil, fmt.Errorf("参数new缺失或类型错误")
	}
	return strings.ReplaceAll(str, oldStr, newStr), nil
}

// Split 字符串分割
func Split(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	sep, ok := params["sep"].(string)
	if !ok {
		return nil, fmt.Errorf("参数sep缺失或类型错误")
	}
	return strings.Split(str, sep), nil
}

// Join 字符串连接
func Join(params map[string]interface{}) (interface{}, error) {
	arr, ok := params["arr"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("参数arr缺失或类型错误")
	}
	sep, ok := params["sep"].(string)
	if !ok {
		return nil, fmt.Errorf("参数sep缺失或类型错误")
	}

	strArr := make([]string, len(arr))
	for i, v := range arr {
		strArr[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(strArr, sep), nil
}

// Contains 检查是否包含子串
func Contains(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}
	substr, ok := params["substr"].(string)
	if !ok {
		return nil, fmt.Errorf("参数substr缺失或类型错误")
	}
	return strings.Contains(str, substr), nil
}

// CountWords 统计单词数量
func CountWords(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}

	words := strings.Fields(str)
	return len(words), nil
}

// IsPalindrome 检查是否是回文
func IsPalindrome(params map[string]interface{}) (interface{}, error) {
	str, ok := params["str"].(string)
	if !ok {
		return nil, fmt.Errorf("参数str缺失或类型错误")
	}

	// 移除空格和标点，转换为小写
	cleanStr := ""
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleanStr += string(unicode.ToLower(r))
		}
	}

	// 检查是否是回文
	runes := []rune(cleanStr)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false, nil
		}
	}
	return true, nil
}

func main() {
	// 注册字符串处理函数
	err := sdk.RegisterFunction("toUpper", ToUpper)
	if err != nil {
		panic(fmt.Sprintf("注册toUpper函数失败: %v", err))
	}

	err = sdk.RegisterFunction("toLower", ToLower)
	if err != nil {
		panic(fmt.Sprintf("注册toLower函数失败: %v", err))
	}

	err = sdk.RegisterFunction("reverse", Reverse)
	if err != nil {
		panic(fmt.Sprintf("注册reverse函数失败: %v", err))
	}

	err = sdk.RegisterFunction("trim", Trim)
	if err != nil {
		panic(fmt.Sprintf("注册trim函数失败: %v", err))
	}

	err = sdk.RegisterFunction("replace", Replace)
	if err != nil {
		panic(fmt.Sprintf("注册replace函数失败: %v", err))
	}

	err = sdk.RegisterFunction("split", Split)
	if err != nil {
		panic(fmt.Sprintf("注册split函数失败: %v", err))
	}

	err = sdk.RegisterFunction("join", Join)
	if err != nil {
		panic(fmt.Sprintf("注册join函数失败: %v", err))
	}

	err = sdk.RegisterFunction("contains", Contains)
	if err != nil {
		panic(fmt.Sprintf("注册contains函数失败: %v", err))
	}

	err = sdk.RegisterFunction("countWords", CountWords)
	if err != nil {
		panic(fmt.Sprintf("注册countWords函数失败: %v", err))
	}

	err = sdk.RegisterFunction("isPalindrome", IsPalindrome)
	if err != nil {
		panic(fmt.Sprintf("注册isPalindrome函数失败: %v", err))
	}

	// 启动插件SDK
	if err := sdk.Start(); err != nil {
		panic(fmt.Sprintf("启动插件SDK失败: %v", err))
	}

	// 阻塞主线程，等待插件停止
	sdk.Wait()
}
