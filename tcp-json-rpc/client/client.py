import json
import socket

class JSONRPCClient:
    def __init__(self, host='localhost', port=1234):
        self.host = host
        self.port = port
        self.socket = None

    def connect(self):
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.connect((self.host, self.port))

    def close(self):
        if self.socket:
            self.socket.close()
            self.socket = None

    def call(self, method, params):
        if not self.socket:
            self.connect()

        # 创建 JSON-RPC 请求
        request = {
            "id": 0,
            "method": method,
            "params": [params]
        }

        # 发送请求
        request_json = json.dumps(request) + "\n"  # 添加换行符作为消息分隔符
        self.socket.sendall(request_json.encode('utf-8'))

        # 接收响应
        response = self.socket.recv(4096).decode('utf-8')
        return json.loads(response)

# 使用示例
if __name__ == "__main__":
    client = JSONRPCClient()
    try:
        # 调用 HelloService.Hello 方法
        response = client.call("HelloService.SayHello", {"Name": "Python TCP Client"})
        message = response['result']['Message']
        print("服务器响应:", message)
    finally:
        client.close()
