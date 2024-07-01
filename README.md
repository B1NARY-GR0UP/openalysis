![OPENALYSIS](./images/OPENALYSIS.png)

OPENALYSIS is a tool for visualizing and analyzing data from the GitHub open-source community.

[![Go Report Card](https://goreportcard.com/badge/github.com/B1NARY-GR0UP/openalysis)](https://goreportcard.com/report/github.com/B1NARY-GR0UP/openalysis)

## Overview

OPENALYSIS only does three things:

1. Sets up scheduled tasks to retrieve data from configured organizations or repositories via the GitHub API ([REST API](https://docs.github.com/en/rest?apiVersion=2022-11-28) and [GraphQL API](https://docs.github.com/en/graphql)).
2. Organizes and stores the retrieved data in a [MySQL](https://www.mysql.com/) database.
3. Queries the database and visualizes the data in charts and other forms through [Grafana](https://grafana.com/grafana/).

OPENALYSIS provides:

- A series of dynamic Grafana dashboards for visualizing and analyzing data from various dimensions, such as organizations, repositories, contributors, etc.
- A command-line tool to quickly start and restart scheduled tasks.
- A series of APIs to configure and use this tool.

OPENALYSIS gives you an overview of the overall data for the open-source community you manage or belong to. We hope OPENALYSIS can help you better build and develop your open-source community.

## Usage

The operation of OPENALYSIS depends on MySQL and Grafana services. The following documents will help you configure the necessary dependencies and run OPENALYSIS through the command line tool or API.

- [How to deploy?](./docs/how-to-deploy.md)
- [How to prepare config file?](./docs/how-to-prepare-config-file.md)
- [How to use CMD tool?](./docs/how-to-use-cmd-tool.md)
- [How to use API?](./docs/how-to-use-api.md)

## Grafana Dashboard Templates

OPENALYSIS offers four dimensions of Grafana Dashboards: 

- [Group](./template/OPENALYSIS-GROUP-TMPL.json)
- [Org](./template/OPENALYSIS-ORG-TMPL.json)
- [Repo](./template/OPENALYSIS-REPO-TMPL.json)
- [Contributor](./template/OPENALYSIS-CONTRIBUTOR-TMPL.json)
 
Each dashboard provides various forms of visual analysis for the corresponding dimension's data.

### Group Template

![group-tmpl](./images/tmpl-group-example.png)

### Organization Template

![org-tmpl](./images/tmpl-org-example.png)

### Repository Template

![repo-tmpl](./images/tmpl-repo-example.png)

### Contributor Template

![contributor-tmpl](./images/tmpl-contributor-example.png)

## Blogs

- [How to Visualize and Analyze Data in Open Source Communities](https://dev.to/justlorain/how-to-visualize-and-analyze-data-in-open-source-communities-1l35)

## Acknowledgement

Sincere appreciation to the [CloudWeGo](https://github.com/cloudwego) community, without whose help this project would not have been possible.

## License

OPENALYSIS is distributed under the [Apache License 2.0](./LICENSE). The licenses of third party dependencies of OPENALYSIS are explained [here](./licenses).

## ECOLOGY

<p align="center">
<img src="https://github.com/justlorain/justlorain/blob/main/images/BINARY-WEB-ECO.png" alt="BINARY-WEB-ECO"/>
<br/><br/>
OPENALYSIS is a Subproject of the <a href="https://github.com/B1NARY-GR0UP">BINARY WEB ECOLOGY</a>
</p>
