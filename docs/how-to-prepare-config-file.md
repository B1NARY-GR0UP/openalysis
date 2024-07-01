# How to prepare config file?

## Example 

A complete configuration file example is as follows:

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
# "`previous company` => `present company`"
# "`previous location` => `present location`"
cleaner:
  - "`@CloudWeGo` => `CloudWeGo`"
  - "`@cloudwego` => `CloudWeGo`"
# "`login`, `company`, `location`"
# e.g. "`justlorain`, `binary`, ``"
# `` means do not update location
marker:
  - "`sample user`, `sample company`, `sample location`"
```

## Fields

- **groups**

Configure the organization or repository you want to analyze on a group basis.

You can set up multiple groups, where each group can contain multiple organizations or repositories.

- **datasource**

The configuration for MySQL, will be used as a datasource for Grafana.

- **backend**

The backend configuration for the service includes settings for `cron`, GitHub `token`, and `retry` attempts.

For instructions on how to obtain a GitHub token, please refer to [creating-a-fine-grained-personal-access-token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token).

- **cleaner**

The configuration for unifying contributor's location and company information.

If your community has multiple contributors from the same company or location, but the statistics are dispersed due to formatting issues, you can use cleaner to configure a series of strategies to standardize the format of the company and location in the contributors' profiles.

For example, the following strategy will standardize the company name of contributors listed as `@github` to `GitHub`:

```yaml
cleaner:
  - "`@github` => `GitHub`"
```

- **marker**

The configuration for tagging specific contributors' company and location information.

If there are some contributors in your community who haven't set their company or location in their profile, but you know their relevant information and want to include their company and location in the statistics, you can use markers to configure a set of strategies to manually set the company and location for these contributors.

For example, the following strategy will update the company and location of the user with the login `octocat` to `GitHub` and `San Francisco, USA`:

```yaml
marker:
  - "`octocat`, `GitHub`, `San Francisco, USA`"
```

