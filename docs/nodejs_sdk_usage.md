# GoProc Node.js SDK 使用指南 / GoProc Node.js SDK Usage Guide

## 概述 / Overview

GoProc Node.js SDK 是一个高性能的插件开发框架，支持跨平台通信（Windows命名管道和Unix域套接字）。SDK提供了简洁的API接口和装饰器语法，让JavaScript/TypeScript开发者能够轻松开发插件。

GoProc Node.js SDK is a high-performance plugin development framework that supports cross-platform communication (Windows Named Pipes and Unix Domain Sockets). The SDK provides a clean API interface and decorator syntax, allowing JavaScript/TypeScript developers to easily develop plugins.

## 🏆 性能表现 / Performance Benchmarks

### 测试环境 / Test Environment
- **设备**: MSI 笔记本 / MSI Laptop
- **处理器**: 13th Gen Intel(R) Core(TM) i7-13700H (2.40 GHz)
- **内存**: 32.0 GB (31.7 GB 可用) / 32.0 GB (31.7 GB available)
- **系统**: 64位操作系统, 基于 x64 的处理器 / 64-bit OS, x64-based processor

### 性能指标 / Performance Metrics
- **QPS**: **100,000+** (10个实例并发测试)
- **延迟**: 低至 **0.01ms** 平均响应时间
- **成功率**: **100%** 无错误率

## 🚀 快速开始 / Quick Start

### 1. 安装依赖 / Install Dependencies

```bash
npm install
```

### 2. 导入SDK / Import SDK

```javascript
const { PluginSDK, registerFunction } = require('./goproc_sdk');
// 或者使用ES6模块语法 / Or use ES6 module syntax
// import { PluginSDK, registerFunction } from './goproc_sdk.js';
```

### 3. 注册函数 / Register Functions

#### 方法一：使用装饰器语法 / Method 1: Using Decorator Syntax

```javascript
const { PluginSDK, registerFunction } = require('./goproc_sdk');

// 加法函数 / Addition function
const add = registerFunction('add', (params) => {
    const { a, b } = params;
    
    if (typeof a !== 'number' || typeof b !== 'number') {
        throw new Error('参数错误: a和b必须为数字');
    }
    
    return a + b;
});

// 减法函数 / Subtraction function
const subtract = registerFunction('subtract', (params) => {
    const { a, b } = params;
    
    if (typeof a !== 'number' || typeof b !== 'number') {
        throw new Error('参数错误: a和b必须为数字');
    }
    
    return a - b;
});

// 启动插件 / Start plugin
const sdk = new PluginSDK();
sdk.start().then(() => {
    console.log('插件启动成功');
}).catch(err => {
    console.error('插件启动失败:', err);
});
```

#### 方法二：使用SDK实例 / Method 2: Using SDK Instance

```javascript
const { PluginSDK } = require('./goproc_sdk');

async function main() {
    const sdk = new PluginSDK();
    
    // 注册乘法函数 / Register multiplication function
    sdk.registerFunction('multiply', (params) => {
        const { a, b } = params;
        
        if (typeof a !== 'number' || typeof b !== 'number') {
            throw new Error('参数错误: a和b必须为数字');
        }
        
        return a * b;
    });
    
    // 注册除法函数 / Register division function
    sdk.registerFunction('divide', (params) => {
        const { a, b } = params;
        
        if (typeof a !== 'number' || typeof b !== 'number') {
            throw new Error('参数错误: a和b必须为数字');
        }
        
        if (b === 0) {
            throw new Error('除数不能为零');
        }
        
        return a / b;
    });
    
    try {
        await sdk.start();
        console.log('插件启动成功');
    } catch (err) {
        console.error('插件启动失败:', err);
    }
}

main();
```

## 📚 API 参考 / API Reference

### 核心类 / Core Classes

#### PluginSDK

```javascript
class PluginSDK {
    constructor()
    registerFunction(name, handler)
    start()
    stop()
    // 私有方法 / Private methods
}
```

插件SDK主要类，提供插件开发的核心功能。

Main plugin SDK class providing core functionality for plugin development.

### 构造函数 / Constructor

#### new PluginSDK()

```javascript
const sdk = new PluginSDK();
```

创建新的插件SDK实例。

Create a new plugin SDK instance.

### 实例方法 / Instance Methods

