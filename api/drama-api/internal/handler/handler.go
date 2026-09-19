package handler

import (
	"net/http"
	"io"
	"strconv"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"short-drama-recommend/api/drama-api/internal/svc"
	dramapb "short-drama-recommend/rpc/drama-rpc/pb"
	behaviorpb "short-drama-recommend/rpc/behavior-rpc/pb"
	recommendpb "short-drama-recommend/rpc/recommend-rpc/pb"
	"short-drama-recommend/pkg/auth"
	paymentpb "short-drama-recommend/rpc/payment-rpc/pb"
)

type Handler struct { svcCtx *svc.ServiceContext }

func NewHandler(svcCtx *svc.ServiceContext) *Handler { return &Handler{svcCtx: svcCtx} }

type listReq struct {
	Country string `form:"country,optional"`
	Language string `form:"language,optional"`
	Page int `form:"page,optional"`
	PageSize int `form:"page_size,optional"`
}
type feedReq struct {
	UserID int64 `form:"user_id,optional"`
	Country string `form:"country,optional"`
	Language string `form:"language,optional"`
	PageSize int `form:"page_size,optional"`
	Cursor string `form:"cursor,optional"`
}
type paymentReq struct {
 UserID int64 `json:"user_id"`
 DramaID int64 `json:"drama_id"`
 Provider string `json:"provider"`
 Currency string `json:"currency"`
 ReturnURL string `json:"return_url"`
 CancelURL string `json:"cancel_url"`
}
type captureReq struct { OrderID int64 `json:"order_id"` ProviderOrderID string `json:"provider_order_id"` }

type behaviorReq struct {
	UserID int64 `json:"user_id"`
	DramaID int64 `json:"drama_id"`
	EpisodeID int64 `json:"episode_id"`
	EventType string `json:"event_type"`
	WatchSeconds int32 `json:"watch_seconds"`
	DurationSeconds int32 `json:"duration_seconds"`
	Country string `json:"country"`
	Language string `json:"language"`
	Device string `json:"device"`
}

func (h *Handler) GetDrama(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(rest.PathValue(r.Context(), "id"), 10, 64)
	if err != nil || id <= 0 { httpx.Error(w, err); return }
	resp, err := h.svcCtx.Drama.GetDrama(r.Context(), &dramapb.GetDramaRequest{Id:id})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, resp)
}

func (h *Handler) ListEpisodes(w http.ResponseWriter,r *http.Request){
 q:=r.URL.Query();uid,_:=strconv.ParseInt(q.Get("user_id"),10,64);if id,e:=auth.Bearer(r.Header.Get("Authorization"));e==nil{uid=id}
 id,err:=strconv.ParseInt(rest.PathValue(r.Context(),"id"),10,64);if err!=nil||id<=0{httpx.Error(w,err);return}
 resp,err:=h.svcCtx.Drama.ListEpisodes(r.Context(),&dramapb.ListEpisodesRequest{DramaId:id,UserId:uid});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}

func (h *Handler) ListDrama(w http.ResponseWriter, r *http.Request) {
	var req listReq
	if err := httpx.Parse(r, &req); err != nil { httpx.Error(w, err); return }
	resp, err := h.svcCtx.Drama.ListDrama(r.Context(), &dramapb.ListDramaRequest{Country:req.Country, Language:req.Language, Page:int32(req.Page), PageSize:int32(req.PageSize)})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, resp)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	var req feedReq
	if err := httpx.Parse(r, &req); err != nil { httpx.Error(w, err); return }
	resp, err := h.svcCtx.Recommend.GetFeed(r.Context(), &recommendpb.FeedRequest{UserId:req.UserID, Country:req.Country, Language:req.Language, PageSize:int32(req.PageSize), Cursor:req.Cursor})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, resp)
}

