package storage

import (
 "context"
 "fmt"
 "net/http"
 "os"
 "strings"
 "time"

 "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
 osscred "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
 "github.com/aws/aws-sdk-go-v2/aws"
 awscfg "github.com/aws/aws-sdk-go-v2/config"
 "github.com/aws/aws-sdk-go-v2/credentials"
 "github.com/aws/aws-sdk-go-v2/service/s3"
 "short-drama-recommend/internal/settings"
)

type Config struct {
 Provider string
 Endpoint string
 Region string
 Bucket string
 AccessKey string
 SecretKey string
 CDNBaseURL string
}

type Storage interface {
 PresignPut(context.Context,string,string,time.Duration)(string,error)
 PresignGet(context.Context,string,time.Duration)(string,error)
 Head(context.Context,string)(int64,error)
 PutFile(context.Context,string,string,string) error
 PublicURL(string) string
}

func FromSettings(ctx context.Context, lookup func(context.Context,string)(string,error), provider string)(Storage,error){
 p:=strings.ToUpper(strings.TrimSpace(provider))
 if p==""||p=="DEFAULT"{p="OSS";if v,e:=lookup(ctx,"default_provider");e==nil&&v!=""{p=strings.ToUpper(v)}}
 var c Config
 c.Provider=p
 c.CDNBaseURL,_=lookup(ctx,"cdn_base_url")
 if p=="OSS"{
  c.Endpoint,_=lookup(ctx,"oss_endpoint");c.Region,_=lookup(ctx,"oss_region");c.Bucket,_=lookup(ctx,"oss_bucket");c.AccessKey,_=lookup(ctx,"oss_access_key")
  enc,e:=lookup(ctx,"oss_secret_key");if e!=nil{return nil,e};c.SecretKey,e=settings.Decrypt(enc);if e!=nil{return nil,e}
  if c.Region==""||c.Bucket==""||c.AccessKey==""||c.SecretKey==""{return nil,fmt.Errorf("OSS storage is not fully configured")}
  return NewOSS(c)
 }
 if p=="S3"{
  c.Endpoint,_=lookup(ctx,"s3_endpoint");c.Region,_=lookup(ctx,"s3_region");c.Bucket,_=lookup(ctx,"s3_bucket");c.AccessKey,_=lookup(ctx,"s3_access_key")
  enc,e:=lookup(ctx,"s3_secret_key");if e!=nil{return nil,e};c.SecretKey,e=settings.Decrypt(enc);if e!=nil{return nil,e}
  if c.Region==""||c.Bucket==""||c.AccessKey==""||c.SecretKey==""{return nil,fmt.Errorf("S3 storage is not fully configured")}
  return NewS3(c)
 }
 return nil,fmt.Errorf("unsupported storage provider %q",p)
}

