# 如何使用 CMD 工具？

## 安装

```shell
go install github.com/B1NARY-GR0UP/openalysis@latest
```

## 使用方法

- **[启动](#start---start-openalysis-service)**: 启动 OPENALYSIS 服务
- **[重启](#restart---restart-openalysis-service)**: 重启 OPENALYSIS 服务

```shell
Usage:        
  openalysis [command] path2config.yaml

Available Commands:
  help        Help about any command
  restart     restart openalysis service
  start       start openalysis service

Flags:
  -h, --help      help for openalysis
  -v, --version   version for openalysis
```

**注意：所有配置基于配置文件，如果设置了标志，它们将覆盖配置文件中的配置。**

### 启动 - 启动 OPENALYSIS 服务

- **用法**

```shell
openalysis start [flags] path2config.yaml
```

- **标志**

| 短标志 | 长标志     | 描述                                                                                                                                                                                |
|-----|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -t  | --token | [您的 GitHub 令牌](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) |
| -c  | --cron  | 您的 Cron 表达式                                                                                                                                                                       |
| -r  | --retry | 重试次数                                                                                                                                                                              |
| -h  | --help  | 启动命令的帮助信息                                                                                                                                                                         |

- **示例**

```shell
openalysis start -c "@hourly" -r "5" config.yaml
```

### 重启 - 重启 OPENALYSIS 服务

- **用法**

```shell
openalysis restart [flags] path2config.yaml
```

- **标志**

| 短标志 | 长标志     | 描述                                                                                                                                                                                |
|-----|---------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -t  | --token | [您的 GitHub 令牌](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) |
| -c  | --cron  | 您的 Cron 表达式                                                                                                                                                                       |
| -r  | --retry | 重试次数                                                                                                                                                                              |
| -h  | --help  | 重启命令的帮助信息                                                                                                                                                                         |

- **示例**

```shell
openalysis restart -t "example-github-token" config.yaml
```