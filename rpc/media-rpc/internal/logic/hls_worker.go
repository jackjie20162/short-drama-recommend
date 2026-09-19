package logic

import (
 "context"
 "fmt"
 "io"
 "net/http"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "time"

 "short-drama-recommend/rpc/media-rpc/internal/storage"
)

func(s *MediaServer)processEpisode(episodeID int64,provider,sourceKey string){
 ctx,cancel:=context.WithTimeout(context.Background(),2*time.Hour);defer cancel()
 fail:=func(err error){_,_=s.svcCtx.DB.ExecContext(context.Background(),"UPDATE episodes SET video_status='FAILED',video_processing_error=?,video_updated_at=NOW() WHERE id=?",err.Error(),episodeID)}
 _,_=s.svcCtx.DB.ExecContext(ctx,"UPDATE episodes SET video_status='PROCESSING',video_processing_error='',video_updated_at=NOW() WHERE id=?",episodeID)
 st,err:=storage.FromSettings(ctx,s.lookup,provider);if err!=nil{fail(err);return}
 srcURL,err:=st.PresignGet(ctx,sourceKey,30*time.Minute);if err!=nil{fail(err);return}
 root,err:=os.MkdirTemp("","short-drama-hls-*");if err!=nil{fail(err);return};defer os.RemoveAll(root)
 src:=filepath.Join(root,"source.mp4")
 if err=download(ctx,srcURL,src);err!=nil{fail(err);return}
 hlsDir:=filepath.Join(root,"hls");if err=os.MkdirAll(hlsDir,0755);err!=nil{fail(err);return}
 master:=filepath.Join(hlsDir,"master.m3u8")
 cmd:=exec.CommandContext(ctx,"ffmpeg","-hide_banner","-loglevel","error","-y","-i",src,"-map","0:v:0?","-map","0:a:0?","-c","copy","-start_number","0","-hls_time","6","-hls_list_size","0","-hls_segment_filename",filepath.Join(hlsDir,"seg_%05d.ts"),master)
 if out,e:=cmd.CombinedOutput();e!=nil{fail(fmt.Errorf("ffmpeg failed: %v: %s",e,strings.TrimSpace(string(out))));return}
 hlsPrefix:=strings.TrimSuffix(sourceKey,filepath.Base(sourceKey))+"hls/"
 err=filepath.Walk(hlsDir,func(p string,info os.FileInfo,e error)error{
  if e!=nil{return e};if info.IsDir(){return nil}
  rel,e:=filepath.Rel(hlsDir,p);if e!=nil{return e}
  ct:="application/octet-stream";if strings.HasSuffix(rel,".m3u8"){ct="application/vnd.apple.mpegurl"}else if strings.HasSuffix(rel,".ts"){ct="video/mp2t"}
  return st.PutFile(ctx,hlsPrefix+filepath.ToSlash(rel),p,ct)
 });if err!=nil{fail(err);return}
 playback:=st.PublicURL(hlsPrefix+"master.m3u8")
 if playback==""{fail(fmt.Errorf("cdn_base_url is required for HLS playback"));return}
 _,err=s.svcCtx.DB.ExecContext(ctx,"UPDATE episodes SET video_format='m3u8',video_status='READY',video_playback_url=?,video_processing_error='',video_updated_at=NOW() WHERE id=?",playback,episodeID);if err!=nil{fail(err)}
}

func download(ctx context.Context,u,dst string)error{
 req,e:=http.NewRequestWithContext(ctx,http.MethodGet,u,nil);if e!=nil{return e}
 resp,e:=http.DefaultClient.Do(req);if e!=nil{return e};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return fmt.Errorf("source download returned %s",resp.Status)}
 f,e:=os.Create(dst);if e!=nil{return e};defer f.Close()
 _,e=io.Copy(f,resp.Body);return e
}
