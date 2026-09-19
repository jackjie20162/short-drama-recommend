# 短剧内容数据契约与推荐特征模型

> 版本：V1.1
> 更新时间：2026-09-19
> 目的：统一业务内容层、搜索层和推荐模型层的数据契约，避免后续 MMoE/ES 接入时重复改模型。

## 1. 分层原则

系统严格区分三层：

1. **业务实体层**：保存可展示、可运营、可审核的原始内容。
2. **搜索/召回层**：把业务实体转换为可检索字段、过滤字段和语义向量。
3. **模型特征层**：把用户与短剧转换成 sparse / dense / semantic 特征，供 MMoE 等排序模型使用。

title、description 等内容字段不直接作为 MMoE 的普通数值输入，而是通过 Feature Builder 转换为语义特征。

## 2. Drama 业务实体

| 字段 | 类型 | 用途 |
|---|---|---|
| drama_id | uint64 | 主键 |
| title | string | 主标题 |
| subtitle | string | 副标题 |
| description | text | 剧情简介 |
| cover | string | 封面 |
| country | string | 内容所属国家/地区 |
| language | string | 内容主语言 |
| genres[] | string | 类型，如 Romance/Drama |
| tags[] | string | 细粒度标签 |
| total_episodes | int | 集数 |
| episode_count | int | 与 total_episodes 兼容，内部统一为 total_episodes |
| popularity | float | 聚合热度 |
| completion_rate | float | 完播率 |
| pay_rate | float | 付费转化率 |
| status | int | 上下架 |
| published_at | datetime | 发布时间 |

## 3. Episode 业务实体

| 字段 | 类型 |
|---|---|
| episode_id | uint64 |
| drama_id | uint64 |
| episode_no | int |
| title | string |
| description | text |
| duration_seconds | int |
| video_url | string |
| poster_url | string |
| is_paid | bool |
| status | int |

推荐单位仍然是 Drama，不是 Episode。

## 4. BehaviorEvent

事件至少支持：

- impression
- click
- play
- pause
- watch_progress
- complete
- next_episode
- favorite
- like
- share
- search
- purchase

关键上下文字段：

- event_id
- event_time
- user_id
- drama_id
- episode_id
- event_type
- watch_seconds
- duration_seconds
- country
- language
- device
- source
- request_id

## 5. RecommendationFeature

### Sparse

- user_id
- drama_id
- country_id
- language_id
- genre_id[]
- tag_id[]

### Dense

- user_watch_seconds
- user_completion
- user_pay_rate
- user_sessions
- drama_popularity
- drama_completion
- drama_pay_rate
- drama_age_days

### Semantic

- title_embedding
- description_embedding
- genre_embedding
- tag_embedding

V1.1 不要求在线 MMoE 立即使用全部 semantic embedding；先完成数据契约、离线构建和接口预留。

## 6. Feature Builder

统一流程：

```text
Drama/User/Behavior
       |
       v
Feature Builder
       |
       +--> SparseFeature
       +--> DenseFeature
       +--> SemanticFeature
       |
       v
Ranker / MMoE
```

语义文本建议：

```text
title + subtitle + description + genres + tags + language + country
```

多语言内容不能简单拼接机器翻译结果覆盖原文；原始语言字段必须保留，翻译/embedding 版本单独记录。

## 7. 冷启动

新剧：

- genre/tag
- language
- country
- published_at
- 全局热度先验

新用户：

- country
- language
- device
- session context
- 热门/新剧候选

没有内容供给时，不使用 country/language 做硬过滤导致空结果。

## 8. 兼容策略

- Proto 只新增字段，不复用旧字段编号。
- MySQL 使用增量 migration，不修改历史 migration。
- ES index 采用 versioned mapping。
- MMoE 输入通过 Feature Builder 生成，业务实体与模型输入解耦。
