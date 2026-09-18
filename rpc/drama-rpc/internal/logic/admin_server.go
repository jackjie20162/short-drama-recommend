package logic

import (
 "context"
 "errors"
 "short-drama-recommend/internal/model"
 "short-drama-recommend/rpc/drama-rpc/internal/svc"
 "short-drama-recommend/rpc/drama-rpc/pb"
)

type DramaAdminServer struct {
 svcCtx *svc.ServiceContext
 pb.UnimplementedDramaAdminServiceServer
}

func NewDramaAdminServer(svcCtx *svc.ServiceContext)*DramaAdminServer{return &DramaAdminServer{svcCtx:svcCtx}}

func(s *DramaAdminServer)CreateDrama(ctx context.Context,in *pb.CreateDramaRequest)(*pb.CreateDramaResponse,error){
 if in.GetTitle()==""{return nil,errors.New("title is required")}
 d:=&model.Drama{Title:in.GetTitle(),Description:in.GetDescription(),Cover:in.GetCover(),Country:in.GetCountry(),Language:in.GetLanguage(),TotalEpisodes:uint32(max32(in.GetTotalEpisodes())),IsPaid:in.GetIsPaid()}
 if err:=s.svcCtx.DramaAdminRepo.Create(ctx,d);err!=nil{return nil,err}
 return &pb.CreateDramaResponse{Drama:toProto(d)},nil
}
func(s *DramaAdminServer)UpdateDrama(ctx context.Context,in *pb.UpdateDramaRequest)(*pb.UpdateDramaResponse,error){
 if in.GetId()<=0{return nil,errors.New("id is required")}
 d:=&model.Drama{ID:uint64(in.GetId()),Title:in.GetTitle(),Description:in.GetDescription(),Cover:in.GetCover(),Country:in.GetCountry(),Language:in.GetLanguage(),TotalEpisodes:uint32(max32(in.GetTotalEpisodes())),IsPaid:in.GetIsPaid()}
 if err:=s.svcCtx.DramaAdminRepo.Update(ctx,d);err!=nil{return nil,err}
 return &pb.UpdateDramaResponse{Drama:toProto(d)},nil
}
func(s *DramaAdminServer)SetDramaStatus(ctx context.Context,in *pb.SetDramaStatusRequest)(*pb.SetDramaStatusResponse,error){
 if in.GetId()<=0{return nil,errors.New("id is required")};if in.GetStatus()<0||in.GetStatus()>2{return nil,errors.New("status must be 0, 1 or 2")}
 if err:=s.svcCtx.DramaAdminRepo.SetStatus(ctx,uint64(in.GetId()),int8(in.GetStatus()));err!=nil{return nil,err};return &pb.SetDramaStatusResponse{Success:true},nil
}
func(s *DramaAdminServer)AdminListDrama(ctx context.Context,in *pb.AdminListDramaRequest)(*pb.AdminListDramaResponse,error){
 size:=int(in.GetPageSize());if size<=0{size=20};if size>100{size=100};page:=int(in.GetPage());if page<=0{page=1}
 var status *int8;if in.GetStatus()>=0&&in.GetStatus()<=2{v:=int8(in.GetStatus());status=&v}
 items,total,err:=s.svcCtx.DramaAdminRepo.List(ctx,in.GetKeyword(),in.GetCountry(),in.GetLanguage(),status,size,(page-1)*size);if err!=nil{return nil,err}
 out:=make([]*pb.Drama,0,len(items));for _,d:=range items{out=append(out,toProto(d))}
 return &pb.AdminListDramaResponse{Items:out,Total:total},nil
}
func(s *DramaAdminServer)CreateEpisode(ctx context.Context,in *pb.CreateEpisodeRequest)(*pb.CreateEpisodeResponse,error){
 if in.GetDramaId()<=0||in.GetEpisodeNo()<=0{return nil,errors.New("drama_id and episode_no are required")}
 e:=&model.Episode{DramaID:uint64(in.GetDramaId()),EpisodeNo:uint32(in.GetEpisodeNo()),Title:in.GetTitle(),DurationSeconds:uint32(max32(in.GetDurationSeconds())),VideoURL:in.GetVideoUrl(),PosterURL:in.GetPosterUrl(),IsPaid:in.GetIsPaid()}
 if err:=s.svcCtx.EpisodeRepo.Create(ctx,e);err!=nil{return nil,err};return &pb.CreateEpisodeResponse{Episode:episodeToProto(e)},nil
}
func(s *DramaAdminServer)ListEpisode(ctx context.Context,in *pb.ListEpisodeRequest)(*pb.ListEpisodeResponse,error){
 if in.GetDramaId()<=0{return nil,errors.New("drama_id is required")}
 var status *int8;if in.GetStatus()>=0&&in.GetStatus()<=2{v:=int8(in.GetStatus());status=&v}
 items,err:=s.svcCtx.EpisodeRepo.ListByDrama(ctx,uint64(in.GetDramaId()),status);if err!=nil{return nil,err}
 out:=make([]*pb.Episode,0,len(items));for _,e:=range items{out=append(out,episodeToProto(e))}
 return &pb.ListEpisodeResponse{Items:out},nil
}
func(s *DramaAdminServer)SetEpisodeStatus(ctx context.Context,in *pb.SetEpisodeStatusRequest)(*pb.SetEpisodeStatusResponse,error){
 if in.GetId()<=0{return nil,errors.New("id is required")};if in.GetStatus()<0||in.GetStatus()>2{return nil,errors.New("status must be 0, 1 or 2")}
 if err:=s.svcCtx.EpisodeRepo.SetStatus(ctx,uint64(in.GetId()),int8(in.GetStatus()));err!=nil{return nil,err};return &pb.SetEpisodeStatusResponse{Success:true},nil
}
func episodeToProto(e *model.Episode)*pb.Episode{return &pb.Episode{Id:int64(e.ID),DramaId:int64(e.DramaID),EpisodeNo:int32(e.EpisodeNo),Title:e.Title,DurationSeconds:int32(e.DurationSeconds),VideoUrl:e.VideoURL,PosterUrl:e.PosterURL,IsPaid:e.IsPaid,Status:int32(e.Status)}}
func max32(v int32)int32{if v<0{return 0};return v}
