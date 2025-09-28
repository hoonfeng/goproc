package main

import (
	"fmt"
	"math"

	"goproc/sdk"
)

// Add 加法函数
func Add(params map[string]interface{}) (interface{}, error) {
	a, ok := params["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数a缺失或类型错误")
	}
	b, ok := params["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数b缺失或类型错误")
	}
	return a + b, nil
}

// Subtract 减法函数
func Subtract(params map[string]interface{}) (interface{}, error) {
	a, ok := params["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数a缺失或类型错误")
	}
	b, ok := params["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数b缺失或类型错误")
	}
	return a - b, nil
}

// Multiply 乘法函数
func Multiply(params map[string]interface{}) (interface{}, error) {
	a, ok := params["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数a缺失或类型错误")
	}
	b, ok := params["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数b缺失或类型错误")
	}
	return a * b, nil
}

// Divide 除法函数
func Divide(params map[string]interface{}) (interface{}, error) {
	a, ok := params["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数a缺失或类型错误")
	}
	b, ok := params["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数b缺失或类型错误")
	}
	if b == 0 {
		return nil, fmt.Errorf("除数不能为零")
	}
	return a / b, nil
}

// Power 幂运算函数
func Power(params map[string]interface{}) (interface{}, error) {
	base, ok := params["base"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数base缺失或类型错误")
	}
	exponent, ok := params["exponent"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数exponent缺失或类型错误")
	}
	return math.Pow(base, exponent), nil
}

// SquareRoot 平方根函数
func SquareRoot(params map[string]interface{}) (interface{}, error) {
	num, ok := params["num"].(float64)
	if !ok {
		return nil, fmt.Errorf("参数num缺失或类型错误")
	}
	if num < 0 {
		return nil, fmt.Errorf("负数没有实数平方根")
	}
	return math.Sqrt(num), nil
}

func main() {
	// 注册数学函数
	err := sdk.RegisterFunction("add", Add)
	if err != nil {
		panic(fmt.Sprintf("注册add函数失败: %v", err))
	}

	err = sdk.RegisterFunction("subtract", Subtract)
	if err != nil {
		panic(fmt.Sprintf("注册subtract函数失败: %v", err))
	}

	err = sdk.RegisterFunction("multiply", Multiply)
	if err != nil {
		panic(fmt.Sprintf("注册multiply函数失败: %v", err))
	}

	err = sdk.RegisterFunction("divide", Divide)
	if err != nil {
		panic(fmt.Sprintf("注册divide函数失败: %v", err))
	}

	err = sdk.RegisterFunction("power", Power)
	if err != nil {
		panic(fmt.Sprintf("注册power函数失败: %v", err))
	}

	err = sdk.RegisterFunction("sqrt", SquareRoot)
	if err != nil {
		panic(fmt.Sprintf("注册sqrt函数失败: %v", err))
	}

	// 启动插件SDK
	if err := sdk.Start(); err != nil {
		panic(err)
	}

	// 等待插件运行直到停止
	sdk.Wait()
}