#### registerFunction

```javascript
sdk.registerFunction(name, handler)
```

注册插件函数到SDK实例。

Register a plugin function to the SDK instance.

**参数 / Parameters:**
- `name` (string): 函数名称 / Function name
- `handler` (function): 函数处理器 / Function handler

**函数处理器签名 / Function Handler Signature:**
```javascript
function handler(params) {
    // params: 函数参数对象 / Function parameter object
    // 返回值: 任意类型 / Return value: any type
    // 抛出错误: throw new Error(message) / Throw error: throw new Error(message)
}
```

#### start

```javascript
await sdk.start()
```

启动SDK实例（异步方法）。

Start the SDK instance (asynchronous method).

**返回值 / Returns:**
- `Promise<void>`: Promise对象 / Promise object

#### stop

```javascript
sdk.stop()
```

停止SDK实例。

Stop the SDK instance.

### 全局函数 / Global Functions

#### registerFunction (装饰器 / Decorator)

```javascript
const functionHandler = registerFunction(name, handler)
```

注册插件函数的装饰器函数。

Decorator function for registering plugin functions.

**参数 / Parameters:**
- `name` (string): 函数名称 / Function name
- `handler` (function): 函数处理器 / Function handler

**返回值 / Returns:**
- `function`: 装饰后的函数 / Decorated function

## 🔧 高级用法 / Advanced Usage

### 复杂数据类型处理 / Complex Data Type Handling

```javascript
const { PluginSDK } = require('./goproc_sdk');

const sdk = new PluginSDK();

// 处理对象数据 / Handle object data
sdk.registerFunction('processUser', (params) => {
    const { user } = params;
    
    if (!user || typeof user !== 'object') {
        throw new Error('无效的用户数据格式');
    }
    
    const { name, age } = user;
    
    if (!name || typeof age !== 'number') {
        throw new Error('用户数据缺少必要字段');
    }
    
    // 处理业务逻辑 / Process business logic
    return {
        processed: true,
        message: `用户 ${name}，年龄 ${age} 岁，处理完成`,
        timestamp: Date.now()
    };
});

// 处理数组数据 / Handle array data
sdk.registerFunction('sumArray', (params) => {
    const { numbers } = params;
    
    if (!Array.isArray(numbers)) {
        throw new Error('参数 numbers 必须是数组');
    }
    
    let sum = 0;
    for (let i = 0; i < numbers.length; i++) {
        if (typeof numbers[i] !== 'number') {
            throw new Error(`数组第 ${i} 个元素不是数字`);
        }
        sum += numbers[i];
    }
    
    return sum;
});

// 处理嵌套数据 / Handle nested data
sdk.registerFunction('processNestedData', (params) => {
    const { data } = params;
    
    if (!data || typeof data !== 'object') {
        throw new Error('参数 data 必须是对象');
    }
    
    const result = {};
    
    // 递归处理嵌套对象 / Recursively process nested objects
    function processObject(obj, prefix = '') {
        for (const [key, value] of Object.entries(obj)) {
            const fullKey = prefix ? `${prefix}.${key}` : key;
            
            if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
                processObject(value, fullKey);
            } else {
                result[fullKey] = value;
            }
        }
    }
    
    processObject(data);
    return result;
});
```

### 异步处理 / Asynchronous Processing

```javascript
const { PluginSDK } = require('./goproc_sdk');

const sdk = new PluginSDK();

// 异步函数处理 / Asynchronous function handling
sdk.registerFunction('asyncProcess', async (params) => {
    const { taskId = `task_${Date.now()}`, delay = 1000 } = params;
    
    console.log(`开始处理任务: ${taskId}`);
    
    // 模拟异步操作 / Simulate asynchronous operation
    await new Promise(resolve => setTimeout(resolve, delay));
    
    console.log(`任务完成: ${taskId}`);
    
    return {
        taskId,
        status: 'completed',
        message: '任务处理完成',
        completedAt: new Date().toISOString()
    };
});

// 文件操作示例 / File operation example
sdk.registerFunction('readFile', async (params) => {
    const fs = require('fs').promises;
    const { filePath } = params;
    
    if (!filePath) {
        throw new Error('文件路径不能为空');
    }
    
    try {
        const content = await fs.readFile(filePath, 'utf8');
        return {
            success: true,
            content,
            size: content.length
        };
    } catch (error) {
        throw new Error(`读取文件失败: ${error.message}`);
    }
});

// HTTP请求示例 / HTTP request example
sdk.registerFunction('httpRequest', async (params) => {
    const https = require('https');
    const { url, method = 'GET' } = params;
    
    if (!url) {
        throw new Error('URL不能为空');
    }
    
    return new Promise((resolve, reject) => {
        const req = https.request(url, { method }, (res) => {
            let data = '';
            
            res.on('data', chunk => {
                data += chunk;
            });
            
            res.on('end', () => {
                resolve({
                    statusCode: res.statusCode,
                    headers: res.headers,
                    body: data
                });
            });
        });
        
        req.on('error', (error) => {
            reject(new Error(`HTTP请求失败: ${error.message}`));
        });
        
        req.end();
    });
});
```

