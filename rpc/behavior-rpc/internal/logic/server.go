package logic

import (
	"context"
	"errors"
	"time"

	"short-drama-recommend/internal/model"
	"short-drama-recommend/rpc/behavior-rpc/internal/svc"
	"short-drama-recommend/rpc/behavior-rpc/pb"
)

type BehaviorServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedBehaviorServiceServer
}

func NewBehaviorServer(svcCtx *svc.ServiceContext) *BehaviorServer { return &BehaviorServer{svcCtx: svcCtx} }

func (s *BehaviorServer) RecordEvent(ctx context.Context, in *pb.RecordEventRequest) (*pb.RecordEventResponse, error) {
	if in.GetDramaId() <= 0 || in.GetEventType() == pb.EventType_EVENT_TYPE_UNSPECIFIED {
		return nil, errors.New("drama_id and event_type are required")
	}
	event := &model.BehaviorEvent{
		UserID:uint64(in.GetUserId()), DramaID:uint64(in.GetDramaId()), EpisodeID:uint64(in.GetEpisodeId()),
		EventType:in.GetEventType().String(), WatchSeconds:uint32(maxInt32(in.GetWatchSeconds())),
		DurationSeconds:uint32(maxInt32(in.GetDurationSeconds())), Country:in.GetCountry(), Language:in.GetLanguage(),
		Device:in.GetDevice(), EventAt:time.Now().UTC(),
	}
	if err := s.svcCtx.BehaviorRepo.Record(ctx, event); err != nil { return nil, err }
	if err := s.svcCtx.Redis.IncrDramaEvent(ctx, in.GetDramaId(), in.GetEventType().String()); err != nil { return nil, err }
	if in.GetUserId() > 0 && in.GetDramaId() > 0 {
		if err := s.svcCtx.Redis.SetUserRecentDrama(ctx, in.GetUserId(), in.GetDramaId()); err != nil { return nil, err }
	}
	return &pb.RecordEventResponse{Success:true}, nil
}

func maxInt32(v int32) int32 { if v < 0 { return 0 }; return v }
