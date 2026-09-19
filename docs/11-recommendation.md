# 全球短剧推荐模块

## 目标

推荐单位是「剧」，不是「集」。用户从首页进入一部剧后连续观看剧集，因此在线推荐返回 drama_id，播放和续播由 drama/episode/media 域负责。

全球化推荐统一把 country、language 作为在线排序特征，而不是简单硬过滤。这样新地区、新语言在没有足够行为数据时仍能使用全球热度完成冷启动。

## 当前在线链路

`client-tauri -> drama-api /api/v1/feed -> recommend-rpc -> MySQL`

当前第一阶段已经实现：

1. 地区匹配加分。
2. 语言匹配加分。
3. 全局热度。
4. 观看深度。
5. 已看降权。
6. 付费历史信号。
7. 匿名冷启动。
8. cursor 分页。
9. 推荐理由。

## 内容数据升级

从 V1.1 开始，推荐系统采用统一内容契约：

```text
Drama
├── title
├── subtitle
├── description
├── genres[]
├── tags[]
├── country
├── language
├── popularity
├── completion_rate
└── pay_rate
       |
       +--> Elasticsearch content document
       |
       +--> RecommendationFeature
```

因此名称、简介、题材和标签属于业务数据；模型侧通过 Feature Builder 转换为 sparse/dense/semantic 特征。

## 当前打分

第一阶段不是机器学习模型，而是可解释的规则精排：

`score = region + language + popularity + watch_depth + entitlement - seen_penalty`

它的目的不是替代最终 MMoE，而是先把数据闭环跑起来：

`曝光 -> 点击/观看 -> watch_seconds -> behavior_events -> 推荐`

## 下一阶段

### recall-rpc

按 country × language 建立 Redis ZSet 热度榜，并增加：

- 新剧召回
- 追更召回
- 相似题材召回
- 全局热门兜底
- ES 内容/语义召回

### feature-rpc / Feature Builder

统一生成：

- sparse features
- dense features
- semantic embeddings

语义输入由 title、subtitle、description、genres、tags、language、country 构成。

### rank-rpc

候选进入轻量粗排，再进入多目标精排：

- pCTR
- 观看时长
- 完播
- 付费解锁

最终：

`score = w1*pCTR + w2*expected_watch_time + w3*pComplete + w4*pPay`

权重按场景、区域、语言配置。

### rerank-rpc

负责同剧去重、题材多样性、运营位、付费引导、区域节日/时区适配及追更位例外。

## 当前限制

尚未宣称完成：

- Elasticsearch 实际部署
- Redis ZSet 热度榜
- Kafka 行为流
- MMoE/ONNX/Triton
- A/B 实验
- 区域化模型 bias
- 运营位配置中心

这些属于后续阶段。
