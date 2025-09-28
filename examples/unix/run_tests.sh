#!/bin/bash

# Unix平台测试运行脚本
echo "=== GoProc Unix平台测试 ==="

# 切换到goproc根目录
cd ../../

echo "当前工作目录: $(pwd)"

# 运行各个测试
echo ""
echo "=== 运行跨平台兼容性测试 ==="
go run examples/unix/test_cross_platform.go

echo ""
echo "=== 运行Python简单测试 ==="
go run examples/unix/test_python_simple.go

echo ""
echo "=== 运行Python插件测试 ==="
go run examples/unix/test_python_plugin.go

echo ""
echo "=== 运行Node.js插件测试 ==="
go run examples/unix/test_nodejs_plugin.go

echo ""
echo "=== 运行UUID测试 ==="
go run examples/unix/test_uuid.go

echo ""
echo "=== 运行通信测试 ==="
go run examples/unix/test_communication.go

echo ""
echo "=== 所有测试完成 ==="