### 错误处理最佳实践 / Error Handling Best Practices

```javascript
const { PluginSDK } = require('./goproc_sdk');

const sdk = new PluginSDK();

// 参数验证函数 / Parameter validation function
function validateParams(params, schema) {
    for (const [key, validator] of Object.entries(schema)) {
        const value = params[key];
        
        if (validator.required && (value === undefined || value === null)) {
            throw new Error(`缺少必需参数: ${key}`);
        }
        
        if (value !== undefined && validator.type && typeof value !== validator.type) {
            throw new Error(`参数 ${key} 类型错误，期望 ${validator.type}，实际 ${typeof value}`);
        }
        
        if (validator.validate && !validator.validate(value)) {
            throw new Error(`参数 ${key} 验证失败: ${validator.message || '无效值'}`);
        }
    }
}

// 带验证的函数示例 / Function example with validation
sdk.registerFunction('validateAndProcess', (params) => {
    try {
        // 定义参数模式 / Define parameter schema
        const schema = {
            name: {
                required: true,
                type: 'string',
                validate: (value) => value.length > 0 && value.length <= 50,
                message: '名称长度必须在1-50个字符之间'
            },
            age: {
                required: true,
                type: 'number',
                validate: (value) => value >= 0 && value <= 150,
                message: '年龄必须在0-150之间'
            },
            email: {
                required: false,
                type: 'string',
                validate: (value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
                message: '邮箱格式无效'
            }
        };
        
        // 验证参数 / Validate parameters
        validateParams(params, schema);
        
        // 处理业务逻辑 / Process business logic
        const { name, age, email } = params;
        
        return {
            success: true,
            data: {
                name: name.trim(),
                age,
                email: email || null,
                processedAt: new Date().toISOString()
            }
        };
        
    } catch (error) {
        // 统一错误处理 / Unified error handling
        console.error('处理错误:', error.message);
        throw error;
    }
});

// 带重试机制的函数 / Function with retry mechanism
sdk.registerFunction('retryableOperation', async (params) => {
    const { operation, maxRetries = 3, delay = 1000 } = params;
    
    let lastError;
    
    for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
            console.log(`尝试第 ${attempt} 次操作`);
            
            // 模拟可能失败的操作 / Simulate potentially failing operation
            if (Math.random() < 0.7) { // 70% 失败率 / 70% failure rate
                throw new Error('操作失败');
            }
            
            return {
                success: true,
                attempt,
                message: '操作成功'
            };
            
        } catch (error) {
            lastError = error;
            console.log(`第 ${attempt} 次尝试失败: ${error.message}`);
            
            if (attempt < maxRetries) {
                await new Promise(resolve => setTimeout(resolve, delay));
            }
        }
    }
    
    throw new Error(`操作失败，已重试 ${maxRetries} 次。最后错误: ${lastError.message}`);
});
```

### TypeScript 支持 / TypeScript Support

