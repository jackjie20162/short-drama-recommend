package logic

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"short-drama-recommend/rpc/recommend-rpc/internal/svc"
	"short-drama-recommend/rpc/recommend-rpc/pb"
)

type RecommendServer struct {
	pb.UnimplementedRecommendServiceServer
	svcCtx *svc.ServiceContext
}

func NewRecommendServer(s *svc.ServiceContext) *RecommendServer {
	return &RecommendServer{svcCtx: s}
}

func (s *RecommendServer) GetFeed(ctx context.Context, r *pb.FeedRequest) (*pb.FeedResponse, error) {
	if s.svcCtx == nil || s.svcCtx.DB == nil {
		return nil, errors.New("recommendation database is not configured")
	}
	size := int(r.GetPageSize())
	if size <= 0 { size = 20 }
	if size > 50 { size = 50 }
	offset, err := decodeCursor(r.GetCursor())
	if err != nil { return nil, err }

	const q = "SELECT d.id, d.title, d.cover, (" +
		"CASE WHEN ? <> '' AND d.country = ? THEN 8 ELSE 0 END + " +
		"CASE WHEN ? <> '' AND d.language = ? THEN 8 ELSE 0 END + " +
		"LEAST(LOG10(1 + COALESCE(p.interactions, 0)) * 3.0, 12.0) + " +
		"LEAST(LOG10(1 + COALESCE(p.watch_seconds, 0)) * 1.5, 8.0) + " +
		"CASE WHEN ? > 0 AND EXISTS (SELECT 1 FROM behavior_events ub WHERE ub.user_id = ? AND ub.drama_id = d.id LIMIT 1) THEN -6 ELSE 0 END + " +
		"CASE WHEN ? > 0 AND EXISTS (SELECT 1 FROM user_entitlements ue WHERE ue.user_id = ? AND ue.drama_id = d.id) THEN 2 ELSE 0 END" +
		") AS score, " +
		"CASE WHEN ? > 0 AND EXISTS (SELECT 1 FROM behavior_events ur WHERE ur.user_id = ? AND ur.drama_id = d.id LIMIT 1) THEN 1 ELSE 0 END AS seen, " +
		"COALESCE(p.interactions, 0) AS interactions, COALESCE(p.watch_seconds, 0) AS watch_seconds " +
		"FROM dramas d LEFT JOIN (SELECT drama_id, COUNT(*) AS interactions, SUM(watch_seconds) AS watch_seconds FROM behavior_events GROUP BY drama_id) p ON p.drama_id = d.id " +
		"WHERE d.status = 1 ORDER BY score DESC, d.id DESC LIMIT ? OFFSET ?"

	rows, err := s.svcCtx.DB.QueryContext(ctx, q,
		r.GetCountry(), r.GetCountry(), r.GetLanguage(), r.GetLanguage(),
		r.GetUserId(), r.GetUserId(), r.GetUserId(), r.GetUserId(),
		r.GetUserId(), r.GetUserId(), size+1, offset,
	)
	if err != nil { return nil, fmt.Errorf("query recommendation candidates: %w", err) }
	defer rows.Close()

	items := make([]*pb.FeedItem, 0, size)
	for rows.Next() {
		var id int64
		var title, cover string
		var score float64
		var seen int
		var interactions, watchSeconds int64
		if err := rows.Scan(&id, &title, &cover, &score, &seen, &interactions, &watchSeconds); err != nil {
			return nil, fmt.Errorf("scan recommendation candidate: %w", err)
		}
		reasons := make([]string, 0, 4)
		if r.GetCountry() != "" { reasons = append(reasons, "region_match") }
		if r.GetLanguage() != "" { reasons = append(reasons, "language_match") }
		if interactions > 0 { reasons = append(reasons, "popular") }
		if r.GetUserId() > 0 && seen == 0 { reasons = append(reasons, "not_watched") }
		if watchSeconds >= 600 { reasons = append(reasons, "high_watch_time") }
		if len(reasons) == 0 { reasons = append(reasons, "global_popular") }
		items = append(items, &pb.FeedItem{DramaId:id, Title:title, Cover:cover, Score:float32(score), Reasons:reasons})
		if len(items) == size { break }
	}
	if err := rows.Err(); err != nil { return nil, fmt.Errorf("iterate recommendation candidates: %w", err) }
	next := ""
	if len(items) == size { next = encodeCursor(offset + size) }
	return &pb.FeedResponse{Items:items, NextCursor:next}, nil
}

func encodeCursor(offset int) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

func decodeCursor(cursor string) (int, error) {
	if strings.TrimSpace(cursor) == "" { return 0, nil }
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil { return 0, errors.New("invalid recommendation cursor") }
	offset, err := strconv.Atoi(string(raw))
	if err != nil || offset < 0 { return 0, errors.New("invalid recommendation cursor") }
	return offset, nil
}
