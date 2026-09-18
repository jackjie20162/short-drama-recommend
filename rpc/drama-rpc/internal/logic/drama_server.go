package logic
import("context";"errors";"short-drama-recommend/rpc/drama-rpc/internal/svc";"short-drama-recommend/rpc/drama-rpc/pb")
type DramaServer struct{pb.UnimplementedDramaServiceServer;svcCtx *svc.ServiceContext}
func NewDramaServer(s *svc.ServiceContext)*DramaServer{return &DramaServer{svcCtx:s}}
func(s *DramaServer)GetDrama(ctx context.Context,req *pb.GetDramaRequest)(*pb.GetDramaResponse,error){if req.GetId()<=0{return nil,errors.New("drama id is required")};d,e:=s.svcCtx.DramaRepo.GetByID(ctx,uint64(req.GetId()));if e!=nil{return nil,e};return &pb.GetDramaResponse{Drama:toPBDrama(d)},nil}
func(s *DramaServer)ListEpisodes(ctx context.Context,req *pb.ListEpisodesRequest)(*pb.ListEpisodesResponse,error){
 if req.GetDramaId()<=0{return nil,errors.New("drama id is required")}
 es,e:=s.svcCtx.EpisodeRepo.ListByDrama(ctx,uint64(req.GetDramaId()),func()*int8{v:=int8(1);return &v}());if e!=nil{return nil,e}
 o:=&pb.ListEpisodesResponse{Items:make([]*pb.Episode,0,len(es))}
 unlocked:=false
 if req.GetUserId()>0{_ = s.svcCtx.DB.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM user_entitlements WHERE user_id=? AND drama_id=?)",req.GetUserId(),req.GetDramaId()).Scan(&unlocked)}
 for _,x:=range es{p:=toPBEpisode(x);p.Unlocked=!x.IsPaid||unlocked;if x.IsPaid&&!p.Unlocked{p.VideoUrl=""};o.Items=append(o.Items,p)}
 return o,nil
}
func(s *DramaServer)ListDrama(ctx context.Context,req *pb.ListDramaRequest)(*pb.ListDramaResponse,error){p:=int(req.GetPage());if p<1{p=1};n:=int(req.GetPageSize());if n<=0{n=20};if n>100{n=100};ds,e:=s.svcCtx.DramaRepo.ListPublished(ctx,req.GetCountry(),req.GetLanguage(),n,(p-1)*n);if e!=nil{return nil,e};o:=&pb.ListDramaResponse{Items:make([]*pb.Drama,0,len(ds))};for _,d:=range ds{o.Items=append(o.Items,toPBDrama(d))};return o,nil}
