# Ent ORM 规范与迁移方案

## 1. 目标

本项目的数据访问层从当前的 `database/sql + go-sql-driver/mysql` 逐步迁移到 Ent，保持与 Simple Admin 的工程规范一致。

Simple Admin 的 Ent 生成入口为：

    make gen-ent

对应命令：

    goctls run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./rpc/ent/template/*.tmpl" ./rpc/ent/schema --feature sql/execquery,intercept,sql/modifier

官方 Ent 文档说明，Schema 用于定义实体字段、关系和索引，并通过代码生成得到类型安全的查询 API。

## 2. 当前架构

迁移前： Repository -> database/sql -> go-sql-driver/mysql -> MySQL

目标： Repository -> Generated Ent Client -> MySQL

上层 RPC/API 契约保持不变。

## 3. Schema

第一阶段已经建立：
- Drama
- Episode
- Genre
- Tag
- DramaGenre
- DramaTag
- User
- BehaviorEvent
- UserTagProfile

其中 DramaGenre、DramaTag 使用 Ent Edge Schema 表达已有的复合主键关系。

## 4. 数据库迁移策略

现阶段不执行 Ent Schema.Create 自动改库。

数据库结构仍以 `deploy/mysql/init/001_schema.sql` 和 `deploy/mysql/init/010_recommend_content.sql` 为准。Ent 负责运行时 ORM 和代码生成。

## 5. 数据库兼容原则

Ent Schema 映射现有 dramas、episodes、genres、drama_genres、tags、drama_tags、users、behavior_events、user_tag_profiles。

暂不把 `drama_feature_snapshots` 纳入第一阶段生成，因为它当前使用 `(drama_id, feature_version)` 复合主键，需要单独确认兼容策略。

`user_entitlements` 也需要先确认实际 SQL migration，再纳入 Ent，避免根据业务代码查询猜表结构。

## 6. 生成与验证

本地安装项目要求的 goctl/goctls、protoc 后运行：

    make gen-ent
    make fmt
    go test ./...

当前环境没有执行这些本地命令，因此不宣称生成代码或测试已经通过。

## 7. 下一步

1. 本地运行 `make gen-ent`，确认模板与 Schema 正常生成。
2. 将 DramaRepository 从手写 SQL 切换到 Ent。
3. 将 EpisodeRepository 切换到 Ent。
4. 将 entitlement 查询统一到 Ent。
5. 增加 Ent Client 初始化。
6. 执行完整测试。
7. 再开始 MySQL → Elasticsearch Content Indexer。

## 8. 推荐数据契约不变

Ent 只改变数据访问方式，不改变推荐数据契约。title、subtitle、description、genres、tags、country、language 和业务指标继续保留。

推荐链路仍为 Business Content → Feature Builder → Sparse/Dense/Semantic → Recall/Rank/MMoE。