package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"short-drama-recommend/recdemo/internal/model"
	"short-drama-recommend/recdemo/internal/rank"
)

type Config struct {
	Address string
	Index string
}

type Client struct {
	address string
	index string
	ranker *rank.Ranker
	http *http.Client
}

func New(c Config, r *rank.Ranker) *Client {
	return &Client{address:strings.TrimRight(c.Address,"/"), index:c.Index, ranker:r, http:&http.Client{Timeout:1500*time.Millisecond}}
}

func (c *Client) Search(ctx context.Context, req model.Query, size int) ([]model.Drama, error) {
	if c.address == "" || c.index == "" { return nil, fmt.Errorf("elasticsearch address/index is empty") }
	body, err := json.Marshal(FunctionScoreQuery(req,size))
	if err != nil { return nil, err }
	reqHTTP, err := http.NewRequestWithContext(ctx,http.MethodPost,c.address+"/"+c.index+"/_search",bytes.NewReader(body))
	if err != nil { return nil, err }
	reqHTTP.Header.Set("Content-Type","application/json")
	resp, err := c.http.Do(reqHTTP)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body,2048))
		return nil, fmt.Errorf("elasticsearch search status=%d body=%s",resp.StatusCode,strings.TrimSpace(string(b)))
	}
	var result struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
				Source model.Drama `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil { return nil, err }
	items := make([]model.Drama,0,len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		d := hit.Source
		if d.ID == 0 && hit.ID != "" {
			if id, err := strconv.ParseInt(hit.ID,10,64); err == nil { d.ID=id }
		}
		items=append(items,d)
	}
	return items,nil
}

func FunctionScoreQuery(req model.Query, size int) map[string]any {
	filters := []map[string]any{{"term":map[string]any{"status":"published"}}}
	if req.Region != "" && req.Region != "0" { filters=append(filters,map[string]any{"term":map[string]any{"region":req.Region}}) }
	if req.Lang != "" && req.Lang != "0" { filters=append(filters,map[string]any{"term":map[string]any{"lang":req.Lang}}) }
	boolQuery := map[string]any{"filter":filters}
	if len(req.SeenIDs)>0 { boolQuery["must_not"]=[]map[string]any{{"terms":map[string]any{"id":req.SeenIDs}}} }
	return map[string]any{
		"size":size,
		"_source":[]string{"id","title","cover","region","lang","status","popularity","completion","pay_rate","published_at"},
		"query":map[string]any{"function_score":map[string]any{
			"query":map[string]any{"bool":boolQuery},
			"functions":[]any{
				map[string]any{"field_value_factor":map[string]any{"field":"popularity","factor":0.35,"modifier":"sqrt","missing":0}},
				map[string]any{"field_value_factor":map[string]any{"field":"completion","factor":0.30,"modifier":"sqrt","missing":0}},
				map[string]any{"field_value_factor":map[string]any{"field":"pay_rate","factor":0.25,"modifier":"sqrt","missing":0}},
				map[string]any{"gauss":map[string]any{"published_at":map[string]any{"origin":"now","scale":"30d","decay":0.5}}},
			},
			"score_mode":"sum","boost_mode":"sum",
		}},
		"sort":[]any{map[string]any{"_score":"desc"},map[string]any{"id":"desc"}},
	}
}
