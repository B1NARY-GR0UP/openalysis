-- TODO: add index

CREATE TABLE `cursors` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`repo_node_id` longtext,`repo_name_with_owner` longtext,`last_update` datetime(3) NULL,`end_cursor` longtext,PRIMARY KEY (`id`),INDEX `idx_cursors_deleted_at` (`deleted_at`))

CREATE TABLE `contributors` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`login` longtext,`node_id` longtext,`company` longtext,`location` longtext,`avatar_url` longtext,`repo_owner` longtext,`repo_name` longtext,`repo_node_id` longtext,`contributions` bigint,PRIMARY KEY (`id`),INDEX `idx_contributors_deleted_at` (`deleted_at`))

CREATE TABLE `groups` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`name` longtext,`issue_count` bigint,`pull_request_count` bigint,`star_count` bigint,`fork_count` bigint,`contributor_count` bigint,PRIMARY KEY (`id`),INDEX `idx_groups_deleted_at` (`deleted_at`))

CREATE TABLE `issues` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`node_id` longtext,`author` longtext,`author_node_id` longtext,`repo_node_id` longtext,`repo_owner` longtext,`repo_name` longtext,`number` bigint,`state` longtext,`issue_created_at` datetime(3) NULL,`issue_closed_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX `idx_issues_deleted_at` (`deleted_at`))

CREATE TABLE `organizations` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`login` longtext,`node_id` longtext,`avatar_url` longtext,`issue_count` bigint,`pull_request_count` bigint,`star_count` bigint,`fork_count` bigint,`contributor_count` bigint,PRIMARY KEY (`id`),INDEX `idx_organizations_deleted_at` (`deleted_at`))

CREATE TABLE `pull_requests` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`node_id` longtext,`author` longtext,`author_node_id` longtext,`repo_node_id` longtext,`repo_owner` longtext,`repo_name` longtext,`number` bigint,`state` longtext,`pr_created_at` datetime(3) NULL,`pr_merged_at` datetime(3) NULL,`pr_closed_at` datetime(3) NULL,PRIMARY KEY (`id`),INDEX `idx_pull_requests_deleted_at` (`deleted_at`))

CREATE TABLE `repositories` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`owner` longtext,`name` longtext,`node_id` longtext,`owner_node_id` longtext,`issue_count` bigint,`pull_request_count` bigint,`star_count` bigint,`fork_count` bigint,`contributor_count` bigint,PRIMARY KEY (`id`),INDEX `idx_repositories_deleted_at` (`deleted_at`))

CREATE TABLE `groups_organizations` (`group_name` longtext,`org_node_id` longtext)

CREATE TABLE `groups_repositories` (`group_name` longtext,`repo_node_id` longtext)

CREATE TABLE `issue_assignees` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`issue_node_id` longtext,`issue_number` bigint,`issue_url` longtext,`issue_repo_name` longtext,`assignee_node_id` longtext,`assignee_login` longtext,PRIMARY KEY (`id`),INDEX `idx_issue_assignees_deleted_at` (`deleted_at`))

CREATE TABLE `pull_request_assignees` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`pull_request_node_id` longtext,`pull_request_number` bigint,`pull_request_url` longtext,`pull_request_repo_name` longtext,`assignee_node_id` longtext,`assignee_login` longtext,PRIMARY KEY (`id`),INDEX `idx_pull_request_assignees_deleted_at` (`deleted_at`))

CREATE TABLE `analyzed_org_contributions` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`login` longtext,`node_id` longtext,`org_login` longtext,`org_node_id` longtext,`contributions` bigint,PRIMARY KEY (`id`),INDEX `idx_analyzed_org_contributions_deleted_at` (`deleted_at`))

CREATE TABLE `analyzed_group_contributions` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`login` longtext,`node_id` longtext,`group_name` longtext,`contributions` bigint,PRIMARY KEY (`id`),INDEX `idx_analyzed_group_contributions_deleted_at` (`deleted_at`))
