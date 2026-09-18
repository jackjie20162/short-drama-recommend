package logic

import (
	"context"
	"errors"

	"short-drama-recommend/internal/model"
	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"
)

type DramaServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedDramaServiceServer
}

func NewDramaServer(svcCtx *svc.ServiceContext) *DramaServer { return &DramaServer{svcCtx: svcCtx} }

func (s *DramaServer) GetDrama(ctx context.Context, in *pb.GetDramaRequest) (*pb.GetDramaResponse, error) {
	if in.GetId() <= 0 {
		return nil, errors.New("id is required")
	}
	drama, err := s.svcCtx.DramaRepo.GetByID(ctx, uint64(in.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.GetDramaResponse{Drama: toProto(drama)}, nil
}

func (s *DramaServer) ListDrama(ctx context.Context, in *pb.ListDramaRequest) (*pb.ListDramaResponse, error) {
	pageSize := int(in.GetPageSize())
	if pageSize <= 0 { pageSize = 20 }
	if pageSize > 100 { pageSize = 100 }
	page := int(in.GetPage())
	if page <= 0 { page = 1 }
	offset := (page - 1) * pageSize

	items, err := s.svcCtx.DramaRepo.ListPublished(ctx, in.GetCountry(), in.GetLanguage(), pageSize, offset)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.Drama, 0, len(items))
	for _, item := range items {
		result = append(result, toProto(item))
	}
	return &pb.ListDramaResponse{Items: result}, nil
}

func toProto(d *model.Drama) *pb.Drama {
	return &pb.Drama{
		Id:int64(d.ID), Title:d.Title, Description:d.Description, Cover:d.Cover,
		Country:d.Country, Language:d.Language, TotalEpisodes:int32(d.TotalEpisodes), IsPaid:d.IsPaid,
	}
}
