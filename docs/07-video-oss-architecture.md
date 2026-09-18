# 视频存储与播放架构

## 目标

短剧视频统一采用 **HLS（M3U8）**。原始视频、M3U8 playlist、分片都放对象存储/CDN，不放 MySQL，也不把大文件长期放应用服务器。

### 对象存储必须支持多云

第一版设计同时支持：

- **阿里云 OSS**
- **AWS S3**

不要把 media 服务写死成 OSS。统一抽象：

```
Media Storage Interface
 ├── Aliyun OSS Adapter
 └── AWS S3 Adapter
```

未来可以继续增加 Cloudflare R2、腾讯云 COS 等实现，而业务层不用改。

## 推荐链路

```
Admin
  ↓
media-api
  ↓
media-rpc
  ↓
Storage Adapter
  ├── OSS
  └── S3
  ↓
Transcode / HLS
  ↓
CDN
  ↓
master.m3u8
  ↓
Tauri / H5 / Web
```

## 对象结构

```
short-drama/{country}/{language}/{drama_id}/{episode_id}/
├── source/original.mp4
├── hls/master.m3u8
├── hls/720p/index.m3u8
├── hls/1080p/index.m3u8
└── poster/cover.jpg
```

## 数据库

episodes 保存：

- video_provider：OSS / S3
- video_storage：oss / s3
- video_region
- video_bucket
- video_object_key
- video_playback_url
- video_format：m3u8
- video_status：UPLOADING / PROCESSING / READY / FAILED
- size / checksum / 更新时间

客户端永远拿不到云存储 AccessKey / Secret。

## 发布短剧必须选择国家和语言

短剧是全球发行，因此发布页面必须明确：

- 国家 / 市场
- 内容语言
- 播放语言
- 后续可以扩展字幕语言

建议内部使用标准 ISO code，例如：

- US / en
- GB / en
- CA / en
- AE / ar
- JP / ja
- KR / ko
- CN / zh-CN

国家和语言既用于内容筛选，也用于推荐、SEO、CDN 路由和后续多语言版本管理。

## 上传设计

不能：

```
Browser → drama-api → OSS/S3
```

应该：

```
Admin → media-api → presigned upload
                         ↓
                    OSS / S3
                         ↓
                    transcode
                         ↓
                       HLS
                         ↓
                       CDN
```

这样 API 不承担视频带宽。

## 播放

付费剧：

```
Client
 ↓
drama-api
 ↓
drama-rpc entitlement check
 ↓
生成短时播放地址
 ↓
CDN
 ↓
M3U8
```

生产环境使用 CDN signed URL / token 防盗链。

## 发布页面

```
国家： [United States ▼]
语言： [English ▼]

视频：
[上传视频]

存储：
[OSS ▼]
或
[AWS S3 ▼]

转码：
[HLS / M3U8]

状态：
上传中 → 转码中 → READY
```

国家和语言属于内容元数据，不应该依赖客户端 IP 推断。

## 下一步

1. media-rpc storage interface
2. OSS adapter
3. AWS S3 adapter
4. presigned upload
5. FFmpeg/HLS 转码任务
6. CDN signed URL
7. Admin 上传组件
8. 视频状态回调
