package handler

import(
 "encoding/json"
 "net/http"
 "strconv"
 "strings"
 "github.com/zeromicro/go-zero/rest"
 "github.com/zeromicro/go-zero/zrpc"
 "short-drama-recommend/api/drama-admin-api/internal/config"
 "short-drama-recommend/rpc/drama-rpc/pb"
)
type ServiceContext struct{Drama pb.DramaAdminServiceClient}
func NewServiceContext(c config.Config)*ServiceContext{return &ServiceContext{Drama:pb.NewDramaAdminServiceClient(zrpc.MustNewClient(c.DramaRpc).Conn())}}

func RegisterRoutes(s *rest.Server,ctx *ServiceContext){
 s.AddRoutes([]rest.Route{
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas",Handler:func(w http.ResponseWriter,r *http.Request){var x pb.CreateDramaRequest;if readJSON(w,r,&x)!=nil{return};v,e:=ctx.Drama.CreateDrama(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/dramas",Handler:func(w http.ResponseWriter,r *http.Request){q:=r.URL.Query();p,_:=strconv.Atoi(q.Get("page"));ps,_:=strconv.Atoi(q.Get("page_size"));st:=-1;if q.Get("status")!=""{st,_=strconv.Atoi(q.Get("status"))};v,e:=ctx.Drama.AdminListDrama(r.Context(),&pb.AdminListDramaRequest{Keyword:q.Get("keyword"),Country:q.Get("country"),Language:q.Get("language"),Status:int32(st),Page:int32(p),PageSize:int32(ps)});write(w,v,e)}},
  {Method:http.MethodPut,Path:"/api/v1/admin/dramas/:id",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.UpdateDramaRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.UpdateDrama(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas/:id/status",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.SetDramaStatusRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.SetDramaStatus(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/dramas/:id/episodes",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};var x pb.CreateEpisodeRequest;if readJSON(w,r,&x)!=nil{return};x.DramaId=id;v,e:=ctx.Drama.CreateEpisode(r.Context(),&x);write(w,v,e)}},
  {Method:http.MethodGet,Path:"/api/v1/admin/dramas/:id/episodes",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/dramas/");if !ok{http.Error(w,"invalid id",400);return};st:=-1;if r.URL.Query().Get("status")!=""{st,_=strconv.Atoi(r.URL.Query().Get("status"))};v,e:=ctx.Drama.ListEpisode(r.Context(),&pb.ListEpisodeRequest{DramaId:id,Status:int32(st)});write(w,v,e)}},
  {Method:http.MethodPost,Path:"/api/v1/admin/episodes/:id/status",Handler:func(w http.ResponseWriter,r *http.Request){id,ok:=pathID(r,"/api/v1/admin/episodes/");if !ok{http.Error(w,"invalid id",400);return};var x pb.SetEpisodeStatusRequest;if readJSON(w,r,&x)!=nil{return};x.Id=id;v,e:=ctx.Drama.SetEpisodeStatus(r.Context(),&x);write(w,v,e)}},
 })
}
func pathID(r *http.Request,prefix string)(int64,bool){x:=strings.TrimPrefix(r.URL.Path,prefix);x=strings.TrimSuffix(x,"/episodes");id,e:=strconv.ParseInt(x,10,64);return id,e==nil&&id>0}
func readJSON(w http.ResponseWriter,r *http.Request,v any)error{defer r.Body.Close();if err:=json.NewDecoder(r.Body).Decode(v);err!=nil{http.Error(w,err.Error(),400);return err};return nil}
func write(w http.ResponseWriter,v any,e error){w.Header().Set("Content-Type","application/json");if e!=nil{http.Error(w,e.Error(),500);return};json.NewEncoder(w).Encode(v)}
