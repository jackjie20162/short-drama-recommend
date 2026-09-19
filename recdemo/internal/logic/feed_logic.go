package logic

import (
	"context"
	"sort"
	"time"

	"short-drama-recommend/recdemo/internal/es"
	"short-drama-recommend/recdemo/internal/model"
	"short-drama-recommend/recdemo/internal/rank"
)

type FeedRequest struct { UserID, Region, Lang string; Size int }
type FeedItem struct { Drama model.Drama `json:"drama"`; Score float64 `json:"score"`; Reasons []string `json:"reasons"` }
type FeedResponse struct { Source string `json:"source"`; Items []FeedItem `json:"items"` }
type FeedLogic struct { seeds []model.Drama; es *es.Client; ranker *rank.Ranker }

func NewFeedLogic(seeds []model.Drama, client *es.Client, ranker *rank.Ranker) *FeedLogic { return &FeedLogic{seeds:seeds,es:client,ranker:ranker} }

func (l *FeedLogic) Feed(ctx context.Context, req FeedRequest) (FeedResponse,error) {
	_ = ctx
	now:=time.Now().Unix()
	candidates:=make([]model.Drama,0,len(l.seeds))
	for _,d:=range l.seeds {
		if d.Status!="published" { continue }
		if req.Region!="0" && req.Region!="" && d.Region!=req.Region { continue }
		if req.Lang!="0" && req.Lang!="" && d.Lang!=req.Lang { continue }
		if req.UserID!="" && d.SeenBy[req.UserID] { continue }
		candidates=append(candidates,d)
	}
	type scored struct{ d model.Drama; score float64; reasons []string }
	out:=make([]scored,0,len(candidates))
	for _,d:=range candidates {
		s:=l.ranker.Score(rank.Features{Popularity:d.Popularity,Completion:d.Completion,PayRate:d.PayRate,PublishedAt:d.PublishedAt,RegionMatch:req.Region!=""&&req.Region!="0"&&d.Region==req.Region,LangMatch:req.Lang!=""&&req.Lang!="0"&&d.Lang==req.Lang},now)
		reasons:=[]string{}
		if req.Region!=""&&req.Region!="0"&&d.Region==req.Region { reasons=append(reasons,"region_match") }
		if req.Lang!=""&&req.Lang!="0"&&d.Lang==req.Lang { reasons=append(reasons,"language_match") }
		if d.Completion>=0.8 { reasons=append(reasons,"high_completion") }
		if d.PayRate>=0.35 { reasons=append(reasons,"high_pay_rate") }
		out=append(out,scored{d,s,reasons})
	}
	sort.SliceStable(out,func(i,j int)bool{if out[i].score==out[j].score{return out[i].d.ID>out[j].d.ID};return out[i].score>out[j].score})
	if len(out)>req.Size {out=out[:req.Size]}
	items:=make([]FeedItem,0,len(out));for _,x:=range out{items=append(items,FeedItem{Drama:x.d,Score:x.score,Reasons:x.reasons})}
	return FeedResponse{Source:"memory",Items:items},nil
}
