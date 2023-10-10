# 1. Go Tools

这是一个包含多个独立 Go 工具的项目

## 1.1. 工具列表

- [go-ssh](./go-ssh): 一个用于处理 SSH 连接的工具
- [godaddy-ddns](./godaddy-ddns): 一个用于更新 GoDaddy 动态 DNS 记录的工具
- [gost-subscribe](./gost-subscribe): 一个用于订阅 Gost 代理服务器的工具

## 1.2. 使用说明

每个工具都包含一个独立的 `main.go` 文件，用于定义和执行该工具的功能您可以按照以下步骤来运行每个工具：

1. 进入工具目录：
   ```
   cd <tool-directory>
   ```

2. 运行工具：
   ```
   go run main.go
   ```

请注意，您需要先安装 Golang 开发环境，并确保您的环境已正确配置

## 1.3. 依赖项

此项目使用 Go Modules 进行依赖项管理项目根目录下的 `go.mod` 文件列出了所有依赖项及其版本在运行每个工具之前，您可以使用以下命令来下载和安装项目的依赖项：

```
go mod download
```

## 1.4. 开发环境
采用`vscode`进行开发
- Go Version: go version go1.21.2 linux/amd64

`.vscode/launch.json`参考如下
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Go SSH",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/ssh-config/main.go",
            "args": [
                "-c=${workspaceFolder}/data/config"
            ],
            "console": "integratedTerminal"
        },
        {
            "name": "Go Subscribe",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/gost-subscribe/main.go",
            "console": "internalConsole",
            "args": [
                "-u=https://example.com/api/v1/client/subscribe?token=mytoken&flag=shadowsocks",
                "-l=300"
            ]
        },
        {
            "name": "Go Godaddy",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/godaddy-ddns/main.go",
            "args": ["--domain=chancel.me","--type=A","--name=test","--shopperid=12345678","--key=7st2btydfasdch63fd","--secret=asdg763vvdafsd","--proxy=http://172.16.21.10:11080"]
        }
    ]
}
```


## 1.5. 贡献

欢迎对该项目进行贡献！如果您发现问题或有改进建议，请提出新的 Issue 或提交 Pull Request