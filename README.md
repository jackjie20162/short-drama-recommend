# Short Drama Recommend

海外短剧推荐系统，基于 Go-zero + gRPC + MySQL + Redis 构建。

## 目标

构建面向海外市场的短剧内容、用户行为与个性化推荐平台，支持 Web / H5 / App。

## 第一阶段

- drama-api：REST API Gateway
- user-rpc：用户服务
- drama-rpc：短剧与剧集服务
- behavior-rpc：播放/点击/收藏等行为采集
- recommend-rpc：规则召回与个性化推荐
- MySQL：业务数据
- Redis：缓存、用户画像与推荐 Feed

## 第一条闭环

用户 → 首页 Feed → 短剧详情 → 播放 → 行为采集 → 用户兴趣画像 → 推荐服务 → 个性化 Feed

## 目录

```text
api/
rpc/
proto/
model/
common/
deploy/
docs/
scripts/
```

## 开发原则

1. API 与 RPC 解耦
2. 推荐算法可插拔
3. 用户行为事件可扩展
4. 多国家、多语言、多币种从数据模型层支持
5. 第一阶段优先保证本地 Docker 环境可运行
