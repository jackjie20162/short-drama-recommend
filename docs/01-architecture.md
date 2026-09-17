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
```

## 2. 服务职责

### drama-api
对外 REST 接口、鉴权、参数校验、聚合 RPC。

### user-rpc
用户、游客、地区、语言、基础兴趣画像。

### drama-rpc
短剧、剧集、分类、标签、上下架及基础内容查询。

### behavior-rpc
记录曝光、点击、播放、暂停、完播、下一集、收藏、点赞、分享、搜索等事件。

### recommend-rpc
候选召回、用户兴趣匹配、热度、新鲜度、地域/语言过滤、排序和多样性处理。

## 3. 推荐 V1

第一版采用可解释的规则推荐，不依赖机器学习模型：

```text
score = personalization
      + popularity
      + freshness
      + completion_rate
      + country_score
      + language_score
      + diversity_score
```

后续可替换为机器学习排序模型，API/RPC 契约保持稳定。

## 4. 国际化

核心实体必须支持 country、language、locale、timezone、currency，避免后期改表。

## 5. 视频

业务服务只保存视频元数据和播放地址，视频文件交给对象存储 + CDN，后续支持 HLS 多码率。
