# 开发日志

## 2026-09-19 — V1.1 内容数据与推荐特征契约

### 本次检查

检查仓库：`jackjie20162/short-drama-recommend`

当前已有：

- Go-zero + gRPC 服务结构
- MySQL 用户、短剧、剧集、标签、行为表
- drama-api / drama-rpc / behavior-rpc / recommend-rpc
- 推荐 V1 规则排序
- Docker/本地测试脚本
- 已有架构、推荐、支付、媒体等项目文档

### 本次决策

将内容数据正式分为：

```text
Business Entity
    |
    +--> Search / Recall Document
    |
    +--> RecommendationFeature
              |
              +--> sparse
              +--> dense
              +--> semantic
```

明确 title / subtitle / description / genres / tags / country / language 属于业务语义层；MMoE 不直接消费原始文本，而由 Feature Builder 生成 embedding 等模型特征。

### 本次落地

- `docs/12-content-data-contract.md`
- `docs/13-elasticsearch-index.md`
- `docs/14-development-log.md`
- `deploy/mysql/init/010_recommend_content.sql`
- 扩展 `proto/drama.proto`
- 扩展 `proto/behavior.proto`
- 扩展 `internal/model/drama.go`
- 扩展 `internal/model/episode.go`
- 扩展 `internal/model/behavior.go`
- 新增 `internal/model/content_feature.go`
- 新增 `internal/model/feature_builder.go`
- 更新 `docs/01-architecture.md`
- 更新 `docs/11-recommendation.md`

### 关键数据结构

业务层：

```text
Drama
  title
  subtitle
  description
  genres[]
  tags[]
  country
  language
  popularity
  completion_rate
  pay_rate
```

模型层：

```text
RecommendationFeature
  sparse
  dense
  semantic
```

### Git 提交记录

本次操作已逐步提交到 `main`，主要提交：

- `813c05e` — recommendation content schema
- `5f18775` — drama content contract
- `76b3eac` — behavior event context
- `b0ffbe6` / `91af4b4` — recommendation feature contracts
- `571a63d` — drama business model
- `0404f65` — episode content model
- `bd84f97` — behavior event model
- `a0e723f` / `3b55c05` — architecture/recommendation docs

### 验证状态

已完成 GitHub 代码结构与契约检查。

当前环境没有直接执行用户本地 Go/Protobuf/Docker 的能力，因此本日志不把 `go test ./...`、`protoc`、Docker 全链路测试标记为已通过。

本次修改保持 Proto 字段向后兼容：仅新增字段，不复用既有字段编号。

### 下一步

1. 让 `drama-rpc` Repository 真正读取 genres/tags/metrics/subtitle/发布时间。
2. 实现 Content Indexer：MySQL -> Elasticsearch。
3. 实现 Feature Builder 的 user + drama 双侧输入。
4. 建立 Redis country × language 热度榜。
5. 生成 MMoE 离线训练样本。
6. 接入训练/评估流程，再进入 rank-rpc。
7. 每一步继续追加本文件，记录代码、数据库、验证结果和 Git commit。
\n\n## 2026-09-19 — V1.2 Drama Repository 内容字段落地\n\n### 本次落地\n\n- `internal/repository/drama.go`：GetByID/ListPublished 读取 subtitle、popularity、completion_rate、pay_rate、published_at、genres、tags。\n- `rpc/drama-rpc/internal/logic/convert.go`：gRPC Drama 响应输出新增内容语义字段与指标；剧集响应输出 description。\n- genres/tags 使用聚合查询，业务库仍作为 source of truth。\n\n### 兼容处理\n\n- published_at 使用 NULL-safe 扫描。\n- 原有 Proto 字段编号不变，仅使用已新增字段。\n\n### 验证状态\n\n- 已完成 GitHub 静态代码检查。\n- 尚未声称用户本地 `go test ./...`、`protoc` 或 Docker E2E 已通过。\n\n### 下一步\n\n1. 建立 Content Indexer：MySQL Drama -> Elasticsearch。\n2. 增加 ES bulk/upsert/delete 与增量同步模型。\n3. 接入 Feature Builder，形成用户侧 + 内容侧特征。\n

