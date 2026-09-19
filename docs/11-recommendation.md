# 全球短剧推荐模块

## 目标

推荐单位是「剧」，不是「集」。用户从首页进入一部剧后连续观看剧集，因此在线推荐返回 drama_id，播放和续播由 drama/episode/media 域负责。

全球化推荐统一把 country、language 作为在线排序特征，而不是简单硬过滤。这样新地区、新语言在没有足够行为数据时仍能使用全球热度完成冷启动。

## 当前在线链路

`client-tauri -> drama-api /api/v1/feed -> recommend-rpc -> MySQL`

当前第一阶段已经实现：

1. **地区匹配**：请求 country 与剧 country 匹配时增加排序分。
2. **语言匹配**：请求 language 与剧 language 匹配时增加排序分。
3. **全局热度**：根据 behavior_events 的互动次数计算热度分。
4. **观看深度**：根据累计 watch_seconds 增加观看质量分。
5. **已看降权**：登录用户已经产生过行为的剧降低排序，避免首页重复轰炸。
6. **付费历史信号**：用户已经拥有该剧 entitlement 时给出轻量加权。
7. **匿名冷启动**：没有 user_id 时仍能返回地区/语言 + 全局热度 feed。
8. **cursor 分页**：当前使用 offset cursor，保持 API 兼容；后续候选规模扩大后可切换 keyset cursor。
9. **推荐理由**：返回 region_match、language_match、popular、not_watched、high_watch_time 等 reason。

## 当前打分

第一阶段不是机器学习模型，而是可解释的规则精排：

`score = region + language + popularity + watch_depth + entitlement - seen_penalty`

它的目的不是替代最终 MMoE，而是先把数据闭环跑起来：

`曝光 -> 点击/观看 -> watch_seconds -> behavior_events -> 推荐`

## 下一阶段

### recall-rpc

按 `country × language` 建立 Redis ZSet 热度榜，并增加：

- 新剧召回
- 追更召回
- 相似题材召回
- 全局热门兜底

### rank-rpc

候选进入轻量粗排，再进入多目标精排：

- pCTR
- 观看时长
- 完播
- 付费解锁

最终：

`score = w1*pCTR + w2*expected_watch_time + w3*pComplete + w4*pPay`

权重按场景、区域、语言进行配置，而不是写死在客户端。

### rerank-rpc

负责：

- 同剧去重
- 题材多样性
- 运营位
- 付费引导
- 区域节日/时区适配
- 追更位例外

### feature-rpc / 模型推理

最终在线链路：

`recall -> rank -> rerank`

模型推理服务预留 Triton/ONNX 接口。第一阶段不直接引入 MMoE，避免在行为样本不足时过早模型化。

## 数据闭环

行为事件至少保留：

- exposure
- click
- watch_start
- watch_progress
- episode_complete
- unlock/pay

其中 watch_seconds 与 duration_seconds 用于计算观看深度/完播特征，付费事件作为稀疏目标单独采样和校准。

## 全球化原则

- 不同地区的热度榜不能直接混用。
- country/language 是特征，不应在没有内容供给时把用户过滤到空结果。
- 区域审核、版权和上下架状态必须在 item/内容域完成后再进入推荐候选。
- 一个全局模型 + region/language bias 优先于每地区独立模型。
- 新地区/新剧使用题材、语言、地区和全局热度进行冷启动。

## 当前限制

当前推荐服务是 MySQL 在线规则推荐，尚未接入：

- Redis ZSet 热度榜
- Kafka 行为流
- MMoE/ONNX/Triton
- A/B 实验
- 区域化模型 bias
- 运营位配置中心

这些属于第二阶段，不伪装成已经完成的功能。
