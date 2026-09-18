# 视频存储与播放架构

## 目标

短剧视频统一采用 HLS（M3U8）播放，视频文件不进入 MySQL，也不建议直接放在应用服务器磁盘。

推荐链路：

管理后台
→ Media/OSS 上传服务
→ OSS Bucket
→ 转码/切片
→ `master.m3u8` / `index.m3u8` + TS/FMP4 分片
→ CDN
→ Client Tauri / H5 / Web
→ HLS 播放器

## OSS 对象结构

```
short-drama/
  {drama_id}/
    {episode_id}/
      source/
        original.mp4
      hls/
        master.m3u8
        720p/
          index.m3u8
          *.m4s
        1080p/
          index.m3u8
          *.m4s
      poster/
        cover.jpg
```

生产环境建议播放域名独立，例如 `https://media.example.com`，由 CDN 回源 OSS。

## 数据库

episodes 保存：

- video_format：默认 m3u8
- video_storage：默认 oss
- video_object_key：OSS 对象 Key
- video_playback_url：CDN 播放地址
- video_status：UPLOADING / PROCESSING / READY / FAILED
- video_size_bytes
- video_checksum
- video_updated_at

旧的 video_url 保留兼容，不再作为长期推荐字段。

## 上传

不能让浏览器把大视频经过 drama-api 转发。

正确方式：

1. admin-api 创建上传任务
2. 后端生成 OSS 临时凭证/预签名上传信息
3. 管理后台直接上传 OSS
4. 上传完成回调 media 服务
5. media 服务触发转码/切片
6. 状态 PROCESSING → READY
7. 返回 CDN M3U8 播放地址
8. episode 保存 playback_url/object_key

这样可以避免 API 网关成为视频带宽瓶颈。

## 播放

付费剧集必须先通过 drama-rpc 权益检查，再返回播放地址。

生产环境进一步使用短时签名 URL / CDN 防盗链；不要把 OSS AccessKey、Secret 写入客户端。

Tauri 客户端：

- Safari/WebKit 原生支持 HLS 时直接播放
- 其他 WebView 使用 hls.js
- 播放地址优先使用 video_playback_url
- 必须支持 HTTPS/CORS

## 后续服务拆分

建议增加独立 `media-rpc`：

```
merchant/admin
   ↓
media-api
   ↓
media-rpc
   ├── OSS
   ├── upload task
   ├── transcode
   ├── HLS packaging
   └── CDN/sign URL
```

这样旅游/短剧未来都可以复用媒体服务。

## 当前开发阶段

先完成：

1. 数据库字段
2. 后台剧集保存 M3U8 / OSS Key
3. 客户端 HLS 播放
4. media-rpc + OSS 上传
5. 转码队列
6. CDN 签名、防盗链

不要把原始 MP4 或 M3U8 分片提交进 Git，也不要存 MySQL BLOB。
