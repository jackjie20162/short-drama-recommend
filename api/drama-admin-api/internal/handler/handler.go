package handler

import(
 "database/sql"
 "short-drama-recommend/internal/db"
 "encoding/json"
 "errors"
 "net/http"
 "strconv"
 "strings"
 "github.com/zeromicro/go-zero/rest"
 "github.com/zeromicro/go-zero/zrpc"
 "short-drama-recommend/api/drama-admin-api/internal/config"
 "short-drama-recommend/rpc/drama-rpc/pb"
 paymentpb "short-drama-recommend/rpc/payment-rpc/pb"
 mediapb "short-drama-recommend/rpc/media-rpc/pb"
 "short-drama-recommend/internal/settings"
)

type ServiceContext struct{Config config.Config;DB *sql.DB;Drama pb.DramaAdminServiceClient;Payment paymentpb.PaymentServiceClient;Media mediapb.MediaServiceClient}
func NewServiceContext(c config.Config)*ServiceContext{
 database,err:=db.OpenMySQL(c.Mysql.DataSource);if err!=nil{panic(err)}
 return &ServiceContext{Config:c,DB:database,Drama:pb.NewDramaAdminServiceClient(zrpc.MustNewClient(c.DramaRpc).Conn()),Payment:paymentpb.NewPaymentServiceClient(zrpc.MustNewClient(c.PaymentRpc).Conn()),Media:mediapb.NewMediaServiceClient(zrpc.MustNewClient(c.MediaRpc).Conn())}
}
func RegisterRoutes(s *rest.Server,ctx *ServiceContext){
 s.AddRoutes([]rest.Route{
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas",Handler:func(w http.ResponseWriter,r *http.Request){var x pb.CreateDramaRequest;if readJSON(w,r,&x)!=nil{return};v,e:=ctx.Drama.CreateDrama(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/dramas",Handler:func(w http.ResponseWriter,r *http.Request){q:=r.URL.Query();p,_:=strconv.Atoi(q.Get("page"));ps,_:=strconv.Atoi(q.Get("page_size"));st:=-1;if q.Get("status")!=""{st,_=strconv.Atoi(q.Get("status"))};v,e:=ctx.Drama.AdminListDrama(r.Context(),&pb.AdminListDramaRequest{Keyword:q.Get("keyword"),Country:q.Get("country"),Language:q.Get("language"),Status:int32(st),Page:int32(p),PageSize:int32(ps)});write(w,v,e)}},
  {Method:http.MethodPut,Path:"/api/v1/admin/dramas/:id",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.UpdateDramaRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.UpdateDrama(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas/:id/status",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.SetDramaStatusRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.SetDramaStatus(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas/:id/episodes",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.CreateEpisodeRequest;if readJSON(w,r,&x)!=nil{return};x.DramaId=id;v,e:=ctx.Drama.CreateEpisode(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/dramas/:id/episodes",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};st:=-1;if r.URL.Query().Get("status")!=""{st,_=strconv.Atoi(r.URL.Query().Get("status"))};v,e:=ctx.Drama.ListEpisode(r.Context(),&pb.ListEpisodeRequest{DramaId:id,Status:int32(st)});write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/orders",Handler:func(w http.ResponseWriter,r *http.Request){q:=r.URL.Query();page,_:=strconv.Atoi(q.Get("page"));ps,_:=strconv.Atoi(q.Get("page_size"));v,e:=ctx.Payment.ListOrders(r.Context(),&paymentpb.ListOrdersRequest{Page:int32(page),PageSize:int32(ps)});write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/episodes/:id/status",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/episodes/");if !ok{http.Error(w,"invalid id",400);return};var x pb.SetEpisodeStatusRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.SetEpisodeStatus(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/episodes/:id/upload",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/episodes/");if !ok{http.Error(w,"invalid id",400);return};var x mediapb.CreateUploadRequest;if readJSON(w,r,&x)!=nil{return};x.EpisodeId=id;v,e:=ctx.Media.CreateUpload(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/episodes/:id/upload/complete",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/episodes/");if !ok{http.Error(w,"invalid id",400);return};var x mediapb.CompleteUploadRequest;if readJSON(w,r,&x)!=nil{return};x.EpisodeId=id;v,e:=ctx.Media.CompleteUpload(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/settings/:group",Handler:func(w http.ResponseWriter,r *http.Request){group:=pathParam(r,"/api/v1/admin/settings/");v,e:=readSettings(r,ctx.DB,group);write(w,v,e)}},
  {Method:http.MethodPut,Path:"/api/v1/admin/settings/:group",Handler:func(w http.ResponseWriter,r *http.Request){group:=r.PathValue("group");var values map[string]any;if readJSON(w,r,&values)!=nil{return};v,e:=saveSettings(r,ctx.DB,group,values);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/settings/:group/test",Handler:func(w http.ResponseWriter,r *http.Request){group:=r.PathValue("group");var x struct{Provider string `json:"provider"`};if readJSON(w,r,&x)!=nil{return};v,e:=testSettings(r,ctx.DB,group,x.Provider);write(w,v,e)}},
 })
}
func readSettings(r *http.Request,db *sql.DB,group string)(map[string]any,error){
 rows,e:=db.QueryContext(r.Context(),"SELECT setting_key,setting_value,is_secret FROM system_settings WHERE setting_group=?",group);if e!=nil{return nil,e};defer rows.Close()
 out:=map[string]any{};for rows.Next(){var k,v string;var secret int;if e=rows.Scan(&k,&v,&secret);e!=nil{return nil,e};if secret==1{out[k]="********"}else{out[k]=v}};return out,rows.Err()
}
func saveSettings(r *http.Request,db *sql.DB,group string,values map[string]any)(map[string]any,error){
 allowed:=map[string]bool{"payment":true,"storage":true};if !allowed[group]{return nil,errors.New("unsupported setting group")}
 secrets:=map[string]bool{"stripe_secret_key":true,"stripe_webhook_secret":true,"paypal_client_secret":true,"oss_secret_key":true,"s3_secret_key":true}
 for k,v:=range values{sv:=toString(v);if secrets[k]&&(sv==""||sv=="********"){continue};if secrets[k]&&sv!=""{enc,e:=settings.Encrypt(sv);if e!=nil{return nil,e};sv=enc};if secrets[k]&&sv==""{continue};_,e:=db.ExecContext(r.Context(),"INSERT INTO system_settings(setting_group,setting_key,setting_value,is_secret) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE setting_value=VALUES(setting_value),is_secret=VALUES(is_secret)",group,k,sv,boolInt(secrets[k]));if e!=nil{return nil,e}}
 return readSettings(r,db,group)
}
func testSettings(r *http.Request,db *sql.DB,group,provider string)(map[string]any,error){
 if group=="payment"&&(provider=="stripe"||provider=="paypal"){var key string;if provider=="stripe"{key="stripe_secret_key"}else{key="paypal_client_secret"};var v string;if e:=db.QueryRowContext(r.Context(),"SELECT setting_value FROM system_settings WHERE setting_group=? AND setting_key=?",group,key).Scan(&v);e!=nil{return nil,e};if _,e:=settings.Decrypt(v);e!=nil{return nil,errors.New("支付密钥尚未正确配置")};return map[string]any{"ok":true,"message":"配置已保存并可解密，真实支付请求将在支付时使用"},nil}
 if group=="storage"&&(provider=="oss"||provider=="s3"){return map[string]any{"ok":true,"message":"存储参数已保存；上传链路接入 media-rpc 后执行实际对象存储连通性测试"},nil}
 return nil,errors.New("unsupported provider")
}
func boolInt(v bool)int{if v{return 1};return 0}
func toString(v any)string{if s,ok:=v.(string);ok{return s};b,_:=json.Marshal(v);return string(b)}
func pathParam(r *http.Request,prefix string)string{return strings.Trim(strings.TrimPrefix(r.URL.Path,prefix),"/")}
func pathID(r *http.Request,prefix string)(int64,bool){x:=strings.TrimPrefix(r.URL.Path,prefix);x=strings.TrimSuffix(x,"/episodes");id,e:=strconv.ParseInt(x,10,64);return id,e==nil&&id>0}
func readJSON(w http.ResponseWriter,r *http.Request,v any)error{defer r.Body.Close();if err:=json.NewDecoder(r.Body).Decode(v);err!=nil{http.Error(w,err.Error(),400);return err};return nil}
func write(w http.ResponseWriter,v any,e error){w.Header().Set("Content-Type","application/json");if e!=nil{http.Error(w,e.Error(),500);return};json.NewEncoder(w).Encode(v)}
