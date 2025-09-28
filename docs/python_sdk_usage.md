# GoProc Python SDK 使用指南 / Python SDK Usage Guide

## 📖 概述 / Overview

GoProc Python SDK 允许开发者使用 Python 语言开发插件，与 GoProc 主程序进行高效通信。SDK 支持跨平台运行，在 Windows 上使用命名管道，在 Unix 系统上使用域套接字。

GoProc Python SDK is the Python language support library for the GoProc plugin system, providing simple and easy-to-use interfaces for developing Python plugins.

## 🏆 **性能表现 / Performance Benchmarks**

### 测试环境 / Test Environment
- **设备**: MSI 笔记本 / MSI Laptop  
- **处理器**: 13th Gen Intel(R) Core(TM) i7-13700H (2.40 GHz)
- **内存**: 32.0 GB (31.7 GB 可用) / 32.0 GB (31.7 GB available)

### Python插件性能 / Python Plugin Performance
- **QPS**: **100,000+** (10个实例并发测试)
- **延迟**: 低至 **1.90ms** 平均响应时间
- **成功率**: **100%** 无错误率

## 🔧 最新修复 (2024年12月)

### Windows 命名管道修复
- ✅ **架构重新设计**: 从客户端模式改为服务器模式
- ✅ **消息发送修复**: 修复 `sendall()` 方法的不完整发送问题
- ✅ **时序优化**: 解决连接建立和消息传输的竞态条件
- ✅ **性能提升**: 响应时间优化 20-30%，连接成功率达到 100%

## 🚀 快速开始

### 1. 环境准备

#### Windows 环境
```cmd
# 安装依赖
pip install -r requirements.txt

# requirements.txt 内容:
# requests==2.31.0
# pywin32==306
```

#### Unix 环境 (Linux/macOS)
```bash
# 安装依赖
pip install requests==2.31.0
```

### 2. 基本插件结构

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os

# 添加SDK路径
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', '..', 'sdk', 'python'))

from goproc_sdk import GoProc

def main():
    # 创建GoProc实例
    goproc = GoProc()
    
    # 注册函数
    goproc.register_function("my_function", my_function)
    
    # 启动插件
    goproc.start()

def my_function(params):
    """
    插件函数示例
    
    Args:
        params (dict): 输入参数字典
        
    Returns:
        dict: 返回结果，必须包含 'result' 或 'error' 字段
    """
    try:
        # 处理业务逻辑
        result = "处理结果"
        return {"result": result}
    except Exception as e:
        return {"error": str(e)}

if __name__ == "__main__":
    main()
```

## 📚 API 参考

### GoProc 类

#### 构造函数
```python
goproc = GoProc()
```

#### 方法

##### register_function(name, func)
注册插件函数

**参数:**
- `name` (str): 函数名称
- `func` (callable): 函数对象

**示例:**
```python
def add_numbers(params):
    a = params.get('a', 0)
    b = params.get('b', 0)
    return {"result": a + b}

goproc.register_function("add", add_numbers)
```

##### start()
启动插件，开始监听来自 GoProc 主程序的请求

**示例:**
```python
goproc.start()  # 阻塞运行
```

### 函数签名规范

插件函数必须遵循以下签名规范:

```python
def function_name(params):
    """
    Args:
        params (dict): 输入参数字典
        
    Returns:
        dict: 返回结果，格式如下:
        - 成功: {"result": any_value}
        - 失败: {"error": "error_message"}
    """
    pass
```

## 🛠 开发示例

### 示例 1: 数学计算插件

```python
#!/usr/bin/env python3
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', '..', 'sdk', 'python'))

from goproc_sdk import GoProc

def add(params):
    """加法运算"""
    try:
        a = params.get('a', 0)
        b = params.get('b', 0)
        result = a + b
        return {"result": result}
    except Exception as e:
        return {"error": f"加法运算失败: {str(e)}"}

def multiply(params):
    """乘法运算"""
    try:
        a = params.get('a', 1)
        b = params.get('b', 1)
        result = a * b
        return {"result": result}
    except Exception as e:
        return {"error": f"乘法运算失败: {str(e)}"}

def main():
    goproc = GoProc()
    
    # 注册多个函数
    goproc.register_function("add", add)
    goproc.register_function("multiply", multiply)
    
    print("数学计算插件已启动...")
    goproc.start()

if __name__ == "__main__":
    main()
```

### 示例 2: 字符串处理插件

```python
#!/usr/bin/env python3
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', '..', 'sdk', 'python'))

from goproc_sdk import GoProc

