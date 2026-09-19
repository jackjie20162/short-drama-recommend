package handler

import (
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"short-drama-recommend/recdemo/internal/logic"
)

type FeedHandler struct { svc *logic.FeedLogic }

func Register(server *rest.Server, svc *logic.FeedLogic) {
	h := &FeedHandler{svc: svc}
	server.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/rec/feed", Handler: h.Feed},
	})
}

func (h *FeedHandler) Feed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	size, _ := strconv.Atoi(q.Get("size"))
	if size <= 0 { size = 10 }
	if size > 50 { size = 50 }
	result, err := h.svc.Feed(r.Context(), logic.FeedRequest{
		UserID: q.Get("user_id"),
		Region: q.Get("region"),
		Lang: q.Get("lang"),
		Size: size,
	})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, result)
}
