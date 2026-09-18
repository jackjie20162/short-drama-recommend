package logic
import("context";"short-drama-recommend/rpc/user-rpc/internal/svc";"short-drama-recommend/rpc/user-rpc/pb")
type UserServer struct{pb.UnimplementedUserServiceServer;svcCtx *svc.ServiceContext}
func NewUserServer(s *svc.ServiceContext)*UserServer{return &UserServer{svcCtx:s}}
func(s *UserServer)GetUser(ctx context.Context,r *pb.GetUserRequest)(*pb.GetUserResponse,error){id:=r.GetId();if id<=0{id=1};return &pb.GetUserResponse{User:&pb.User{Id:id,Country:"US",Language:"en",Locale:"en-US",Timezone:"UTC"}},nil}
