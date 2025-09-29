#!/usr/bin/env python3
"""
GoProc插件SDK - Python版本
提供简单的函数注册和插件启动功能
"""

import json
import socket
import struct
import sys
import os
import platform
import threading
import time
from typing import Dict, Callable, Any, Optional

# Windows命名管道支持
try:
    import win32file
    import win32pipe
    WINDOWS_PIPE_AVAILABLE = True
except ImportError:
    WINDOWS_PIPE_AVAILABLE = False

# Windows命名管道连接包装器
class PipeConnection:
    """Windows命名管道连接包装器"""
    
    def __init__(self, handle):
        self.handle = handle
        self._closed = False
    
    def send(self, data):
        """发送数据 / Send data"""
        if self._closed:
            raise IOError("管道连接已关闭")
        
        try:
            # WriteFile返回(error_code, bytes_written)
            # WriteFile returns (error_code, bytes_written)
            error_code, bytes_written = win32file.WriteFile(self.handle, data)
            if error_code == 0:  # 成功 / Success
                return bytes_written
            else:
                raise IOError(f"写入失败，错误代码: {error_code}")
        except Exception as e:
            raise IOError(f"发送数据失败: {e}")

    def sendall(self, data):
        """发送所有数据 / Send all data"""
        if self._closed:
            raise IOError("管道连接已关闭")
        
        total_sent = 0
        while total_sent < len(data):
            try:
                # WriteFile返回(error_code, bytes_written)
                # WriteFile returns (error_code, bytes_written)
                error_code, bytes_written = win32file.WriteFile(self.handle, data[total_sent:])
                if error_code == 0:  # 成功 / Success
                    total_sent += bytes_written
                else:
                    raise IOError(f"写入失败，错误代码: {error_code}")
            except Exception as e:
                raise IOError(f"发送数据失败: {e}")
        
        # 刷新缓冲区确保数据立即发送
        # Flush buffer to ensure data is sent immediately
        try:
            win32file.FlushFileBuffers(self.handle)
        except:
            pass  # 忽略刷新错误
        
        return total_sent
    
    def recv(self, bufsize):
        """接收数据 / Receive data"""
        if self._closed:
            raise IOError("管道连接已关闭")
        
        try:
            result, data = win32file.ReadFile(self.handle, bufsize)
            if result == 0:  # 成功 / Success
                return data
            else:
                raise IOError(f"读取失败，错误代码: {result}")
        except Exception as e:
            raise IOError(f"接收数据失败: {e}")

    def settimeout(self, timeout):
        """设置超时（兼容socket接口）"""
        # Windows命名管道不支持超时设置，这里只是接口兼容
        pass
    
    def close(self):
        """关闭连接"""
        if not self._closed:
            win32file.CloseHandle(self.handle)
            self._closed = True
    
    def fileno(self):
        """获取文件描述符"""
        return self.handle
    
    def getpeername(self):
        """获取对端地址"""
        return "named-pipe"
    
    def getsockname(self):
        """获取本地地址"""
        return "named-pipe"

# 消息类型常量
MESSAGE_TYPE_CALL = "call"
MESSAGE_TYPE_RESULT = "result"
MESSAGE_TYPE_ERROR = "error"
MESSAGE_TYPE_PING = "ping"
MESSAGE_TYPE_PONG = "pong"
MESSAGE_TYPE_REGISTER = "register"
MESSAGE_TYPE_REGISTER_ACK = "register_ack"
MESSAGE_TYPE_STOP = "stop"

