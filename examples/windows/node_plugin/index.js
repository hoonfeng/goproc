const { PluginSDK } = require('../../../sdk/javascript/goproc_sdk.js');

// 创建插件SDK实例
const sdk = new PluginSDK();

// 注册jsonStringify函数 - JSON序列化
sdk.registerFunction('jsonStringify', (params) => {
    try {
        // 处理测试中传入的对象格式 {obj: {...}}
        const obj = params.obj || params;
        return JSON.stringify(obj);
    } catch (error) {
        throw new Error(`JSON序列化失败: ${error.message}`);
    }
});

// 注册jsonParse函数 - JSON解析
sdk.registerFunction('jsonParse', (params) => {
    try {
        const jsonStr = params.json || params;
        return JSON.parse(jsonStr);
    } catch (error) {
        throw new Error(`JSON解析失败: ${error.message}`);
    }
});

// 注册arraySum函数 - 数组求和
sdk.registerFunction('arraySum', (params) => {
    try {
        // 处理测试中传入的数组格式 {array: [...]}
        const array = params.array || params;
        if (!Array.isArray(array)) {
            throw new Error('参数必须是数组');
        }
        return array.reduce((sum, num) => {
            const value = parseFloat(num);
            if (isNaN(value)) {
                throw new Error(`无效的数字: ${num}`);
            }
            return sum + value;
        }, 0);
    } catch (error) {
        throw new Error(`数组求和失败: ${error.message}`);
    }
});

// 注册arrayFilter函数 - 数组过滤
sdk.registerFunction('arrayFilter', (params) => {
    try {
        const array = params.array || [];
        const condition = params.condition || 'positive';
        
        if (!Array.isArray(array)) {
            throw new Error('参数必须是数组');
        }
        
        switch (condition) {
            case 'positive':
                return array.filter(num => parseFloat(num) > 0);
            case 'negative':
                return array.filter(num => parseFloat(num) < 0);
            case 'even':
                return array.filter(num => parseFloat(num) % 2 === 0);
            case 'odd':
                return array.filter(num => parseFloat(num) % 2 !== 0);
            default:
                return array;
        }
    } catch (error) {
        throw new Error(`数组过滤失败: ${error.message}`);
    }
});

// 注册timestamp函数 - 获取当前时间戳
sdk.registerFunction('timestamp', (params) => {
    try {
        const format = params?.format || 'unix';
        const now = new Date();
        
        switch (format) {
            case 'unix':
                return Math.floor(now.getTime() / 1000);
            case 'milliseconds':
                return now.getTime();
            case 'iso':
                return now.toISOString();
            case 'readable':
                return now.toLocaleString();
            default:
                return Math.floor(now.getTime() / 1000);
        }
    } catch (error) {
        throw new Error(`获取时间戳失败: ${error.message}`);
    }
});

// 注册uuid函数 - 生成UUID
sdk.registerFunction('uuid', (params) => {
    try {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    } catch (error) {
        throw new Error(`生成UUID失败: ${error.message}`);
    }
});

// 注册httpGet函数 - 简单的URL信息返回（同步版本）
sdk.registerFunction('httpGet', (params) => {
    try {
        const url = params.url;
        if (!url) {
            throw new Error('缺少URL参数');
        }
        
        // 返回URL信息而不是实际请求（避免异步复杂性）
        const urlObj = new URL(url);
        return {
            url: url,
            protocol: urlObj.protocol,
            hostname: urlObj.hostname,
            port: urlObj.port,
            pathname: urlObj.pathname,
            search: urlObj.search,
            message: '这是一个简化的HTTP GET函数，返回URL信息'
        };
    } catch (error) {
        throw new Error(`URL解析失败: ${error.message}`);
    }
});

// 异步启动插件
async function main() {
    try {
        
        
        // 启动插件（等待连接建立）
        const success = await sdk.start();
        
        if (success) {
            
            
            // 保持进程运行，等待消息
            process.on('SIGINT', () => {
                
                sdk.stop();
                process.exit(0);
            });
            
            process.on('SIGTERM', () => {
                
                sdk.stop();
                process.exit(0);
            });
            
            // 保持进程运行
           
        } else {
            
            process.exit(1);
        }
    } catch (error) {
        
        process.exit(1);
    }
}

// 启动插件
main();