# 如何准备配置文件？

## 示例

一个完整的配置文件示例如下：

```yaml
groups:
  - name: "cloudwego"
    orgs:
      - "cloudwego"
      - "kitex-contrib"
      - "hertz-contrib"
      - "volo-rs"
    repos:
      - "bytedance/sonic"
      - "bytedance/monoio"
datasource:
  mysql:
    host: "mysql"
    port: "3306"
    user: "openalysis"
    password: "openalysis"
    database: "openalysis?charset=utf8&parseTime=True&loc=Local"
backend:
  cron: "@daily"
  token: "your-github-token"
  retry: 3
cleaner:
  - "`@CloudWeGo` => `CloudWeGo`"
  - "`@cloudwego` => `CloudWeGo`"
marker:
  - "`sample user`, `sample company`, `sample location`"
```

## 字段说明

- **groups**

按组设置您想要分析的组织或仓库。

您可以设置多个组，每个组可以包含多个组织或仓库。

- **datasource**

MySQL 的配置，将作为 Grafana 的数据源使用。

- **backend**

服务的后端配置，包括 `cron` 定时任务设置，GitHub 的 `token`，和重试尝试次数。

有关如何获取 GitHub token 的说明，请参阅 [creating-a-fine-grained-personal-access-token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)。

- **cleaner**

用于统一贡献者位置和公司信息的配置。

如果您的社区中有多个来自同一公司或地点的贡献者，但由于格式问题导致统计数据分散，您可以使用 cleaner 配置一系列策略来标准化贡献者档案中公司和位置的格式。

例如，以下策略将会将以 `@github` 列出的贡献者的公司名称标准化为 `GitHub`：

```yaml
cleaner:
  - "`@github` => `GitHub`"
```

- **marker**

用于标记特定贡献者公司和位置信息的配置。

如果您社区中有一些贡献者没有在其档案中设置公司或位置信息，但您知道他们的相关信息并希望包含在统计数据中，您可以使用 marker 配置一组策略，手动设置这些贡献者的公司和位置。

例如，以下策略将会将登录名为 `octocat` 的用户的公司和位置更新为 `GitHub` 和 `美国加利福尼亚州旧金山`：

```yaml
marker:
  - "`octocat`, `GitHub`, `San Francisco, USA`"
```