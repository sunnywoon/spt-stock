# SpartanHost VPS 库存监控器

这是一个命令行工具，用于监控 SpartanHost VPS（虚拟专用服务器）的库存情况。它会定期检查 VPS 的可用性，并在检测到变化时提供通知。

## 功能

- 监控 SpartanHost VPS 的库存情况
- 可自定义检测间隔时间
- 支持设置请求代理
- 支持通过设置环境变量传入 Webhook URL

## 使用方法
该工具支持以下命令参数：

-sleep：指定检测间隔时间（秒）（默认值：60）
-proxy：设置请求代理地址（可选）

### 示例：

使用 30 秒的检测间隔运行监控器：
```
./spt-stock -sleep 30
```

使用代理运行监控器
```
./spt-stock -proxy http://proxy-server:port
```

## webhook通知
您还可以通过设置环境变量来传递 Webhook URL：
在项目根目录下创建一个名为 .env 的文件，并添加以下内容：
```
WEBHOOK_URL=<Your Webhook URL>
```
将 <Your Webhook URL> 替换为您的实际 Webhook URL。

当库存不为 0 时，工具会向指定的 Webhook_URL Post一个 JSON 消息。消息的示例格式如下：

```
{
  "e3_list": [
    {
      "name": "512MB SEAHKVM", 
      "num": 1, // 库存
      "price": "$3.60 USD", // 每月价格
      "link": "/store/ddos-protected-hdd-e3-kvm-vps-seattle/512mb-seahkvm" // 下单地址
    },
    ...
  ],
  "e5_list": [
    ...
  ],
  "amd_list": [
    ...
  ]
}
```