func (h *Handler) RecordBehavior(w http.ResponseWriter, r *http.Request) {
	var req behaviorReq
	if err := httpx.Parse(r, &req); err != nil { httpx.Error(w, err); return }
	event, ok := behaviorpb.EventType_value[req.EventType]
	if !ok { httpx.Error(w, strconv.ErrSyntax); return }
	resp, err := h.svcCtx.Behavior.RecordEvent(r.Context(), &behaviorpb.RecordEventRequest{
		UserId:req.UserID, DramaId:req.DramaID, EpisodeId:req.EpisodeID, EventType:behaviorpb.EventType(event),
		WatchSeconds:req.WatchSeconds, DurationSeconds:req.DurationSeconds, Country:req.Country, Language:req.Language, Device:req.Device,
	})
	if err != nil { httpx.Error(w, err); return }
	httpx.OkJson(w, resp)
}


func (h *Handler) PaymentWebhook(w http.ResponseWriter,r *http.Request){
 body,err:=io.ReadAll(r.Body);if err!=nil{httpx.Error(w,err);return};defer r.Body.Close()
 provider:=r.Header.Get("X-Payment-Provider");if provider==""{provider=r.URL.Query().Get("provider")}
 req:=&paymentpb.WebhookRequest{Provider:provider,Payload:string(body),Signature:r.Header.Get("Stripe-Signature"),TransmissionId:r.Header.Get("PAYPAL-TRANSMISSION-ID"),TransmissionTime:r.Header.Get("PAYPAL-TRANSMISSION-TIME"),CertUrl:r.Header.Get("PAYPAL-CERT-URL"),AuthAlgo:r.Header.Get("PAYPAL-AUTH-ALGO"),TransmissionSig:r.Header.Get("PAYPAL-TRANSMISSION-SIG")}
 resp,err:=h.svcCtx.Payment.HandleWebhook(r.Context(),req);if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}

func (h *Handler) CreatePaymentOrder(w http.ResponseWriter,r *http.Request){
 uid,e:=auth.Bearer(r.Header.Get("Authorization"));if e!=nil{httpx.Error(w,e);return}
 var req paymentReq
 if err:=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return};req.UserID=uid
 provider:=paymentpb.Provider_STRIPE
 if req.Provider=="paypal"||req.Provider=="PAYPAL"{provider=paymentpb.Provider_PAYPAL}
 resp,err:=h.svcCtx.Payment.CreateOrder(r.Context(),&paymentpb.CreateOrderRequest{UserId:req.UserID,DramaId:req.DramaID,Provider:provider,Currency:req.Currency,ReturnUrl:req.ReturnURL,CancelUrl:req.CancelURL})
 if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}
func (h *Handler) CapturePayment(w http.ResponseWriter,r *http.Request){
 if _,e:=auth.Bearer(r.Header.Get("Authorization"));e!=nil{httpx.Error(w,e);return}
 var req captureReq
 if err:=httpx.Parse(r,&req);err!=nil{httpx.Error(w,err);return}
 resp,err:=h.svcCtx.Payment.CapturePayment(r.Context(),&paymentpb.CapturePaymentRequest{OrderId:req.OrderID,ProviderOrderId:req.ProviderOrderID,UserId:func()int64{id,_:=auth.Bearer(r.Header.Get("Authorization"));return id}()})
 if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}
func (h *Handler) ListPaymentOrders(w http.ResponseWriter,r *http.Request){
 uid,e:=auth.Bearer(r.Header.Get("Authorization"));if e!=nil{httpx.Error(w,e);return}
 q:=r.URL.Query();_ = q;strconv.ParseInt(q.Get("user_id"),10,64);page,_:=strconv.Atoi(q.Get("page"));size,_:=strconv.Atoi(q.Get("page_size"))
 resp,err:=h.svcCtx.Payment.ListOrders(r.Context(),&paymentpb.ListOrdersRequest{UserId:uid,Page:int32(page),PageSize:int32(size)});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}

func (h *Handler) GetPaymentOrder(w http.ResponseWriter,r *http.Request){
 id,err:=strconv.ParseInt(rest.PathValue(r.Context(),"id"),10,64);if err!=nil||id<=0{httpx.Error(w,err);return}
 uid,e:=auth.Bearer(r.Header.Get("Authorization"));if e!=nil{httpx.Error(w,e);return};resp,err:=h.svcCtx.Payment.GetOrder(r.Context(),&paymentpb.GetOrderRequest{OrderId:id,UserId:uid});if err!=nil{httpx.Error(w,err);return};httpx.OkJson(w,resp)
}