## V1.3 — Simple Admin Ent ORM 基础设施

- 按用户提供的 Simple Admin Makefile 对齐 Ent 生成规范。
- 新增 `rpc/ent/schema/`：Drama、Episode、Genre、Tag、DramaGenre、DramaTag、User、BehaviorEvent、UserTagProfile。
- 复制 Simple Admin Ent templates：`pagination.tmpl`、`set_not_nil.tmpl`。
- Makefile 新增 `gen-ent`，使用 `goctls run -mod=mod entgo.io/ent/cmd/ent generate`。
- go.mod 增加 Ent ORM 依赖。
- 第一阶段不启用 Ent 自动迁移，继续以 deploy/mysql/init SQL 为数据库结构来源。
- 当前环境没有执行 `make gen-ent` 或 `go test ./...`，因此不声称生成代码或测试已经通过。
- 下一步：本地生成 Ent client 后，把 Drama/Episode/Entitlement Repository 切换到 Ent，再进入 Elasticsearch Content Indexer。


## 2026-09-19 — V1.4 按 go-zero 官方规范校准 API-first 契约

### 本次依据

本阶段按 go-zero 官方推荐的 API/RPC 生成与分层方式推进：

- REST：.api 作为 HTTP 契约来源，保持 Handler → Logic → ServiceContext 的边界。
- RPC：.proto 作为 gRPC 契约来源，保持 Logic → RPC Server → ServiceContext。
- API Gateway 通过生成的 zRPC client 注入 ServiceContext 调用 drama-rpc，不直接访问数据库。
- 数据访问继续下沉到 Repository / Ent，业务 Logic 不拼接 SQL。
- 数据库结构仍以 deploy/mysql/init/*.sql 为 source of truth，Ent 不承担自动迁移。

官方参考：

- go-zero 项目结构与 API/RPC 生成规范：https://go-zero.dev/zh-cn/concepts/project-structure/
- go-zero RPC 服务指南：https://go-zero.dev/zh-cn/guides/quickstart/rpc-service/
- go-zero Proto DSL：https://go-zero.dev/zh-cn/reference/proto-dsl/
- go-zero RPC Client 配置：https://go-zero.dev/guides/grpc/client/configuration/

### 本次落地

- api/drama.api
  - 补齐 Drama 业务语义字段：subtitle、genres、tags、price_cents、currency、popularity、completion_rate、pay_rate、published_at。
  - 补齐行为事件上下文：source、request_id。
  - 保持既有 endpoint 与字段兼容，仅增加 API 契约字段。

### Ent 校准结论

已对照 Ent 官方 Edge Schema 文档确认：

- field.ID("drama_id", "genre_id") 作为 Edge Schema 的复合主键声明是合法模式，不再修改为普通 schema annotation。
- DramaGenre / DramaTag 继续采用 Edge Schema + .Through(...)。
- user_tag_profiles、behavior_events 的复合索引保持与 MySQL migration 一致。

### 当前验证状态

本次只完成 GitHub 静态契约校准，没有声称：

- goctl api go 已在本环境执行成功；
- goctl rpc protoc / protoc 已执行成功；
- make gen-ent 已执行成功；
- go test ./... 已通过；
- Docker/MySQL/Redis/Elasticsearch 全链路已通过。

下一步按官方生成链继续处理：

1. 校准 .api 与现有生成代码的对应关系；
2. 校准 .proto 与 generated pb/zRPC client；
3. Ent Schema 与真实 MySQL migration 做逐表一致性检查；
4. 生成 Ent client；
5. 将 drama-rpc 的 Drama/Episode/Entitlement Repository 迁移到 Ent；
6. 再进入 Content Indexer、Redis Recall 与 Recommendation Logic。