```typescript
// types.ts
interface PluginParams {
    [key: string]: any;
}

interface UserData {
    name: string;
    age: number;
    email?: string;
}

interface ProcessResult {
    success: boolean;
    data?: any;
    error?: string;
}

// plugin.ts
import { PluginSDK } from './goproc_sdk';

const sdk = new PluginSDK();

// 类型安全的函数注册 / Type-safe function registration
sdk.registerFunction('processUserTyped', (params: PluginParams): ProcessResult => {
    try {
        const userData = params.user as UserData;
        
        if (!userData || typeof userData !== 'object') {
            throw new Error('无效的用户数据');
        }
        
        const { name, age, email } = userData;
        
        if (!name || typeof name !== 'string') {
            throw new Error('用户名无效');
        }
        
        if (typeof age !== 'number' || age < 0) {
            throw new Error('年龄无效');
        }
        
        return {
            success: true,
            data: {
                name: name.trim(),
                age,
                email: email || null,
                processedAt: new Date().toISOString()
            }
        };
        
    } catch (error) {
        return {
            success: false,
            error: error instanceof Error ? error.message : '未知错误'
        };
    }
});

// 泛型函数示例 / Generic function example
function createTypedHandler<T, R>(
    validator: (params: PluginParams) => T,
    processor: (data: T) => R
) {
    return (params: PluginParams): R => {
        const validatedData = validator(params);
        return processor(validatedData);
    };
}

// 使用泛型处理器 / Using generic handler
const mathHandler = createTypedHandler(
    (params): { a: number; b: number; operation: string } => {
        const { a, b, operation } = params;
        
        if (typeof a !== 'number' || typeof b !== 'number') {
            throw new Error('a和b必须是数字');
        }
        
        if (typeof operation !== 'string') {
            throw new Error('operation必须是字符串');
        }
        
        return { a, b, operation };
    },
    ({ a, b, operation }) => {
        switch (operation) {
            case 'add': return a + b;
            case 'subtract': return a - b;
            case 'multiply': return a * b;
            case 'divide':
                if (b === 0) throw new Error('除数不能为零');
                return a / b;
            default:
                throw new Error(`不支持的操作: ${operation}`);
        }
    }
);

sdk.registerFunction('math', mathHandler);
```

## 🌍 跨平台支持 / Cross-Platform Support

SDK自动检测运行平台并使用相应的通信机制：

The SDK automatically detects the running platform and uses the appropriate communication mechanism:

- **Windows**: 使用命名管道 (Named Pipes) / Uses Named Pipes
- **Linux/macOS/FreeBSD**: 使用Unix域套接字 (Unix Domain Sockets) / Uses Unix Domain Sockets

开发者无需关心底层通信细节，SDK会自动处理平台差异。

Developers don't need to worry about underlying communication details; the SDK automatically handles platform differences.

## 🔍 调试和日志 / Debugging and Logging

### 启用调试模式 / Enable Debug Mode

```javascript
const { PluginSDK } = require('./goproc_sdk');

// 设置调试模式 / Set debug mode
process.env.DEBUG = 'goproc:*';

const sdk = new PluginSDK();

// 添加日志中间件 / Add logging middleware
const originalRegisterFunction = sdk.registerFunction.bind(sdk);
sdk.registerFunction = function(name, handler) {
    const wrappedHandler = (params) => {
        console.log(`[${new Date().toISOString()}] 调用函数: ${name}`);
        console.log(`[${new Date().toISOString()}] 参数:`, JSON.stringify(params, null, 2));
        
        try {
            const result = handler(params);
            console.log(`[${new Date().toISOString()}] 返回结果:`, JSON.stringify(result, null, 2));
            return result;
        } catch (error) {
            console.error(`[${new Date().toISOString()}] 函数执行错误:`, error.message);
            throw error;
        }
    };
    
    return originalRegisterFunction(name, wrappedHandler);
};

// 注册测试函数 / Register test function
sdk.registerFunction('debugTest', (params) => {
    return {
        message: '调试测试成功',
        timestamp: Date.now(),
        params
    };
});

sdk.start().then(() => {
    console.log('插件已启动，等待调用...');
}).catch(err => {
    console.error('插件启动失败:', err);
});
```

### 性能监控 / Performance Monitoring

