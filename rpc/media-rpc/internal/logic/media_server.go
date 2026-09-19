package logic
import ("context";"errors";"fmt";"net/url";"path";"strings";"short-drama-recommend/rpc/media-rpc/internal/svc";"short-drama-recommend/rpc/media-rpc/pb")
type MediaServer struct{pb.UnimplementedMediaServiceServer;svcCtx *svc.ServiceContext}
func NewMediaServer(c *svc.ServiceContext)*MediaServer{return &MediaServer{svcCtx:c}}
func(s *MediaServer)CreateUpload(ctx context.Context,r *pb.CreateUploadRequest)(*pb.CreateUploadResponse,error){
 if r.GetDramaId()<=0||r.GetEpisodeId()<=0{return nil,errors.New("drama_id and episode_id are required")}
 provider:=strings.ToUpper(r.GetProvider());if provider==""||provider=="DEFAULT"{provider="OSS";if v,ok:=s.setting(ctx,"storage","default_provider");ok&&v!=""{provider=strings.ToUpper(v)}}
 if provider!="OSS"&&provider!="S3"{return nil,errors.New("unsupported storage provider")}
 filename:=path.Base(r.GetFilename());if filename==""||filename=="."{filename="source.mp4"}
 key:=fmt.Sprintf("short-drama/%s/%s/%d/%d/source/%s",sanitize(r.GetCountry(),"global"),sanitize(r.GetLanguage(),"und"),r.GetDramaId(),r.GetEpisodeId(),filename)
 base,_:=s.setting(ctx,"storage","cdn_base_url");playback:=strings.TrimRight(base,"/")+"/"+key
 return &pb.CreateUploadResponse{Provider:provider,ObjectKey:key,UploadUrl:playback,PlaybackUrl:playback,ExpiresSeconds:3600},nil
}
func(s *MediaServer)CompleteUpload(ctx context.Context,r *pb.CompleteUploadRequest)(*pb.MediaAsset,error){
 if r.GetEpisodeId()<=0||r.GetObjectKey()==""{return nil,errors.New("episode_id and object_key are required")}
 provider:=strings.ToUpper(r.GetProvider());if provider==""{provider="OSS"};base,_:=s.setting(ctx,"storage","cdn_base_url");playback:=strings.TrimRight(base,"/")+"/"+r.GetObjectKey()
 _,err:=s.svcCtx.DB.ExecContext(ctx,"UPDATE episodes SET video_provider=?,video_storage=?,video_object_key=?,video_playback_url=?,video_format='m3u8',video_status='PROCESSING',video_size_bytes=?,video_checksum=?,video_updated_at=NOW() WHERE id=?",provider,strings.ToLower(provider),r.GetObjectKey(),playback,r.GetSizeBytes(),r.GetChecksum(),r.GetEpisodeId());if err!=nil{return nil,err}
 return &pb.MediaAsset{Provider:provider,ObjectKey:r.GetObjectKey(),PlaybackUrl:playback,Status:"PROCESSING"},nil
}
func(s *MediaServer)GetPlayback(ctx context.Context,r *pb.GetPlaybackRequest)(*pb.PlaybackResponse,error){
 if r.GetEpisodeId()<=0{return nil,errors.New("episode_id is required")};var u,f string;var paid int
 if err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT video_playback_url,video_format,is_paid FROM episodes WHERE id=? AND status=1",r.GetEpisodeId()).Scan(&u,&f,&paid);err!=nil{return nil,err}
 if paid==1{if r.GetUserId()<=0{return nil,errors.New("login required")};var n int;if err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT COUNT(*) FROM user_entitlements ue JOIN episodes e ON e.drama_id=ue.drama_id WHERE ue.user_id=? AND e.id=?",r.GetUserId(),r.GetEpisodeId()).Scan(&n);err!=nil{return nil,err};if n==0{return nil,errors.New("episode is locked")}}
 if u==""{return nil,errors.New("playback url is not ready")};return &pb.PlaybackResponse{Url:u,Format:f,ExpiresSeconds:300},nil
}
func(s *MediaServer)setting(ctx context.Context,g,k string)(string,bool){var v string;var secret int;if err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT setting_value,is_secret FROM system_settings WHERE setting_group=? AND setting_key=? LIMIT 1",g,k).Scan(&v,&secret);err!=nil{return "",false};return v,true}
func sanitize(v,f string)string{v=strings.TrimSpace(v);if v==""{return f};return url.PathEscape(v)}