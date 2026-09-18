package main

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"short-drama-recommend/api/drama-api/internal/handler"
	"short-drama-recommend/api/drama-api/internal/svc"
)

func registerRoutes(server *rest.Server, svcCtx *svc.ServiceContext) {
	h := handler.NewHandler(svcCtx)
	server.AddRoutes([]rest.Route{
		{Method:http.MethodGet, Path:"/api/v1/dramas/:id", Handler:h.GetDrama},
		{Method:http.MethodGet, Path:"/api/v1/dramas/:id/episodes", Handler:h.ListEpisodes},
		{Method:http.MethodGet, Path:"/api/v1/dramas", Handler:h.ListDrama},
		{Method:http.MethodGet, Path:"/api/v1/feed", Handler:h.GetFeed},
		{Method:http.MethodPost, Path:"/api/v1/behaviors", Handler:h.RecordBehavior},
		{Method:http.MethodPost, Path:"/api/v1/payments/orders", Handler:h.CreatePaymentOrder},
		{Method:http.MethodPost, Path:"/api/v1/payments/webhook", Handler:h.PaymentWebhook},
		{Method:http.MethodPost, Path:"/api/v1/payments/capture", Handler:h.CapturePayment},
		{Method:http.MethodGet, Path:"/api/v1/payments/orders/:id", Handler:h.GetPaymentOrder},
	})
}