```javascript
const { PluginSDK } = require('./goproc_sdk');

class PerformanceMonitor {
    constructor() {
        this.stats = new Map();
    }
    
    wrapFunction(name, handler) {
        return (params) => {
            const startTime = process.hrtime.bigint();
            const startMemory = process.memoryUsage();
            
            try {
                const result = handler(params);
                
                const endTime = process.hrtime.bigint();
                const endMemory = process.memoryUsage();
                
                this.recordStats(name, {
                    duration: Number(endTime - startTime) / 1000000, // 转换为毫秒
                    memoryDelta: endMemory.heapUsed - startMemory.heapUsed,
                    success: true
                });
                
                return result;
            } catch (error) {
                const endTime = process.hrtime.bigint();
                
                this.recordStats(name, {
                    duration: Number(endTime - startTime) / 1000000,
                    memoryDelta: 0,
                    success: false,
                    error: error.message
                });
                
                throw error;
            }
        };
    }
    
    recordStats(functionName, stats) {
        if (!this.stats.has(functionName)) {
            this.stats.set(functionName, {
                callCount: 0,
                totalDuration: 0,
                avgDuration: 0,
                maxDuration: 0,
                minDuration: Infinity,
                successCount: 0,
                errorCount: 0
            });
        }
        
        const funcStats = this.stats.get(functionName);
        funcStats.callCount++;
        funcStats.totalDuration += stats.duration;
        funcStats.avgDuration = funcStats.totalDuration / funcStats.callCount;
        funcStats.maxDuration = Math.max(funcStats.maxDuration, stats.duration);
        funcStats.minDuration = Math.min(funcStats.minDuration, stats.duration);
        
        if (stats.success) {
            funcStats.successCount++;
        } else {
            funcStats.errorCount++;
        }
        
        console.log(`[性能] ${functionName}: ${stats.duration.toFixed(2)}ms`);
    }
    
    getStats() {
        return Object.fromEntries(this.stats);
    }
    
    printStats() {
        console.log('\n=== 性能统计 ===');
        for (const [name, stats] of this.stats) {
            console.log(`函数: ${name}`);
            console.log(`  调用次数: ${stats.callCount}`);
            console.log(`  平均耗时: ${stats.avgDuration.toFixed(2)}ms`);
            console.log(`  最大耗时: ${stats.maxDuration.toFixed(2)}ms`);
            console.log(`  最小耗时: ${stats.minDuration.toFixed(2)}ms`);
            console.log(`  成功率: ${(stats.successCount / stats.callCount * 100).toFixed(2)}%`);
            console.log('');
        }
    }
}

const monitor = new PerformanceMonitor();
const sdk = new PluginSDK();

// 注册带性能监控的函数 / Register function with performance monitoring
sdk.registerFunction('monitoredFunction', monitor.wrapFunction('monitoredFunction', (params) => {
    // 模拟一些计算 / Simulate some computation
    const start = Date.now();
    while (Date.now() - start < 10) {
        // 忙等待10ms / Busy wait for 10ms
    }
    
    return { result: 'success', timestamp: Date.now() };
}));

// 定期打印统计信息 / Periodically print statistics
setInterval(() => {
    monitor.printStats();
}, 30000); // 每30秒打印一次
```

## 📝 完整示例 / Complete Example

