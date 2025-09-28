#!/bin/bash

echo "========================================"
echo "GoProc 示例程序构建脚本"
echo "========================================"
echo

# 检查Go环境
echo "检查Go环境..."
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go编译器，请先安装Go语言环境"
    exit 1
fi

# 检查Python环境
echo "检查Python环境..."
if ! command -v python3 &> /dev/null; then
    echo "警告: 未找到Python3，部分示例可能无法运行"
fi

# 检查Node.js环境
echo "检查Node.js环境..."
if ! command -v node &> /dev/null; then
    echo "警告: 未找到Node.js，部分示例可能无法运行"
fi

echo
echo "开始构建Go插件..."

# 构建基础数学插件
echo "构建基础数学插件..."
cd basic_math_plugin
go build -o math_plugin math_plugin.go
if [ $? -ne 0 ]; then
    echo "错误: 基础数学插件构建失败"
    exit 1
fi
cd ..

# 构建并发处理插件
echo "构建并发处理插件..."
cd concurrent_processor
go build -o concurrent_processor concurrent_processor.go
if [ $? -ne 0 ]; then
    echo "错误: 并发处理插件构建失败"
    exit 1
fi
cd ..

# 构建错误处理插件
echo "构建错误处理插件..."
cd error_handler
go build -o error_handler error_handler.go
if [ $? -ne 0 ]; then
    echo "错误: 错误处理插件构建失败"
    exit 1
fi
cd ..

# 构建演示应用程序
echo "构建演示应用程序..."
cd demo_app
go build -o demo_app main.go
if [ $? -ne 0 ]; then
    echo "错误: 演示应用程序构建失败"
    exit 1
fi
cd ..

# 构建性能测试程序
echo "构建性能测试程序..."
cd performance_test
go build -o performance_test main.go
if [ $? -ne 0 ]; then
    echo "错误: 性能测试程序构建失败"
    exit 1
fi
cd ..

echo
echo "========================================"
echo "构建完成！"
echo "========================================"
echo
echo "已构建的程序:"
echo "  - basic_math_plugin/math_plugin"
echo "  - concurrent_processor/concurrent_processor"
echo "  - error_handler/error_handler"
echo "  - demo_app/demo_app"
echo "  - performance_test/performance_test"
echo
echo "运行示例:"
echo "  演示程序: ./demo_app/demo_app"
echo "  性能测试: ./performance_test/performance_test"
echo
echo "给程序添加执行权限..."
chmod +x basic_math_plugin/math_plugin
chmod +x concurrent_processor/concurrent_processor
chmod +x error_handler/error_handler
chmod +x demo_app/demo_app
chmod +x performance_test/performance_test

echo "完成！"