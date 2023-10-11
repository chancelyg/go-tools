# 1. Go Tools

这是一个包含多个独立 Go 可执行程序的工具类项目

每个文件夹都包含一个独立的 `main.go` 文件，可以参考根目录下的`goreleaser.example.yaml`单独打包需要的工具程序

## 1.1. 工具说明
### 1.1.1. go-ssh
该程序用于将`ssh/config`中的ssh服务器以列表的形式在shell中呈现

使用效果如下
```bash
❯ ./go-ssh
************************ Hi, Welcome to use Go-SSH Tool *****************************

+-----+------------------------------+-------------------------+------------------------------------------+
| id  | Host                         | username                | address                                  |
+-----+------------------------------+-------------------------+------------------------------------------+
| 1   | 192.168.122.4                | chancel                 | 192.168.122.4                            |
| 2   | 192.168.4.15                 | chancel                 | 192.168.4.15                             |
| 3   | 192.168.11.2                 | chancel                 | 192.168.11.2                             |
| 4   | 192.168.11.3                 | chancel                 | 192.168.11.3                             |
| 5   | 192.168.11.12                | chancel                 | 192.168.11.12                            |
| 6   | 192.168.4.10                 | chancel                 | 192.168.4.10                             |
+-----+------------------------------+-------------------------+------------------------------------------+

Tips: Press a number between 1 and 5 to select the host to connect, or "q" to quit.

# 

```

你可以使用一个数字选择要连接的主机，然后会自动使用SSH连接到该主机

### 1.1.2. gost-subscribe

该程序用于生成 gost 配置文件，根据给定的SS订阅链接，随机选择服务器并生成[Gost](https://gost.run)配置文件。

可用的命令行选项：

- `-h`：显示帮助信息
- `-u`：订阅链接（默认为 http://www.subcriptionurl.com）
- `-l`：要获取的服务器数量（默认为 10）
- `-p`：TCP 代理端口号（默认为 11080）
- `-r`：RED 代理端口号（默认为 11081）
- `-s`：代理策略（round|rand|fifo|hash，默认为 fifo）
- `-t`：代理失败超时时间（以秒为单位，默认为 600）
- `-m`：代理最大失败次数（默认为 1）
- `-V`：显示版本信息
- `-f`：过滤包含关键字的订阅（默认为 "套餐|重置|剩余|更新"）
- `-o`：输出文件路径（默认为 config.yml）

以下示例演示如何使用此程序生成配置文件：
```shell
./gost-config-generator -u http://www.example.com/subscription -o config.yml
```

### 1.1.3. godaddy-ddns

该程序可以使用Godaddy API来更新域名的DNS记录（即动态DDNS）

原理：获取运行设备的外网地址（IPV4/IPV6）并更新指定域名的解析结果

可用的命令行选项：
- `--domain`：域名（如chancel.me）
- `--type`：域名类型，默认为 A 类型（A类型是IPV4地址，AAAA类型是IPV6地址）
- `--name`：子域名，多个子域名以逗号分隔
- `--record`：域名对应的值，如果为空，则自动获取 IPV4/IPV6 值
- `--shopperid`：Godaddy API 的 shopper id
- `--key`：Godaddy API 的 key
- `--secret`：Godaddy API 的 secret
- `--proxy`：HTTP 请求的代理

使用如下
```bash
# 更新单个子域名
godaddy-ddns --domain=chancel.me --name=test --shopperid=123456789 --key=wm2uSn3udbsHu_s38ndy1hdj8 --secret=Abweid78b3nx9nHSDg3gh

# 更新多个子域名
godaddy-ddns --domain=chancel.me --name=test1,test2,test3,test4 --shopperid=123456789 --key=wm2uSn3udbsHu_s38ndy1hdj8 --secret=Abweid78b3nx9nHSDg3gh

# 指定更新的ip值
godaddy-ddns --domain=chancel.me --name=test --record=192.168.1.1 --shopperid=123456789 --key=wm2uSn3udbsHu_s38ndy1hdj8 --secret=Abweid78b3nx9nHSDg3gh
```

## 1.2. 依赖项

此项目使用 Go Modules 进行依赖项管理项目根目录下的 `go.mod` 文件列出了所有依赖项及其版本在运行每个工具之前，您可以使用以下命令来下载和安装项目的依赖项：

```
go mod tidy
```

## 1.3. 开发环境
采用`vscode`进行开发
- Go Version: go version go1.21.2 linux/amd64

环境配置信息`.vscode/launch.json`参考
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


## 1.4. 贡献

欢迎对该项目进行贡献！如果您发现问题或有改进建议，请提出新的 Issue 或提交 Pull Request