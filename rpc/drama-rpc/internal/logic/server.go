package logic

import (
	"context"
	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"
)

type DramaServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedDramaServiceServer
}

func NewDramaServer(svcCtx *svc.ServiceContext) *DramaServer { return &DramaServer{svcCtx: svcCtx} }

func (s *DramaServer) GetDrama(ctx context.Context, in *pb.GetDramaRequest) (*pb.GetDramaResponse, error) {
	return &pb.GetDramaResponse{Drama: &pb.Drama{Id: in.Id}}, nil
}

func (s *DramaServer) ListDrama(ctx context.Context, in *pb.ListDramaRequest) (*pb.ListDramaResponse, error) {
	return &pb.ListDramaResponse{Items: []*pb.Drama{}}, nil
}
