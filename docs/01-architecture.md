# 海外短剧推荐系统架构

## 1. 总体架构

```text
Web / H5 / App
      |
      v
 drama-api (REST)
      |
  +---+---------+----------------+
  |             |                |
  v             v                v
user-rpc     drama-rpc      recommend-rpc
  |             |                |
  +-------------+----------------+
                |
           behavior-rpc
                |
        MySQL + Redis
                |
        Content Indexer
                |
          Elasticsearch
```

## 2. 服务职责

### drama-api
对外 REST 接口、鉴权、参数校验、聚合 RPC。

### user-rpc
用户、游客、地区、语言、基础兴趣画像。

### drama-rpc
短剧、剧集、分类、标签、上下架及基础内容查询。MySQL 是内容事实源。

### behavior-rpc
记录曝光、点击、播放、暂停、进度、完播、下一集、收藏、点赞、分享、搜索、购买等事件。

### recommend-rpc
当前负责规则推荐；后续拆分 recall/rank/rerank/feature 能力。

## 3. 内容数据分层

业务实体不等于模型输入：

```text
Drama
  |
  +--> MySQL business entity
  |
  +--> ES content document
  |
  +--> Feature Builder
          |
          +--> sparse
          +--> dense
          +--> semantic
                    |
                    v
                  MMoE
```

业务层必须保留 title、subtitle、description、genres、tags、country、language 等原始语义字段。

## 4. 推荐链路

V1：

```text
DB -> rule recall/rank -> Feed
```

目标链路：

```text
ES / Redis / DB
      |
      v
    recall
      |
      v
 candidate
      |
      v
 feature-rpc / Feature Builder
      |
      v
 rank-rpc / MMoE
      |
      v
  rerank
      |
      v
   Feed
```

## 5. RecommendationFeature

### Sparse

user_id、drama_id、country_id、language_id、genre_id、tag_id。

### Dense

用户观看时长、用户完播、用户付费率、session 数，以及短剧热度、完播率、付费率、剧龄等。

### Semantic

title_embedding、description_embedding、genre_embedding、tag_embedding。

原始文本不直接作为普通 dense 数值特征；通过 embedding pipeline 转换。

## 6. 国际化

核心实体支持 country、language、locale、timezone、currency。

country/language 在推荐中默认作为特征或 boost，而不是无条件硬过滤，以避免新市场/冷启动时返回空结果。

## 7. 搜索与索引

MySQL -> Content Indexer -> Elasticsearch。

ES 用于搜索、内容召回、标签过滤、语义检索；MySQL 仍是业务事实源。

上下架状态必须在进入推荐候选前校验。

## 8. 视频

业务服务只保存视频元数据和播放地址，视频文件交给对象存储 + CDN，支持 HLS 多码率。
