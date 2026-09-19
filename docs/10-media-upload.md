# Media upload and HLS processing

## Current closed loop

1. Admin calls POST /api/v1/admin/episodes/:id/upload.
2. drama-admin-api calls media-rpc.CreateUpload.
3. media-rpc reads storage settings and creates a real OSS or S3 pre-signed PUT URL.
4. Admin uploads the original MP4 directly to object storage. The API does not carry video bytes.
5. Admin calls POST /api/v1/admin/episodes/:id/upload/complete.
6. media-rpc verifies the object with a signed HEAD request and records provider, key and size.
7. media-rpc starts the FFmpeg worker asynchronously.
8. FFmpeg creates HLS master.m3u8 and TS segments.
9. The worker uploads all HLS artifacts to object storage and sets the episode to READY.
10. Playback returns a short-lived signed URL for the HLS master.

## Object layout

short-drama/{country}/{language}/{drama_id}/{episode_id}/source/original.mp4
short-drama/{country}/{language}/{drama_id}/{episode_id}/source/hls/master.m3u8
short-drama/{country}/{language}/{drama_id}/{episode_id}/source/hls/seg_00000.ts

## Required production settings

For storage:
- default_provider
- OSS: oss_region, oss_bucket, oss_access_key, oss_secret_key, optional oss_endpoint
- S3: s3_region, s3_bucket, s3_access_key, s3_secret_key, optional s3_endpoint
- cdn_base_url

HLS playback should use a CDN/base URL that can serve the master playlist and relative segments. Without cdn_base_url, the worker marks the episode FAILED instead of exposing a broken HLS URL.

## Runtime

The media Docker image includes FFmpeg. HLS processing is asynchronous so the upload RPC does not block for the duration of transcoding.

## Security

Storage secret keys are read from encrypted system_settings values. Browser clients receive only pre-signed upload URLs. Paid playback still checks user_entitlements before a playback URL is returned.
