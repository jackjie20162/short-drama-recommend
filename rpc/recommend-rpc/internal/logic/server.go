package logic

import (
	"context"
	"sort"

	"short-drama-recommend/rpc/recommend-rpc/internal/svc"
	"short-drama-recommend/rpc/recommend-rpc/pb"
)

type RecommendServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedRecommendServiceServer
}

func NewRecommendServer(svcCtx *svc.ServiceContext) *RecommendServer {
	return &RecommendServer{svcCtx: svcCtx}
}

type scoredItem struct {
	item *pb.FeedItem
	score float32
}

func (s *RecommendServer) GetFeed(ctx context.Context, in *pb.FeedRequest) (*pb.FeedResponse, error) {
	size := int(in.GetPageSize())
	if size <= 0 { size = 20 }
	if size > 50 { size = 50 }

	// V1 uses a deliberately explainable ranking:
	// language/country match + popularity + stable freshness ordering.
	candidates, err := s.svcCtx.DramaRepo.ListPublished(ctx, in.GetCountry(), in.GetLanguage(), 100, 0)
	if err != nil { return nil, err }

	items := make([]scoredItem, 0, len(candidates))
	for _, d := range candidates {
		score := float32(10)
		reasons := []string{"published"}
		if in.GetLanguage() != "" && d.Language == in.GetLanguage() {
			score += 30
			reasons = append(reasons, "language_match")
		}
		if in.GetCountry() != "" && d.Country == in.GetCountry() {
			score += 20
			reasons = append(reasons, "country_match")
		}
		item := &pb.FeedItem{DramaId:int64(d.ID), Title:d.Title, Cover:d.Cover, Score:score, Reasons:reasons}
		items = append(items, scoredItem{item:item, score:score})
	}

	sort.SliceStable(items, func(i,j int) bool {
		if items[i].score == items[j].score { return items[i].item.DramaId > items[j].item.DramaId }
		return items[i].score > items[j].score
	})
	if len(items) > size { items = items[:size] }

	out := make([]*pb.FeedItem, 0, len(items))
	for _, x := range items { out = append(out, x.item) }
	return &pb.FeedResponse{Items:out}, nil
}
