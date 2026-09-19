# 开发日志

## 2026-09-19 — V1.1 内容数据与推荐特征契约

### 本次检查

检查仓库：`jackjie20162/short-drama-recommend`

当前已有：

- Go-zero + gRPC 服务结构
- MySQL 用户、短剧、剧集、标签、行为表
- drama-api / drama-rpc / behavior-rpc / recommend-rpc
- 推荐 V1 规则排序
- docs/01-architecture.md、docs/11-recommendation.md 等项目文档

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

- 增加 `docs/12-content-data-contract.md`
- 增加 `docs/13-elasticsearch-index.md`
- 增加本开发日志
- 增加增量数据库 migration：`010_recommend_content.sql`
- 扩展 Drama / Episode protobuf 的内容字段
- 为下一阶段 Feature Builder 预留统一输入结构

### 验证状态

本次通过 GitHub 仓库静态检查完成架构与契约落地。

未宣称已经完成：

- Elasticsearch 实际部署
- embedding 模型训练
- MMoE 训练
- Triton/ONNX 在线推理
- Kafka 行为流
- A/B 实验

### 下一步

1. Feature Builder 实现
2. Drama 内容同步 ES
3. Redis country × language 热度榜
4. MMoE 离线样本生成
5. 模型训练与离线评估
6. 在线 rank-rpc 接入