```javascript
const { PluginSDK, registerFunction } = require('./goproc_sdk');

// 文本处理插件 / Text processing plugin
class TextProcessor {
    constructor() {
        this.sdk = new PluginSDK();
        this.setupFunctions();
    }
    
    setupFunctions() {
        // 文本反转 / Text reverse
        this.sdk.registerFunction('reverse', (params) => {
            const { text } = params;
            
            if (typeof text !== 'string') {
                throw new Error('参数 text 必须是字符串');
            }
            
            return text.split('').reverse().join('');
        });
        
        // 文本转大写 / Text to uppercase
        this.sdk.registerFunction('uppercase', (params) => {
            const { text } = params;
            
            if (typeof text !== 'string') {
                throw new Error('参数 text 必须是字符串');
            }
            
            return text.toUpperCase();
        });
        
        // 文本转小写 / Text to lowercase
        this.sdk.registerFunction('lowercase', (params) => {
            const { text } = params;
            
            if (typeof text !== 'string') {
                throw new Error('参数 text 必须是字符串');
            }
            
            return text.toLowerCase();
        });
        
        // 文本统计 / Text statistics
        this.sdk.registerFunction('stats', (params) => {
            const { text } = params;
            
            if (typeof text !== 'string') {
                throw new Error('参数 text 必须是字符串');
            }
            
            const words = text.trim().split(/\s+/).filter(word => word.length > 0);
            const lines = text.split('\n');
            const chars = text.length;
            const charsNoSpaces = text.replace(/\s/g, '').length;
            
            return {
                characters: chars,
                charactersNoSpaces: charsNoSpaces,
                words: words.length,
                lines: lines.length,
                paragraphs: text.split(/\n\s*\n/).filter(p => p.trim().length > 0).length
            };
        });
        
        // 文本格式化 / Text formatting
        this.sdk.registerFunction('format', (params) => {
            const { text, options = {} } = params;
            
            if (typeof text !== 'string') {
                throw new Error('参数 text 必须是字符串');
            }
            
            let result = text;
            
            // 移除多余空格 / Remove extra spaces
            if (options.trimSpaces !== false) {
                result = result.replace(/\s+/g, ' ').trim();
            }
            
            // 首字母大写 / Capitalize first letter
            if (options.capitalize) {
                result = result.charAt(0).toUpperCase() + result.slice(1);
            }
            
            // 移除空行 / Remove empty lines
            if (options.removeEmptyLines) {
                result = result.split('\n').filter(line => line.trim().length > 0).join('\n');
            }
            
            // 添加行号 / Add line numbers
            if (options.addLineNumbers) {
                const lines = result.split('\n');
                result = lines.map((line, index) => `${index + 1}: ${line}`).join('\n');
            }
            
            return result;
        });
        
        // 文本搜索和替换 / Text search and replace
        this.sdk.registerFunction('searchReplace', (params) => {
            const { text, search, replace, options = {} } = params;
            
            if (typeof text !== 'string' || typeof search !== 'string') {
                throw new Error('text 和 search 参数必须是字符串');
            }
            
            const replaceValue = replace || '';
            let flags = 'g'; // 全局替换
            
            if (options.ignoreCase) {
                flags += 'i';
            }
            
            if (options.multiline) {
                flags += 'm';
            }
            
            try {
                const regex = new RegExp(search, flags);
                const result = text.replace(regex, replaceValue);
                const matches = text.match(regex) || [];
                
                return {
                    result,
                    matchCount: matches.length,
                    originalLength: text.length,
                    newLength: result.length
                };
            } catch (error) {
                throw new Error(`正则表达式错误: ${error.message}`);
            }
        });
    }
    
    async start() {
        try {
            await this.sdk.start();
            console.log('文本处理插件启动成功');
            console.log('可用功能: reverse, uppercase, lowercase, stats, format, searchReplace');
        } catch (error) {
            console.error('插件启动失败:', error);
            throw error;
        }
    }
    
    stop() {
        this.sdk.stop();
        console.log('文本处理插件已停止');
    }
}

// 启动插件 / Start plugin
const processor = new TextProcessor();
processor.start().catch(console.error);

// 优雅关闭 / Graceful shutdown
process.on('SIGINT', () => {
    console.log('\n正在关闭插件...');
    processor.stop();
    process.exit(0);
});

process.on('SIGTERM', () => {
    console.log('\n正在关闭插件...');
    processor.stop();
    process.exit(0);
});
```

## 🚨 注意事项 / Important Notes

1. **异步支持**: SDK支持异步函数，可以使用`async/await`语法 / The SDK supports asynchronous functions and can use `async/await` syntax

2. **错误处理**: 使用`throw new Error()`抛出错误，SDK会自动捕获并传递 / Use `throw new Error()` to throw errors; the SDK will automatically catch and pass them

3. **参数类型**: JSON反序列化后，数字类型为JavaScript的`number`类型 / After JSON deserialization, numeric types are JavaScript `number` types

4. **内存管理**: 注意避免内存泄漏，及时清理不需要的引用 / Pay attention to avoiding memory leaks and clean up unnecessary references promptly

5. **并发处理**: SDK是线程安全的，但要注意共享状态的并发访问 / The SDK is thread-safe, but pay attention to concurrent access to shared state

6. **平台兼容**: 确保代码在不同平台上的兼容性 / Ensure code compatibility across different platforms

## 🔗 相关链接 / Related Links

- [项目主页 / Project Homepage](https://github.com/hoonfeng/goproc)
- [Go SDK 使用指南 / Go SDK Usage Guide](go_sdk_usage.md)
- [Python SDK 使用指南 / Python SDK Usage Guide](python_sdk_usage.md)
- [跨平台使用指南 / Cross-Platform Usage Guide](CROSS_PLATFORM_USAGE_GUIDE.md)