def reverse_string(params):
    """字符串反转"""
    try:
        text = params.get('text', '')
        result = text[::-1]
        return {"result": result}
    except Exception as e:
        return {"error": f"字符串反转失败: {str(e)}"}

def to_upper(params):
    """转换为大写"""
    try:
        text = params.get('text', '')
        result = text.upper()
        return {"result": result}
    except Exception as e:
        return {"error": f"大写转换失败: {str(e)}"}

def word_count(params):
    """单词计数"""
    try:
        text = params.get('text', '')
        words = text.split()
        result = len(words)
        return {"result": result}
    except Exception as e:
        return {"error": f"单词计数失败: {str(e)}"}

def main():
    goproc = GoProc()
    
    goproc.register_function("reverse", reverse_string)
    goproc.register_function("upper", to_upper)
    goproc.register_function("count_words", word_count)
    
    print("字符串处理插件已启动...")
    goproc.start()

if __name__ == "__main__":
    main()
```

## 🔍 调试和故障排除

### 常见问题

#### 1. 连接失败 (Windows)
```
Failed to connect to named pipe
```

**解决方案:**
- 确保安装了 `pywin32==306`
- 检查 GoProc 主程序是否正在运行
- 验证命名管道地址格式

#### 2. 模块导入错误
```
ModuleNotFoundError: No module named 'goproc_sdk'
```

**解决方案:**
```python
# 确保正确添加SDK路径
import sys
import os
sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', '..', 'sdk', 'python'))
```

#### 3. 函数注册失败
```
Function registration failed
```

**解决方案:**
- 检查函数签名是否正确
- 确保函数返回格式符合规范
- 验证函数名称是否唯一

### 调试技巧

#### 启用详细日志
```python
import logging

# 配置日志
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

# 在函数中添加日志
def my_function(params):
    logging.info(f"收到参数: {params}")
    # ... 处理逻辑
    logging.info(f"返回结果: {result}")
    return {"result": result}
```

#### 参数验证
```python
def validate_params(params, required_keys):
    """验证参数"""
    for key in required_keys:
        if key not in params:
            raise ValueError(f"缺少必需参数: {key}")
    return True

def my_function(params):
    try:
        validate_params(params, ['input_data'])
        # ... 处理逻辑
        return {"result": result}
    except ValueError as e:
        return {"error": str(e)}
```

## 📊 性能优化

### 最佳实践

1. **避免阻塞操作**: 在插件函数中避免长时间阻塞操作
2. **异常处理**: 始终包含适当的异常处理
3. **参数验证**: 验证输入参数的类型和范围
4. **资源管理**: 及时释放文件句柄和网络连接
5. **日志记录**: 适度使用日志，避免过度输出

### 性能监控

```python
import time

def timed_function(params):
    """带性能监控的函数"""
    start_time = time.time()
    
    try:
        # 业务逻辑
        result = process_data(params)
        
        end_time = time.time()
        execution_time = end_time - start_time
        
        return {
            "result": result,
            "execution_time": execution_time
        }
    except Exception as e:
        return {"error": str(e)}
```

## 🔐 安全考虑

### 输入验证
```python
def secure_function(params):
    """安全的函数示例"""
    # 验证参数类型
    if not isinstance(params, dict):
        return {"error": "参数必须是字典类型"}
    
    # 验证必需字段
    if 'data' not in params:
        return {"error": "缺少data字段"}
    
    # 验证数据类型
    data = params['data']
    if not isinstance(data, str):
        return {"error": "data字段必须是字符串"}
    
    # 验证数据长度
    if len(data) > 1000:
        return {"error": "数据长度不能超过1000字符"}
    
    # 处理数据
    result = process_secure_data(data)
    return {"result": result}
```

### 避免代码注入
```python
import re

def safe_eval_function(params):
    """安全的表达式计算"""
    expression = params.get('expression', '')
    
    # 只允许数字、基本运算符和括号
    if not re.match(r'^[0-9+\-*/().\s]+$', expression):
        return {"error": "表达式包含非法字符"}
    
    try:
        # 使用安全的eval替代方案
        result = eval(expression, {"__builtins__": {}})
        return {"result": result}
    except Exception as e:
        return {"error": f"表达式计算失败: {str(e)}"}
```

## 📖 更多资源

- [GoProc 主文档](../README.md)
- [跨平台使用指南](CROSS_PLATFORM_USAGE_GUIDE.md)
- [技术修复详情](technical_fix_details.md)
- [示例代码](../examples/)

---

**版本**: 1.0.0  
**最后更新**: 2024年12月  
**维护者**: GoProc开发团队