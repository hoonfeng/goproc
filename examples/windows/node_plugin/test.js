// 测试Node.js插件是否能独立运行
const net = require('net');

// 模拟插件管理器连接
const client = net.createConnection({
    path: '\\\\.\\pipe\\test-node-plugin'
});

client.on('connect', () => {
    console.log('✓ 连接成功');
    
    // 发送注册消息
    const registerMsg = {
        type: 'register',
        params: { 
            functions: ['jsonStringify', 'arraySum', 'timestamp', 'uuid'] 
        }
    };
    
    const data = JSON.stringify(registerMsg);
    const length = Buffer.byteLength(data, 'utf8');
    const header = Buffer.alloc(4);
    header.writeUInt32BE(length, 0);
    
    const buffer = Buffer.concat([header, Buffer.from(data, 'utf8')]);
    client.write(buffer);
    
    console.log('✓ 注册消息已发送');
});

client.on('data', (data) => {
    console.log('收到响应:', data.toString());
});

client.on('error', (err) => {
    console.log('❌ 连接错误:', err.message);
});

client.on('close', () => {
    console.log('连接已关闭');
});

// 设置超时
setTimeout(() => {
    console.log('测试完成');
    client.end();
}, 5000);