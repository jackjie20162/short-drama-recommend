package logic

import (
 "context"
 "errors"
 "fmt"
 "net/url"
 "path"
 "strings"
 "time"

 "short-drama-recommend/internal/settings"
 "short-drama-recommend/rpc/media-rpc/internal/storage"
 "short-drama-recommend/rpc/media-rpc/internal/svc"
 "short-drama-recommend/rpc/media-rpc/pb"
)

type MediaServer struct{pb.UnimplementedMediaServiceServer;svcCtx *svc.ServiceContext}
func NewMediaServer(c *svc.ServiceContext)*MediaServer{return &MediaServer{svcCtx:c}}

func(s *MediaServer)CreateUpload(ctx context.Context,r *pb.CreateUploadRequest)(*pb.CreateUploadResponse,error){
 if r.GetDramaId()<=0||r.GetEpisodeId()<=0{return nil,errors.New("drama_id and episode_id are required")}
 filename:=path.Base(r.GetFilename());if filename==""||filename=="."{filename="source.mp4"}
 key:=fmt.Sprintf("short-drama/%s/%s/%d/%d/source/%s",sanitize(r.GetCountry(),"global"),sanitize(r.GetLanguage(),"und"),r.GetDramaId(),r.GetEpisodeId(),filename)
 st,e:=storage.FromSettings(ctx,s.lookup,"");if e!=nil{return nil,e}
 ct:=r.GetContentType();if ct==""{ct="video/mp4"}
 u,e:=st.PresignPut(ctx,key,ct,60*time.Minute);if e!=nil{return nil,e}
 return &pb.CreateUploadResponse{Provider:providerName(st),ObjectKey:key,UploadUrl:u,PlaybackUrl:st.PublicURL(key),ExpiresSeconds:3600},nil
}

func(s *MediaServer)CompleteUpload(ctx context.Context,r *pb.CompleteUploadRequest)(*pb.MediaAsset,error){
 if r.GetEpisodeId()<=0||r.GetObjectKey()==""{return nil,errors.New("episode_id and object_key are required")}
 st,e:=storage.FromSettings(ctx,s.lookup,r.GetProvider());if e!=nil{return nil,e}
 size,e:=st.Head(ctx,r.GetObjectKey());if e!=nil{return nil,fmt.Errorf("uploaded object verification failed: %w",e)}
 if r.GetSizeBytes()>0&&size!=r.GetSizeBytes(){return nil,fmt.Errorf("uploaded size mismatch: expected %d got %d",r.GetSizeBytes(),size)}
 playback:=st.PublicURL(r.GetObjectKey())
 if playback==""{playback,_=st.PresignGet(ctx,r.GetObjectKey(),5*time.Minute)}
 _,e=s.svcCtx.DB.ExecContext(ctx,"UPDATE episodes SET video_provider=?,video_storage=?,video_object_key=?,video_playback_url=?,video_format='source',video_status='UPLOADED',video_size_bytes=?,video_checksum=?,video_updated_at=NOW() WHERE id=?",providerName(st),strings.ToLower(providerName(st)),r.GetObjectKey(),playback,size,r.GetChecksum(),r.GetEpisodeId());if e!=nil{return nil,e}
 go s.processEpisode(r.GetEpisodeId(),providerName(st),r.GetObjectKey())
 return &pb.MediaAsset{Provider:providerName(st),ObjectKey:r.GetObjectKey(),PlaybackUrl:playback,Status:"PROCESSING"},nil
}

func(s *MediaServer)GetPlayback(ctx context.Context,r *pb.GetPlaybackRequest)(*pb.PlaybackResponse,error){
 if r.GetEpisodeId()<=0{return nil,errors.New("episode_id is required")}
 var key,f string;var paid int
 if e:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT video_object_key,video_format,is_paid FROM episodes WHERE id=? AND status=1",r.GetEpisodeId()).Scan(&key,&f,&paid);e!=nil{return nil,e}
 if paid==1{if r.GetUserId()<=0{return nil,errors.New("login required")};var n int;if e:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT COUNT(*) FROM user_entitlements ue JOIN episodes ep ON ep.drama_id=ue.drama_id WHERE ue.user_id=? AND ep.id=?",r.GetUserId(),r.GetEpisodeId()).Scan(&n);e!=nil{return nil,e};if n==0{return nil,errors.New("episode is locked")}}
 if key==""{return nil,errors.New("playback url is not ready")}
 st,e:=storage.FromSettings(ctx,s.lookup,"");if e!=nil{return nil,e}
 u,e:=st.PresignGet(ctx,key,5*time.Minute);if e!=nil{return nil,e}
 return &pb.PlaybackResponse{Url:u,Format:"m3u8",ExpiresSeconds:300},nil
}
func(s *MediaServer)lookup(ctx context.Context,k string)(string,error){var v string;var secret int;e:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT setting_value,is_secret FROM system_settings WHERE setting_group='storage' AND setting_key=? LIMIT 1",k).Scan(&v,&secret);if e!=nil{return "",e};if secret==1{return settings.Decrypt(v)};return v,nil}
func providerName(st storage.Storage)string{if strings.Contains(fmt.Sprintf("%T",st),"s3Storage"){return "S3"};return "OSS"}
func sanitize(v,f string)string{v=strings.TrimSpace(v);if v==""{return f};return url.PathEscape(v)}
