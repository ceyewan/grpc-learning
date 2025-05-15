import requests
import json

url = "http://localhost:1234/jsonrpc"

payload = {
    "method": "HelloService.SayHello",
    "params": [{"Name": "Python HTTP Client"}],
    "id": 0
}
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, data=json.dumps(payload), headers=headers)
result = response.json()
message = result.get("result", {}).get("Message")
print("Message:", message)