type ossStorage struct{client *oss.Client;bucket,cdn string}
func NewOSS(c Config)(Storage,error){
 provider:=osscred.NewStaticCredentialsProvider(c.AccessKey,c.SecretKey)
 cfg:=oss.LoadDefaultConfig().WithCredentialsProvider(provider).WithRegion(c.Region)
 if c.Endpoint!=""{cfg=cfg.WithEndpoint(c.Endpoint)}
 return &ossStorage{client:oss.NewClient(cfg),bucket:c.Bucket,cdn:c.CDNBaseURL},nil
}
func(s *ossStorage)PresignPut(ctx context.Context,key,contentType string,ttl time.Duration)(string,error){
 r,e:=s.client.Presign(ctx,&oss.PutObjectRequest{Bucket:oss.Ptr(s.bucket),Key:oss.Ptr(key),ContentType:oss.Ptr(contentType)},oss.PresignExpires(ttl));if e!=nil{return "",e};return r.URL,nil
}
func(s *ossStorage)PresignGet(ctx context.Context,key string,ttl time.Duration)(string,error){
 r,e:=s.client.Presign(ctx,&oss.GetObjectRequest{Bucket:oss.Ptr(s.bucket),Key:oss.Ptr(key)},oss.PresignExpires(ttl));if e!=nil{return "",e};return r.URL,nil
}
func(s *ossStorage)Head(ctx context.Context,key string)(int64,error){
 u,e:=s.PresignGet(ctx,key,2*time.Minute);if e!=nil{return 0,e};req,e:=http.NewRequestWithContext(ctx,http.MethodHead,u,nil);if e!=nil{return 0,e};resp,e:=http.DefaultClient.Do(req);if e!=nil{return 0,e};defer resp.Body.Close();if resp.StatusCode/100!=2{return 0,fmt.Errorf("storage HEAD returned %s",resp.Status)};return resp.ContentLength,nil
}
func(s *ossStorage)PutFile(ctx context.Context,key,filePath,contentType string)error{f,e:=os.Open(filePath);if e!=nil{return e};defer f.Close();_,e=s.client.PutObject(ctx,&oss.PutObjectRequest{Bucket:oss.Ptr(s.bucket),Key:oss.Ptr(key),Body:f,ContentType:oss.Ptr(contentType)});return e}
func(s *ossStorage)PublicURL(key string)string{if s.cdn!=""{return strings.TrimRight(s.cdn,"/")+"/"+strings.TrimLeft(key,"/")};return ""}

type s3Storage struct{client *s3.Client;presign *s3.PresignClient;bucket,cdn string}
func NewS3(c Config)(Storage,error){
 cfg,e:=awscfg.LoadDefaultConfig(context.Background(),awscfg.WithRegion(c.Region),awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKey,c.SecretKey,"")));if e!=nil{return nil,e}
 client:=s3.NewFromConfig(cfg,func(o *s3.Options){o.UsePathStyle=true;if c.Endpoint!=""{o.BaseEndpoint=aws.String(c.Endpoint)}})
 return &s3Storage{client:client,presign:s3.NewPresignClient(client),bucket:c.Bucket,cdn:c.CDNBaseURL},nil
}
func(s *s3Storage)PresignPut(ctx context.Context,key,contentType string,ttl time.Duration)(string,error){
 r,e:=s.presign.PresignPutObject(ctx,&s3.PutObjectInput{Bucket:aws.String(s.bucket),Key:aws.String(key),ContentType:aws.String(contentType)},func(o *s3.PresignOptions){o.Expires=ttl});if e!=nil{return "",e};return r.URL,nil
}
func(s *s3Storage)PresignGet(ctx context.Context,key string,ttl time.Duration)(string,error){
 r,e:=s.presign.PresignGetObject(ctx,&s3.GetObjectInput{Bucket:aws.String(s.bucket),Key:aws.String(key)},func(o *s3.PresignOptions){o.Expires=ttl});if e!=nil{return "",e};return r.URL,nil
}
func(s *s3Storage)Head(ctx context.Context,key string)(int64,error){
 u,e:=s.PresignGet(ctx,key,2*time.Minute);if e!=nil{return 0,e};req,e:=http.NewRequestWithContext(ctx,http.MethodHead,u,nil);if e!=nil{return 0,e};resp,e:=http.DefaultClient.Do(req);if e!=nil{return 0,e};defer resp.Body.Close();if resp.StatusCode/100!=2{return 0,fmt.Errorf("storage HEAD returned %s",resp.Status)};return resp.ContentLength,nil
}
func(s *s3Storage)PutFile(ctx context.Context,key,filePath,contentType string)error{f,e:=os.Open(filePath);if e!=nil{return e};defer f.Close();_,e=s.client.PutObject(ctx,&s3.PutObjectInput{Bucket:aws.String(s.bucket),Key:aws.String(key),Body:f,ContentType:aws.String(contentType)});return e}
func(s *s3Storage)PublicURL(key string)string{if s.cdn!=""{return strings.TrimRight(s.cdn,"/")+"/"+strings.TrimLeft(key,"/")};return ""}
