![OPENALYSIS](./images/OPENALYSIS.png)

OPENALYSIS 是一款对 GitHub 开源社区的数据进行可视化和分析的工具。

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/openalysis)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/openalysis)

[View English Document](./README.md)

## 概览

OPENALYSIS 只做了三件事：

1. 设置定时任务通过 GitHub API ([REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) 和 [GraphQL API](https://docs.github.com/en/graphql)) 获取配置的组织或者仓库的数据；
2. 将获取的数据整理并存储在 [MySQL](https://www.mysql.com/) 数据库中；
3. 通过 [Grafana](https://grafana.com/grafana/) 查询数据库并将数据可视化为图表等形式；

OPENALYSIS 提供了：

- 一系列动态的 Grafana Dashboard 来对配置的组织，仓库，贡献者等不同维度的数据进行可视化和分析；
- 一个命令行工具来快速启动和重启定时任务；
- 一系列 API 来配置和使用这个工具；

OPENALYSIS 可以让您对您所管理的或所在的开源社区的整体数据有一个概览，我们希望 OPENALYSIS 可以帮助您更好的对开源社区进行建设和发展。

## 使用

OPENALYSIS 的运行依赖于 MySQL 和 Grafana 服务。以下文档将帮助您配置必要的依赖项，并通过命令行工具或 API 运行 OPENALYSIS。

- [如何部署？](./docs/how-to-deploy-zh.md)
- [如何准备配置文件？](./docs/how-to-prepare-config-file-zh.md)
- [如何使用命令行工具？](./docs/how-to-use-cmd-tool-zh.md)
- [如何使用 API？](./docs/how-to-use-api-zh.md)

## Grafana 仪表盘模板

OPENALYSIS 提供了四个维度的 Grafana 仪表板：

- [Group](./template/OPENALYSIS-GROUP-TMPL.json)
- [Org](./template/OPENALYSIS-ORG-TMPL.json)
- [Repo](./template/OPENALYSIS-REPO-TMPL.json)
- [Contributor](./template/OPENALYSIS-CONTRIBUTOR-TMPL.json)

每个仪表板都为相应维度的数据提供各种形式的可视化分析。

### Group 模板

![group-tmpl](./images/tmpl-group-example.png)

在 Group 模板中，面板按以下顺序从左到右、从上到下排列：

- **Star Count:** 组内所有组织和仓库的总标星数。
- **Contributor Count:** 组内所有组织和仓库的总贡献者数，计算时会去重重复的贡献者。
- **Issue Count:** 组内所有组织和仓库的总问题数。
- **PR Count:** 组内所有组织和仓库的总拉取请求数。
- **Fork Count:** 组内所有组织和仓库的总分叉数。
- **Star Count:** 显示组内所有组织和仓库的标星数变化趋势图。
- **Contributor Count:** 显示组内所有组织和仓库的贡献者数变化趋势图，计算时会去重重复的贡献者。
- **Issue Count:** 显示组内所有组织和仓库的问题数变化趋势图。
- **PR Count:** 显示组内所有组织和仓库的拉取请求数变化趋势图。
- **Fork Count:** 显示组内所有组织和仓库的分叉数变化趋势图。
- **Contributor Company:** 显示组内所有组织和仓库的贡献者公司分布的饼图。
- **Contributor Location:** 显示组内所有组织和仓库的贡献者地点分布的饼图。
- **Leaderboard:** 按贡献数量对组内所有组织和仓库的贡献者进行排名，`Ranged Contributions` 字段计算在 Grafana 仪表板中设定的时间范围内的贡献差异。
- **Issue Assignees:** 显示组内所有组织和仓库中具有负责人且处于 OPEN 状态的问题。
- **PR Assignees:** 显示组内所有组织和仓库中具有负责人且处于 OPEN 状态的拉取请求。

### Organization 模板

![org-tmpl](./images/tmpl-org-example.png)

在 Organization 模板中，面板按以下顺序从左到右、从上到下排列：

- **Profile:** 组织头像。
- **Contributor Company (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的贡献者公司分布的饼图。
- **Contributor Location (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的贡献者地点分布的饼图。
- **Leaderboard:** 按贡献数量对组织下所有仓库的贡献者进行排名，`Ranged Contributions` 字段计算在 Grafana 仪表板中设定的时间范围内的贡献差异。
- **Star Count:** 组织下所有仓库的总标星数。
- **Contributor Count:** 组织下所有仓库的总贡献者数，计算时会去重重复的贡献者。
- **Issue Count:** 组织下所有仓库的总问题数。
- **PR Count:** 组织下所有仓库的总拉取请求数。
- **Fork Count:** 组织下所有仓库的总分叉数。
- **Star Count:** 显示组织下所有仓库标星数变化的趋势图。
- **Contributor Count:** 显示组织下所有仓库贡献者数变化的趋势图，计算时会去重重复的贡献者。
- **Issue Count:** 显示组织下所有仓库问题数变化的趋势图。
- **PR Count:** 显示组织下所有仓库拉取请求数变化的趋势图。
- **Fork Count:** 显示组织下所有仓库分叉数变化的趋势图。
- **Issue Assignees (support repos):** 显示由 `repos` 变量指定的一个或多个仓库中具有负责人且处于OPEN状态的问题。
- **PR Assignees (support repos):** 显示由 `repos` 变量指定的一个或多个仓库中具有负责人且处于OPEN状态的拉取请求。
- **Star Count (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的标星数变化的趋势图。
- **Contributor Count (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的贡献者数变化的趋势图，计算时会去重重复的贡献者。
- **Fork Count (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的分叉数变化的趋势图。
- **Issue Count (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的问题数变化的趋势图。
- **PR Count (support repos):** 显示由 `repos` 变量指定的一个或多个仓库的拉取请求数变化的趋势图。

### Repository 模板

![repo-tmpl](./images/tmpl-repo-example.png)

在 Repository 模板中，面板按以下顺序从左到右、从上到下排列：

- **Contributor Company:** 显示仓库中所有贡献者公司分布的饼图。
- **Contributor Location:** 显示仓库中所有贡献者地点分布的饼图。
- **Leaderboard:** 根据贡献数量对仓库中所有贡献者进行排名，`Ranged Contributions` 字段计算在 Grafana 仪表板中设定的时间范围内的贡献差异。
- **Star Count:** 仓库中的总标星数。
- **Contributor Count:** 仓库中的总贡献者数，计算时会去重重复的贡献者。
- **Issue Count:** 仓库中的总问题数。
- **PR Count:** 仓库中的总拉取请求数。
- **Fork Count:** 仓库中的总分叉数。
- **Star Count:** 显示仓库中标星数变化的趋势图。
- **Contributor Count:** 显示仓库中贡献者数变化的趋势图，计算时会去重重复的贡献者。
- **Issue Count:** 显示仓库中问题数变化的趋势图。
- **PR Count:** 显示仓库中拉取请求数变化的趋势图。
- **Fork Count:** 显示仓库中分叉数变化的趋势图。
- **Issue Assignees:** 显示已分配负责人且处于OPEN状态的仓库中的所有问题。
- **PR Assignees:** 显示已分配负责人且处于OPEN状态的仓库中的所有拉取请求。

### Contributor 模板

![contributor-tmpl](./images/tmpl-contributor-example.png)

在 Contributor 模板中，面板按以下顺序从左到右、从上到下排列：

- **Profile:** 贡献者的头像和其他信息。
- **PR History:** 贡献者的拉取请求历史。
- **Contributions:** 贡献者对每个仓库的贡献统计。
- **Assigned Issues:** 分配给贡献者的问题。
- **Assigned PRs:** 分配给贡献者的拉取请求。
- **Issue Count:** 贡献者创建的问题统计。
- **PR Count:** 贡献者创建的拉取请求统计。

## 博客

- [如何对开源社区的数据进行可视化分析](https://juejin.cn/post/7359882185362948135)

## 致谢

真诚感谢 [CloudWeGo](https://github.com/cloudwego) 社区的帮助，没有他们的支持，这个项目将无法实现。

## 许可证

OPENALYSIS 使用 [Apache License 2.0](./LICENSE) 进行分发。OPENALYSIS 的第三方依赖项的许可证说明在[此处](./licenses)。

## 生态

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BINARY-WEB-ECO.png" alt="BINARY-WEB-ECO"/>
<br/><br/>
OPENALYSIS 是 <a href="https://github.com/B1NARY-GR0UP"> BINARY 网络生态 </a> 的一个子项目
</p>