class PluginSDK:
    """插件SDK主类"""
    
    def __init__(self):
        self.functions: Dict[str, Callable] = {}
        self.conn: Optional[socket.socket] = None
        self.address: str = ""
        self.running: bool = False
        self.registered: bool = False
        self.message_thread: Optional[threading.Thread] = None
    
    def register_function(self, name: str, handler: Callable) -> None:
        """注册函数"""
        self.functions[name] = handler
    
    def function(self, name: str = None):
        """
        函数装饰器，用于简化函数注册
        
        用法1（使用函数名作为装饰器参数）:
            @sdk.function("my_function")
            def my_function(params):
                return {"result": "success"}
        
        用法2（使用函数名作为装饰器）:
            @sdk.function
            def my_function(params):
                return {"result": "success"}
        """
        def decorator(func):
            # 如果name为None，使用函数名作为注册名
            function_name = name if name is not None else func.__name__
            self.register_function(function_name, func)
            return func
        
        # 如果name是函数，说明是用法2
        if callable(name):
            func = name
            function_name = func.__name__
            self.register_function(function_name, func)
            return func
        
        return decorator
    
    def encode_message(self, msg: Dict[str, Any]) -> bytes:
        """编码消息"""
        try:
            json_str = json.dumps(msg, ensure_ascii=False)
            return json_str.encode('utf-8')
        except Exception as e:
            return b''
    
    def decode_message(self, data: bytes) -> Optional[Dict[str, Any]]:
        """解码消息"""
        try:
            json_str = data.decode('utf-8')
            return json.loads(json_str)
        except Exception as e:
            return None
    
    def send_message(self, msg: Dict[str, Any]) -> bool:
        """发送消息"""
        if not self.conn:
            return False
        
        try:
            data = self.encode_message(msg)
            if not data:
                return False
            
            # 添加长度前缀
            length = len(data)
            header = struct.pack('>I', length)
            
            # 发送头部和数据
            self.conn.sendall(header + data)
            return True
        except Exception:
            return False
    
    def receive_message(self) -> Optional[Dict[str, Any]]:
        """接收消息"""
        if not self.conn:
            return None
        
        try:
            # 读取消息长度（4字节）
            header = b''
            while len(header) < 4:
                chunk = self.conn.recv(4 - len(header))
                if not chunk:
                    return None
                header += chunk
            
            # 解析消息长度
            length = struct.unpack('>I', header)[0]
            
            # 读取消息体
            data = b''
            while len(data) < length:
                chunk = self.conn.recv(length - len(data))
                if not chunk:
                    return None
                data += chunk
            
            # 解码消息
            return self.decode_message(data)
        except Exception:
            return None
    
    def handle_call_message(self, msg: Dict[str, Any]) -> None:
        """处理调用消息"""
        function_name = msg.get('function', '')
        params = msg.get('params', {})
        msg_id = msg.get('id', '')
        
        if not function_name or not msg_id:
            return
        
        # 查找函数
        handler = self.functions.get(function_name)
        if not handler:
            # 函数不存在，返回错误
            error_msg = {
                'type': MESSAGE_TYPE_ERROR,
                'id': msg_id,
                'error': f"函数 {function_name} 不存在"
            }
            self.send_message(error_msg)
            return
        
        try:
            # 调用函数
            result = handler(params)
            
            # 返回结果
            result_msg = {
                'type': MESSAGE_TYPE_RESULT,
                'id': msg_id,
                'result': result
            }
            self.send_message(result_msg)
        except Exception as e:
            # 返回错误
            error_msg = {
                'type': MESSAGE_TYPE_ERROR,
                'id': msg_id,
                'error': str(e)
            }
            self.send_message(error_msg)
    
    def handle_ping_message(self, msg: Dict[str, Any]) -> None:
        """处理ping消息"""
        pong_msg = {
            'type': MESSAGE_TYPE_PONG
        }
        self.send_message(pong_msg)
    
    def message_loop(self) -> None:
        """消息处理循环"""
        while self.running:
            try:
                msg = self.receive_message()
                if not msg:
                    # 如果连接断开，自动停止插件
                    self.running = False
                    break
                
                msg_type = msg.get('type', '')
                
                if msg_type == MESSAGE_TYPE_CALL:
                    self.handle_call_message(msg)
                elif msg_type == MESSAGE_TYPE_PING:
                    self.handle_ping_message(msg)
                elif msg_type == MESSAGE_TYPE_REGISTER_ACK:
                    # 处理注册确认消息
                    self.registered = True
                elif msg_type == MESSAGE_TYPE_STOP:
                    # 收到停止消息，优雅退出
                    self.running = False
                    break
            except Exception:
                # 处理异常，但不退出循环
                import time
                time.sleep(0.1)
    
    def create_listener_and_wait(self, address: str) -> bool:
        """创建监听器并等待连接"""
        try:
            self.address = address
            
            if address.startswith('\\\\.\\pipe'):
                # Windows命名管道
                if not WINDOWS_PIPE_AVAILABLE:
                    return False
                
                try:
                    # 创建命名管道
                    pipe_handle = win32pipe.CreateNamedPipe(
                        address,
                        win32pipe.PIPE_ACCESS_DUPLEX,
                        win32pipe.PIPE_TYPE_BYTE | win32pipe.PIPE_READMODE_BYTE | win32pipe.PIPE_WAIT,
                        1,  # 最大实例数
                        65536,  # 输出缓冲区大小
                        65536,  # 输入缓冲区大小
                        0,  # 默认超时
                        None  # 安全属性
                    )
                    
                    if pipe_handle == win32file.INVALID_HANDLE_VALUE:
                        return False
                    
                    # 等待客户端连接
                    result = win32pipe.ConnectNamedPipe(pipe_handle, None)
                    if result == 0:  # 成功
                        self.conn = PipeConnection(pipe_handle)
                        return True
                    else:
                        win32file.CloseHandle(pipe_handle)
                        return False
                    
                except Exception:
                    return False
                
            else:
                # Unix域套接字
                # 确保socket文件所在目录存在
                socket_dir = os.path.dirname(address)
                if socket_dir and not os.path.exists(socket_dir):
                    os.makedirs(socket_dir, exist_ok=True)
                
                if os.path.exists(address):
                    os.unlink(address)
                
                self.listener = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
                self.listener.bind(address)
                self.listener.listen(1)
                
                self.conn, _ = self.listener.accept()
            
            return True
            
        except Exception:
            return False
    
    def start(self) -> bool:
        """启动插件"""
        # 获取通信地址
        address = os.getenv('GOPROC_PLUGIN_ADDRESS', '')
        if not address and len(sys.argv) >= 2:
            address = sys.argv[1]
        
        if not address:
            return False
        
        # 创建监听器并等待连接
        if not self.create_listener_and_wait(address):
            return False
        
        # 立即发送注册消息，在启动消息循环之前
        if not self.send_register_message_simple():
            self.conn.close()
            return False
        
        # 启动消息处理循环
        self.running = True
        self.message_thread = threading.Thread(target=self.message_loop)
        self.message_thread.daemon = True
        self.message_thread.start()
        
        # 等待注册完成（最多10秒）
        start_time = time.time()
        while not self.registered and self.running and (time.time() - start_time < 10):
            time.sleep(0.1)
        
        return self.registered and self.running
    
    def stop(self) -> None:
        """停止插件"""
        self.running = False
        if self.conn:
            self.conn.close()
    
    def wait(self) -> None:
        """等待插件停止"""
        while self.running:
            time.sleep(0.1)

    def send_register_message(self) -> bool:
        """发送注册消息并等待确认"""
        try:
            function_names = list(self.functions.keys())
            
            # 按照Go SDK协议，使用params字段包含函数列表
            register_msg = {
                'type': MESSAGE_TYPE_REGISTER,
                'params': {
                    'functions': function_names
                }
            }
            
            # 发送消息
            result = self.send_message(register_msg)
            if not result:
                return False
            
            # 等待注册确认响应
            # 设置超时（10秒）
            start_time = time.time()
            while time.time() - start_time < 10:
                try:
                    # 设置读取超时（1秒）
                    if hasattr(self.conn, 'settimeout'):
                        self.conn.settimeout(1.0)
                    
                    response = self.receive_message()
                    if response:
                        if response.get('type') == MESSAGE_TYPE_REGISTER_ACK:
                            return True
                    
                    time.sleep(0.1)
                except socket.timeout:
                    # 超时是正常的，继续等待
                    continue
                except Exception:
                    return False
            
            return False
            
        except Exception:
            return False
    
    def send_register_message_simple(self) -> bool:
         """发送简化的注册消息"""
         try:
             function_names = list(self.functions.keys())
             
             register_msg = {
                 'type': MESSAGE_TYPE_REGISTER,
                 'params': {
                     'functions': function_names
                 }
             }
             
             return self.send_message(register_msg)
             
         except Exception as e:
             return False

# 全局SDK实例
_global_sdk = PluginSDK()

def register_function(name: str, handler: Callable) -> None:
    """注册函数（全局函数）"""
    _global_sdk.register_function(name, handler)

def start_plugin() -> bool:
    """启动插件（全局函数）"""
    return _global_sdk.start()

def stop_plugin() -> None:
    """停止插件（全局函数）"""
    _global_sdk.stop()

def wait_plugin() -> None:
    """等待插件停止（全局函数）"""
    _global_sdk.wait()

# 简化导入
register = register_function
start = start_plugin
stop = stop_plugin
wait = wait_plugin