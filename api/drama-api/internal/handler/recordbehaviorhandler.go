// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"short-drama-recommend/api/drama-api/internal/logic"
	"short-drama-recommend/api/drama-api/internal/svc"
	"short-drama-recommend/api/drama-api/internal/types"
)

func RecordBehaviorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecordBehaviorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRecordBehaviorLogic(r.Context(), svcCtx)
		resp, err := l.RecordBehavior(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
