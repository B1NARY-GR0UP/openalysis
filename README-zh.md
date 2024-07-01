![OPENALYSIS](./images/OPENALYSIS.png)

OPENALYSIS 是一款对 GitHub 开源社区的数据进行可视化和分析的工具。

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/openalysis)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/openalysis)

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


## Grafana Dashboard Templates

OPENALYSIS offers four dimensions of Grafana Dashboards:

- [Group](./template/OPENALYSIS-GROUP-TMPL.json)
- [Org](./template/OPENALYSIS-ORG-TMPL.json)
- [Repo](./template/OPENALYSIS-REPO-TMPL.json)
- [Contributor](./template/OPENALYSIS-CONTRIBUTOR-TMPL.json)

Each dashboard provides various forms of visual analysis for the corresponding dimension's data.

### Group Template

![group-tmpl](./images/tmpl-group-example.png)

In the Group Template, the panels are arranged in the following order from left to right, top to bottom:

- **Star Count:** Group 中所有 Organizations 和 Repositories 的 stargazer 数总和。
- **Contributor Count:** Group 中所有 Organizations 和 Repositories 的 contributor 数总和，统计时对重复的贡献者进行了去重。
- **Issue Count:** Group 中所有 Organizations 和 Repositories 的 issue 数总和。
- **PR Count:** Group 中所有 Organizations 和 Repositories 的 pull request 数总和。
- **Fork Count:** Group 中所有 Organizations 和 Repositories 的 fork 数总和。
- **Star Count:** Group 中所有 Organizations 和 Repositories 的 stargazer 数变化趋势图。
- **Contributor Count:** Group 中所有 Organizations 和 Repositories 的 contributor 数变化趋势图，统计时对重复的贡献者进行了去重。
- **Issue Count:** Group 中所有 Organizations 和 Repositories 的 issue 数变化趋势图。
- **PR Count:** Group 中所有 Organizations 和 Repositories 的 pull request 数变化趋势图。
- **Fork Count:** Group 中所有 Organizations 和 Repositories 的 fork 数变化趋势图。
- **Contributor Company:** Group 中所有 Organizations 和 Repositories 的 contributor 的 company 分布饼图。
- **Contributor Location:** Group 中所有 Organizations 和 Repositories 的 contributor 的 location 分布饼图。
- **Leaderboard:** Group 中所有 Organizations 和 Repositories 的 contributor 的贡献数排名，其中 `Ranged Contributions` 字段会对 Grafana Dashboard 设置的时间范围内的贡献做差值统计。
- **Issue Assignees:** Group 中所有 Organizations 和 Repositories 的被分配有 Assignee 的 OPEN issue。
- **PR Assignees:** Group 中所有 Organizations 和 Repositories 的被分配有 Assignee 的 OPEN pull requests。

### Organization Template

![org-tmpl](./images/tmpl-org-example.png)

In the Organization Template, the panels are arranged in the following order from left to right, top to bottom:

### Repository Template

![repo-tmpl](./images/tmpl-repo-example.png)

In the Repository Template, the panels are arranged in the following order from left to right, top to bottom:

### Contributor Template

![contributor-tmpl](./images/tmpl-contributor-example.png)

In the Contributor Template, the panels are arranged in the following order from left to right, top to bottom:

## 博客

- [如何对开源社区的数据进行可视化分析](https://juejin.cn/post/7359882185362948135)

## 致谢

Sincere appreciation to the [CloudWeGo](https://github.com/cloudwego) community, without whose help this project would not have been possible.

## 许可证

OPENALYSIS is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of OPENALYSIS are explained [here](./licenses).

## 生态

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BINARY-WEB-ECO.png" alt="BINARY-WEB-ECO"/>
<br/><br/>
OPENALYSIS is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">BINARY WEB ECOLOGY</a>
</p>