package logic
import("context";"short-drama-recommend/rpc/recommend-rpc/internal/svc";"short-drama-recommend/rpc/recommend-rpc/pb")
type RecommendServer struct{pb.UnimplementedRecommendServiceServer;svcCtx *svc.ServiceContext}
func NewRecommendServer(s *svc.ServiceContext)*RecommendServer{return &RecommendServer{svcCtx:s}}
func(s *RecommendServer)GetFeed(ctx context.Context,r *pb.FeedRequest)(*pb.FeedResponse,error){n:=int(r.GetPageSize());if n<=0{n=20};if n>100{n=100};ds,e:=s.svcCtx.DramaRepo.ListPublished(ctx,r.GetCountry(),r.GetLanguage(),n,0);if e!=nil{return nil,e};o:=&pb.FeedResponse{Items:make([]*pb.FeedItem,0,len(ds))};for i,d:=range ds{o.Items=append(o.Items,&pb.FeedItem{DramaId:int64(d.ID),Title:d.Title,Cover:d.Cover,Score:float64(len(ds)-i),Reasons:[]string{"language_match","country_match"}})};return o,nil}
