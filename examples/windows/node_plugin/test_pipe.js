const net = require('net');

// 测试Windows命名管道
const pipePath = '\\\\.\\pipe\\test-node-pipe-123';

console.log('尝试创建命名管道服务器:', pipePath);

const server = net.createServer((socket) => {
    console.log('客户端已连接');
    
    socket.on('data', (data) => {
        console.log('收到数据:', data.toString());
        socket.write('Hello from server!');
    });
    
    socket.on('end', () => {
        console.log('客户端断开连接');
    });
});

server.listen(pipePath, () => {
    console.log('命名管道服务器正在监听:', pipePath);
});

server.on('error', (error) => {
    console.error('服务器错误:', error.message);
});

// 保持服务器运行
console.log('服务器启动完成，等待连接...');