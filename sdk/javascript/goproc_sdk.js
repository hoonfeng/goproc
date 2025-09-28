#!/usr/bin/env node

const net = require('net');
const fs = require('fs');

// 装饰器工厂函数（简化版，使用函数包装器）
function pluginFunction(name) {
    return function(target, context) {
        // 这是一个简化的装饰器实现，使用函数包装器
        // 在实际使用中，用户可以使用 @pluginFunction('functionName') 语法
        // 但需要启用TypeScript或实验性装饰器功能
        
        // 返回一个包装函数，用于注册
        return function(...args) {
            return target.apply(this, args);
        };
    };
}

// 简化的装饰器替代方案
function createPluginFunction(sdk, name) {
    return function(handler) {
        const functionName = name || handler.name || 'anonymous';
        sdk.registerFunction(functionName, handler);
        return handler;
    };
}

// 消息类型常量
const MESSAGE_TYPE_CALL = "call";
const MESSAGE_TYPE_RESULT = "result";
const MESSAGE_TYPE_ERROR = "error";
const MESSAGE_TYPE_PING = "ping";
const MESSAGE_TYPE_PONG = "pong";
const MESSAGE_TYPE_REGISTER = "register";
const MESSAGE_TYPE_STOP = "stop";

class PluginSDK {
    constructor() {
        this.functions = {};
        this.conn = null;
        this.address = "";
        this.running = false;
        
        // 支持装饰器语法
        this.pluginFunction = function(name) {
            return function(target, context) {
                // 这是一个简化的装饰器实现
                // 在实际使用中，用户可以使用 @sdk.pluginFunction('functionName') 语法
                
                // 如果提供了函数名，直接注册
                if (name && typeof target === 'function') {
                    this.registerFunction(name, target);
                    return target;
                }
                
                // 返回一个包装函数，用于注册
                return function(...args) {
                    return target.apply(this, args);
                };
            }.bind(this);
        };
    }

    registerFunction(name, handler) {
        this.functions[name] = handler;
    }

    registerFunctionDecorator(name, handler) {
        /**
         * 函数装饰器，用于简化函数注册
         * 
         * 用法1（使用函数名作为参数）:
         *     sdk.registerFunctionDecorator("myFunction")(function(params) {
         *         return { result: "success" };
         *     });
         * 
         * 用法2（使用函数名作为装饰器）:
         *     sdk.registerFunctionDecorator(function myFunction(params) {
         *         return { result: "success" };
         *     });
         * 
         * 用法3（链式调用）:
         *     sdk.registerFunctionDecorator("myFunction", function(params) {
         *         return { result: "success" };
         *     });
         */
        
        // 如果第一个参数是函数，说明是用法2
        if (typeof name === 'function') {
            const func = name;
            const functionName = func.name || 'anonymous';
            this.registerFunction(functionName, func);
            return this;
        }
        
        // 如果第二个参数是函数，说明是用法3
        if (arguments.length === 2 && typeof handler === 'function') {
            this.registerFunction(name, handler);
            return this;
        }
        
        // 否则返回装饰器函数（用法1）
        return (handlerFunc) => {
            const functionName = name || handlerFunc.name || 'anonymous';
            this.registerFunction(functionName, handlerFunc);
            return handlerFunc;
        };
    }

    encodeMessage(msg) {
        const data = JSON.stringify(msg);
        const length = Buffer.byteLength(data, 'utf8');
        const header = Buffer.alloc(4);
        header.writeUInt32BE(length, 0);
        
        return Buffer.concat([header, Buffer.from(data, 'utf8')]);
    }

    decodeMessage(buffer) {
        const data = buffer.toString('utf8');
        return JSON.parse(data);
    }

    sendMessage(msg) {
        if (!this.conn) {
            return false;
        }

        try {
            const data = this.encodeMessage(msg);
            this.conn.write(data);
            return true;
        } catch (error) {
            return false;
        }
    }

