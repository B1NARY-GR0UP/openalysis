# How to use CMD tool?

## Install

```shell
go install github.com/B1NARY-GR0UP/openalysis@latest
```

## Usage

- **[Start](#start---start-openalysis-service)**: Start OPENALYSIS service
- **[Restart](#restart---restart-openalysis-service)**: Restart OPENALYSIS service

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

**NOTE: All configurations are based on the configuration file, and if flags are set, they will override the configurations in the configuration file.**

### Start - Start OPENALYSIS service

- **Usage**

```shell
openalysis start [flags] path2config.yaml
```

- **Flags**

| Short | Long    | Description                                                                                                                                                                            |
|-------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -t    | --token | [Your GitHub Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) |
| -c    | --cron  | Your Cron Spec                                                                                                                                                                         |
| -r    | --retry | Retry Times                                                                                                                                                                            |
| -h    | --help  | Help for Start                                                                                                                                                                         |

- **Example**

```shell
openalysis start -c "@hourly" -r "5" config.yaml
```

### Restart - Restart OPENALYSIS service

- **Usage**

```shell
openalysis restart [flags] path2config.yaml
```

- **Flags**

| Short | Long    | Description                                                                                                                                                                            |
|-------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -t    | --token | [Your GitHub Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token) |
| -c    | --cron  | Your Cron Spec                                                                                                                                                                         |
| -r    | --retry | Retry Times                                                                                                                                                                            |
| -h    | --help  | Help for Restart                                                                                                                                                                       |

- **Example**

```shell
openalysis restart -t "example-github-token" config.yaml
```
