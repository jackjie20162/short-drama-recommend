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