    // 设置消息接收器（异步事件驱动）
    setupMessageReceiver() {
        if (!this.conn) return;
        
        let buffer = Buffer.alloc(0);
        let expectedLength = 0;
        
        this.conn.on('data', (data) => {
            //console.log(`接收到数据块，长度: ${data.length}`);
            buffer = Buffer.concat([buffer, data]);
            
            // 处理所有完整的消息
            while (buffer.length >= 4) {
                if (expectedLength === 0) {
                    // 读取消息头
                    expectedLength = buffer.readUInt32BE(0);
                    //console.log(`消息头长度: ${expectedLength}`);
                }
                
                // 检查消息是否完整
                if (buffer.length >= 4 + expectedLength) {
                    // 提取消息体
                    const messageData = buffer.slice(4, 4 + expectedLength);
                    const msg = this.decodeMessage(messageData);
                    //console.log(`解码消息:`, JSON.stringify(msg));
                    
                    // 处理消息
                    this.handleIncomingMessage(msg);
                    
                    // 移除已处理的消息
                    buffer = buffer.slice(4 + expectedLength);
                    expectedLength = 0;
                } else {
                    // 消息不完整，等待更多数据
                    break;
                }
            }
        });
        
        this.conn.on('error', (error) => {
            console.error('连接错误:', error);
        });
        
        this.conn.on('close', () => {
            //console.log('连接已关闭');
            this.running = false;
        });
    }

    async handleCallMessage(msg) {
        const functionName = msg.function || '';
        const params = msg.params || {};
        const msgId = msg.id || '';

        //console.log(`[SDK] 处理函数调用: ${functionName}, 参数:`, JSON.stringify(params));

        if (!this.functions[functionName]) {
            const errorMsg = {
                type: MESSAGE_TYPE_ERROR,
                id: msgId,
                error: `函数 ${functionName} 不存在`
            };
            //console.log(`[SDK] 函数不存在: ${functionName}`);
            this.sendMessage(errorMsg);
            return;
        }

        try {
            // 调用函数并获取结果
            let result;
            try {
                result = this.functions[functionName](params);
            } catch (syncError) {
                //console.log(`[SDK] 函数 ${functionName} 同步执行失败:`, syncError.message);
                const errorMsg = {
                    type: MESSAGE_TYPE_ERROR,
                    id: msgId,
                    error: syncError.message
                };
                this.sendMessage(errorMsg);
                return;
            }
            
            //console.log(`[SDK] 函数 ${functionName} 执行结果类型: ${typeof result}`);
            //console.log(`[SDK] 函数 ${functionName} 是否为Promise: ${result instanceof Promise}`);
            //console.log(`[SDK] 函数 ${functionName} 是否有then方法: ${result && typeof result.then === 'function'}`);
            
            // 更严格的Promise检测
            if (result instanceof Promise || (result && typeof result.then === 'function' && typeof result.catch === 'function')) {
                //console.log(`[SDK] 检测到异步函数: ${functionName}`);
                
                // 等待异步函数完成
                try {
                    const actualResult = await result;
                    //console.log(`[SDK] 异步函数 ${functionName} 执行完成，结果:`, JSON.stringify(actualResult));
                    
                    const resultMsg = {
                        type: MESSAGE_TYPE_RESULT,
                        id: msgId,
                        result: actualResult
                    };
                    this.sendMessage(resultMsg);
                } catch (asyncError) {
                    //console.log(`[SDK] 异步函数 ${functionName} 执行失败:`, asyncError.message);
                    
                    const errorMsg = {
                        type: MESSAGE_TYPE_ERROR,
                        id: msgId,
                        error: asyncError.message
                    };
                    this.sendMessage(errorMsg);
                }
            } else {
                // 同步函数，直接发送结果
                //console.log(`[SDK] 同步函数 ${functionName} 执行完成，结果:`, JSON.stringify(result));
                
                const resultMsg = {
                    type: MESSAGE_TYPE_RESULT,
                    id: msgId,
                    result: result
                };
                this.sendMessage(resultMsg);
            }
        } catch (error) {
            //console.log(`[SDK] 函数 ${functionName} 处理失败:`, error.message);
            
            const errorMsg = {
                type: MESSAGE_TYPE_ERROR,
                id: msgId,
                error: error.message
            };
            this.sendMessage(errorMsg);
        }
    }

    handlePingMessage(msg) {
        const pongMsg = {
            type: MESSAGE_TYPE_PONG,
            id: msg.id || ''
        };
        this.sendMessage(pongMsg);
    }

