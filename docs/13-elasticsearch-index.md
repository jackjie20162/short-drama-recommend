# Elasticsearch 短剧索引契约

> 版本：V1.1
> 更新时间：2026-09-19

## 1. 索引定位

ES 用于：

- 短剧关键词搜索
- country/language 召回
- genre/tag 过滤
- 新剧召回
- 相似内容召回
- title/description 语义检索

MySQL 仍是业务事实源；ES 不作为最终业务状态源。

## 2. 建议索引

```text
drama_content_v1
```

核心 document：

```json
{
  "drama_id": 20086,
  "title": "Love After Divorce",
  "subtitle": "A second chance at love",
  "description": "After discovering her husband's betrayal...",
  "country": "US",
  "language": "en",
  "genres": ["Romance", "Drama", "Revenge"],
  "tags": ["divorce", "CEO", "second-chance-love"],
  "total_episodes": 80,
  "status": 1,
  "published_at": "2026-09-01T10:00:00Z",
  "popularity": 0.82,
  "completion_rate": 0.64,
  "pay_rate": 0.18
}
```

## 3. 字段用途

| 字段 | ES 用途 |
|---|---|
| title | BM25 + 语义检索 |
| subtitle | BM25 + 语义检索 |
| description | BM25 + 语义检索 |
| genres | filter + boost |
| tags | filter + boost |
| country | filter/boost |
| language | filter/boost |
| status | 必须过滤 |
| published_at | freshness |
| popularity | recall/rank |
| completion_rate | rank |
| pay_rate | rank |

## 4. 语义向量

V1.1 先预留：

- title_vector
- description_vector
- content_vector

向量模型版本必须作为 metadata 保存，例如：

```text
embedding_model
embedding_version
embedding_updated_at
```

这样模型升级时可以重建向量而不污染业务数据。

## 5. 同步

```text
Drama Admin
   |
   v
MySQL
   |
   v
Content Indexer
   |
   v
Elasticsearch
```

MySQL 更新成功后，通过异步事件或 outbox 推送索引更新。

上下架必须优先同步；推荐/搜索服务不得返回 status != published 的内容。

## 6. 与推荐系统关系

```text
ES / Redis / DB
      |
      v
Recall
      |
      v
Candidate
      |
      v
Feature Builder
      |
      v
Rank / MMoE
      |
      v
Rerank
```

ES 负责内容理解和召回，不替代 MMoE 排序。
