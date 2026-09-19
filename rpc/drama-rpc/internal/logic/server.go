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

func NewDramaServer(svcCtx *svc.ServiceContext) *DramaServer {
	return &DramaServer{svcCtx: svcCtx}
}

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
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	page := int(in.GetPage())
	if page <= 0 {
		page = 1
	}

	items, err := s.svcCtx.DramaRepo.ListPublished(ctx, in.GetCountry(), in.GetLanguage(), pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Drama, 0, len(items))
	for _, item := range items {
		result = append(result, toProto(item))
	}
	return &pb.ListDramaResponse{Items: result}, nil
}

func (s *DramaServer) ListEpisodes(ctx context.Context, in *pb.ListEpisodesRequest) (*pb.ListEpisodesResponse, error) {
	if in.GetDramaId() <= 0 {
		return nil, errors.New("drama id is required")
	}

	status := int8(1)
	items, err := s.svcCtx.EpisodeRepo.ListByDrama(ctx, uint64(in.GetDramaId()), &status)
	if err != nil {
		return nil, err
	}

	unlocked := false
	if in.GetUserId() > 0 {
		_ = s.svcCtx.DB.QueryRowContext(
			ctx,
			"SELECT EXISTS(SELECT 1 FROM user_entitlements WHERE user_id=? AND drama_id=?)",
			in.GetUserId(),
			in.GetDramaId(),
		).Scan(&unlocked)
	}

	result := make([]*pb.DramaEpisode, 0, len(items))
	for _, item := range items {
		episode := &pb.DramaEpisode{
			Id:              int64(item.ID),
			DramaId:         int64(item.DramaID),
			EpisodeNo:       int32(item.EpisodeNo),
			Title:           item.Title,
			Description:     item.Description,
			DurationSeconds: int32(item.DurationSeconds),
			VideoUrl:        item.VideoURL,
			PosterUrl:       item.PosterURL,
			IsPaid:          item.IsPaid,
			Status:          int32(item.Status),
			Unlocked:        !item.IsPaid || unlocked,
		}
		if item.IsPaid && !episode.Unlocked {
			episode.VideoUrl = ""
		}
		result = append(result, episode)
	}

	return &pb.ListEpisodesResponse{Items: result}, nil
}

func toProto(d *model.Drama) *pb.Drama {
	return &pb.Drama{
		Id:             int64(d.ID),
		Title:          d.Title,
		Subtitle:       d.Subtitle,
		Description:    d.Description,
		Cover:          d.Cover,
		Country:        d.Country,
		Language:       d.Language,
		Genres:         d.Genres,
		Tags:            d.Tags,
		TotalEpisodes:  int32(d.TotalEpisodes),
		IsPaid:         d.IsPaid,
		PriceCents:     int64(d.PriceCents),
		Currency:       d.Currency,
		Popularity:     d.Popularity,
		CompletionRate: d.CompletionRate,
		PayRate:        d.PayRate,
		PublishedAt:    d.PublishedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