    // 处理接收到的消息
    handleIncomingMessage(msg) {
        if (!msg || !msg.type) {
            //console.log('收到无效消息');
            return;
        }

        const msgType = msg.type || '';
        //console.log(`处理消息类型: ${msgType}, ID: ${msg.id}`);

        switch (msgType) {
            case MESSAGE_TYPE_CALL:
                //console.log(`处理函数调用: ${msg.function}`);
                // 异步调用handleCallMessage，不等待结果
                this.handleCallMessage(msg).catch(error => {
                    //console.log(`[SDK] 处理函数调用异常: ${error.message}`);
                    // 发送错误响应
                    const errorMsg = {
                        type: MESSAGE_TYPE_ERROR,
                        id: msg.id || '',
                        error: `处理函数调用异常: ${error.message}`
                    };
                    this.sendMessage(errorMsg);
                });
                break;
            case MESSAGE_TYPE_PING:
                this.handlePingMessage(msg);
                break;
            case MESSAGE_TYPE_STOP:
                // 收到停止消息，优雅退出
                //console.log('收到停止消息');
                this.running = false;
                break;
            default:
                //console.log(`未知消息类型: ${msgType}`);
        }
    }

    createListenerAndWait() {
            // 获取通信地址
            let address = process.env.GOPROC_PLUGIN_ADDRESS;
            
            if (!address && process.argv.length >= 3) {
                address = process.argv[2]; // 命令行参数是第三个（0: node, 1: 脚本文件, 2: 地址）
            }

            //console.log(`创建监听器并等待连接，地址: ${address}`);

            if (!address) {
                throw new Error('未提供通信地址');
            }

            this.address = address;

            return new Promise((resolve, reject) => {
                try {
                    // 根据地址格式判断通信类型
                    if (address.includes('pipe') || address.startsWith('\\\\\\.\\\\pipe') || address.startsWith('\\.\\pipe')) {
                        // Windows命名管道 - 创建服务器
                        //console.log('使用Windows命名管道监听');
                        
                        // 创建命名管道服务器
                        const server = net.createServer((socket) => {
                            //console.log('命名管道客户端已连接');
                            this.conn = socket;
                            resolve(); // 连接建立后resolve
                        });
                        
                        // 在Windows上，net.createServer可以直接监听命名管道
                        server.listen(address, () => {
                            //console.log(`命名管道服务器正在监听: ${address}`);
                        });
                        
                        server.on('error', (error) => {
                            //console.log(`命名管道服务器错误: ${error.message}`);
                            reject(error);
                        });
                        
                    } else {
                        // Unix域套接字（默认）
                        //console.log('使用Unix域套接字监听');
                        
                        // 删除已存在的socket文件
                        if (fs.existsSync(address)) {
                            fs.unlinkSync(address);
                        }
                        
                        const server = net.createServer((socket) => {
                            //console.log('Unix域套接字客户端已连接');
                            this.conn = socket;
                            resolve(); // 连接建立后resolve
                        });
                        
                        server.listen(address, () => {
                            //console.log(`Unix域套接字服务器正在监听: ${address}`);
                        });
                        
                        server.on('error', (error) => {
                            //console.log(`Unix域套接字服务器错误: ${error.message}`);
                            reject(error);
                        });
                    }

                } catch (error) {
                    //console.log(`创建监听器异常: ${error.message}`);
                    reject(error);
                }
            });
    }

    sendRegisterMessage() {
        const functionNames = Object.keys(this.functions);
        
        const registerMsg = {
            type: MESSAGE_TYPE_REGISTER,
            params: {
                functions: functionNames
            }
        };
        
        return this.sendMessage(registerMsg);
    }

    async start() {
        try {
            // 创建监听器并等待连接（服务器模式）
            await this.createListenerAndWait();
            
            // 设置消息接收器（异步事件驱动）
            this.setupMessageReceiver();
            
            // 发送注册消息
            if (!this.sendRegisterMessage()) {
                this.conn.end();
                return false;
            }
            
            // 启动运行状态
            this.running = true;
            //console.log('插件启动成功，等待消息...');
            
            return true;
        } catch (error) {
            //console.log('start方法异常:', error.message);
            return false;
        }
    }

    stop() {
        this.running = false;
        if (this.conn) {
            this.conn.end();
        }
    }

    wait() {
        while (this.running) {
            // 同步延迟
            const start = Date.now();
            while (Date.now() - start < 100) {
                // 空循环实现同步延迟
            }
        }
    }
}

module.exports = { PluginSDK };