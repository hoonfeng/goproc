#!/usr/bin/env python3
"""
Python插件示例 - 使用标准SDK实现
"""

import sys
import os

# 添加SDK路径到Python路径
sdk_path = os.path.join(os.path.dirname(__file__), "..", "..", "..", "sdk", "python")
if os.path.exists(sdk_path):
    sys.path.insert(0, sdk_path)

# 导入标准SDK
try:
    from goproc_sdk import PluginSDK, register_function, start_plugin, wait_plugin
except ImportError:
    sys.exit(1)

# 定义插件函数
def add_numbers(params):
    """加法函数"""
    a = params.get('a', 0)
    b = params.get('b', 0)
    return a + b

def subtract_numbers(params):
    """减法函数"""
    a = params.get('a', 0)
    b = params.get('b', 0)
    return a - b

def multiply_numbers(params):
    """乘法函数"""
    a = params.get('a', 0)
    b = params.get('b', 0)
    return a * b

def divide_numbers(params):
    """除法函数"""
    a = params.get('a', 0)
    b = params.get('b', 1)
    if b == 0:
        raise ValueError("除数不能为零")
    return a / b

def get_current_time(params):
    """获取当前时间"""
    import datetime
    return datetime.datetime.now().isoformat()

def reverse_string(params):
    """反转字符串"""
    text = params.get('text', '')
    return text[::-1]

def uppercase_string(params):
    """转换为大写"""
    text = params.get('text', '')
    return text.upper()

def lowercase_string(params):
    """转换为小写"""
    text = params.get('text', '')
    return text.lower()

def get_string_length(params):
    """获取字符串长度"""
    text = params.get('text', '')
    return len(text)

def text_processing(params):
    """文本处理函数"""
    text = params.get('text', '')
    operation = params.get('operation', 'word_count')
    
    if operation == 'word_count':
        # 简单的单词计数（按空格分割）
        words = text.split()
        return len(words)
    elif operation == 'character_count':
        # 字符计数
        return len(text)
    elif operation == 'uppercase':
        # 转换为大写
        return text.upper()
    elif operation == 'lowercase':
        # 转换为小写
        return text.lower()
    else:
        # 默认返回字符计数
        return len(text)

def fibonacci_sequence(params):
    """生成斐波那契数列"""
    n = params.get('n', 10)
    if n <= 0:
        return []
    elif n == 1:
        return [0]
    elif n == 2:
        return [0, 1]
    
    sequence = [0, 1]
    for i in range(2, n):
        sequence.append(sequence[i-1] + sequence[i-2])
    
    return sequence

# 注册函数
register_function("add", add_numbers)
register_function("subtract", subtract_numbers)
register_function("multiply", multiply_numbers)
register_function("divide", divide_numbers)
register_function("datetime_utils", get_current_time)  # 修复函数名以匹配综合演示程序
register_function("reverse", reverse_string)
register_function("uppercase", uppercase_string)
register_function("lowercase", lowercase_string)
register_function("length", get_string_length)
register_function("fibonacci", fibonacci_sequence)
register_function("text_processing", text_processing)  # 添加text_processing函数以支持并发测试

if __name__ == "__main__":
    # 启动插件
    if start_plugin():
        # 等待插件停止
        from goproc_sdk import wait_plugin
        wait_plugin()
    else:
        sys.exit(1)