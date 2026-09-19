package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"short-drama-recommend/recdemo/internal/model"
	"short-drama-recommend/recdemo/internal/rank"
)

type Client struct { address, index string; ranker *rank.Ranker; http *http.Client }

func New(c struct{ Address string; Index string }, r *rank.Ranker) *Client {
	return &Client{address:strings.TrimRight(c.Address,"/"),index:c.Index,ranker:r,http:&http.Client{Timeout:1500*time.Millisecond}}
}

func (c *Client) Search(ctx context.Context, req model.Query, size int) ([]model.Drama, error) {
	_ = req; _ = size
	return nil, fmt.Errorf("ES adapter initialized; use memory fallback until ES document mapping is provisioned")
}

func FunctionScoreQuery(req model.Query, size int) map[string]any {
	filters := []map[string]any{
		{"term": map[string]any{"status": "published"}},
	}
	if req.Region != "" && req.Region != "0" { filters=append(filters,map[string]any{"term":map[string]any{"region":req.Region}}) }
	if req.Lang != "" && req.Lang != "0" { filters=append(filters,map[string]any{"term":map[string]any{"lang":req.Lang}}) }
	return map[string]any{
		"size":size,
		"query":map[string]any{"function_score":map[string]any{
			"query":map[string]any{"bool":map[string]any{"filter":filters}},
			"functions":[]any{
				map[string]any{"field_value_factor":map[string]any{"field":"popularity","factor":0.35,"modifier":"sqrt","missing":0}},
				map[string]any{"field_value_factor":map[string]any{"field":"completion","factor":0.30,"modifier":"sqrt","missing":0}},
				map[string]any{"field_value_factor":map[string]any{"field":"pay_rate","factor":0.25,"modifier":"sqrt","missing":0}},
			},
			"score_mode":"sum","boost_mode":"sum",
		}},
		"sort":[]any{map[string]any{"_score":"desc"},map[string]any{"_id":"desc"}},
	}
}
