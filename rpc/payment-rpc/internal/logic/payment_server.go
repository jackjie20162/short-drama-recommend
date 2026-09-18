package logic

import (
 "bytes"
 "context"
 "database/sql"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strconv"
 "strings"
 "time"
 "short-drama-recommend/rpc/payment-rpc/internal/svc"
 "short-drama-recommend/rpc/payment-rpc/pb"
)

type PaymentServer struct { pb.UnimplementedPaymentServiceServer; svcCtx *svc.ServiceContext }
func NewPaymentServer(svcCtx *svc.ServiceContext) *PaymentServer { return &PaymentServer{svcCtx:svcCtx} }

func (s *PaymentServer) CreateOrder(ctx context.Context,r *pb.CreateOrderRequest)(*pb.CreateOrderResponse,error){
 if r.GetUserId()<=0||r.GetDramaId()<=0{return nil,errors.New("user_id and drama_id are required")}
 var cents int64;var currency string;var paid bool
 if err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT price_cents,currency,is_paid FROM dramas WHERE id=? AND status=1",r.GetDramaId()).Scan(&cents,&currency,&paid);err!=nil{return nil,err}
 if !paid||cents<=0{return nil,errors.New("drama is not paid")};if r.GetCurrency()!=""{currency=strings.ToUpper(r.GetCurrency())}
 provider:=r.GetProvider();if provider==pb.Provider_PROVIDER_UNSPECIFIED{provider=pb.Provider_STRIPE}
 orderNo:=fmt.Sprintf("SD%d%d",time.Now().UnixNano(),r.GetUserId()%1000)
 res,err:=s.svcCtx.DB.ExecContext(ctx,"INSERT INTO orders(order_no,user_id,drama_id,provider,amount_cents,currency,status) VALUES(?,?,?,?,?,?,?)",orderNo,r.GetUserId(),r.GetDramaId(),providerName(provider),cents,currency,"PENDING");if err!=nil{return nil,err}
 orderID,err:=res.LastInsertId();if err!=nil{return nil,err}
 out:=&pb.CreateOrderResponse{OrderId:orderID,OrderNo:orderNo,Amount:fmt.Sprintf("%.2f",float64(cents)/100),Currency:currency,Status:"PENDING"}
 if provider==pb.Provider_STRIPE{
  id,secret,err:=s.createStripe(ctx,cents,currency,orderNo);if err!=nil{return nil,err};out.ProviderOrderId=id;out.ClientSecret=secret
  _,err=s.svcCtx.DB.ExecContext(ctx,"UPDATE orders SET provider_order_id=?,client_secret=? WHERE id=?",id,secret,orderID);if err!=nil{return nil,err}
 }else if provider==pb.Provider_PAYPAL{
  id,approve,err:=s.createPayPal(ctx,cents,currency,orderNo,r.GetReturnUrl(),r.GetCancelUrl());if err!=nil{return nil,err};out.ProviderOrderId=id;out.ApproveUrl=approve
  _,err=s.svcCtx.DB.ExecContext(ctx,"UPDATE orders SET provider_order_id=?,approve_url=? WHERE id=?",id,approve,orderID);if err!=nil{return nil,err}
 }else{return nil,errors.New("unsupported provider")}
 return out,nil
}

func (s *PaymentServer) createStripe(ctx context.Context,cents int64,currency,orderNo string)(string,string,error){
 if s.svcCtx.Config.Stripe.SecretKey==""{return "","",errors.New("STRIPE_SECRET_KEY is not configured")}
 form:=url.Values{};form.Set("amount",strconv.FormatInt(cents,10));form.Set("currency",strings.ToLower(currency));form.Set("metadata[order_no]",orderNo)
 req,_:=http.NewRequestWithContext(ctx,http.MethodPost,"https://api.stripe.com/v1/payment_intents",strings.NewReader(form.Encode()));req.SetBasicAuth(s.svcCtx.Config.Stripe.SecretKey,"");req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return "","",err};defer resp.Body.Close();body,_:=io.ReadAll(resp.Body);if resp.StatusCode/100!=2{return "","",fmt.Errorf("stripe: %s",body)}
 var x map[string]any;if err=json.Unmarshal(body,&x);err!=nil{return "","",err};id,_:=x["id"].(string);secret,_:=x["client_secret"].(string);if id==""||secret==""{return "","",errors.New("stripe response missing payment intent fields")};return id,secret,nil
}

func (s *PaymentServer) paypalBase()string{if strings.EqualFold(s.svcCtx.Config.Paypal.Env,"live"){return "https://api-m.paypal.com"};return "https://api-m.sandbox.paypal.com"}
func (s *PaymentServer) paypalToken(ctx context.Context)(string,error){
 if s.svcCtx.Config.Paypal.ClientID==""||s.svcCtx.Config.Paypal.ClientSecret==""{return "",errors.New("PayPal credentials are not configured")}
 req,_:=http.NewRequestWithContext(ctx,http.MethodPost,s.paypalBase()+"/v1/oauth2/token",strings.NewReader("grant_type=client_credentials"));req.SetBasicAuth(s.svcCtx.Config.Paypal.ClientID,s.svcCtx.Config.Paypal.ClientSecret);req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return "",err};defer resp.Body.Close();body,_:=io.ReadAll(resp.Body);if resp.StatusCode/100!=2{return "",fmt.Errorf("paypal token: %s",body)}
 var x map[string]any;if err=json.Unmarshal(body,&x);err!=nil{return "",err};token,_:=x["access_token"].(string);return token,nil
}
func (s *PaymentServer) createPayPal(ctx context.Context,cents int64,currency,orderNo,returnURL,cancelURL string)(string,string,error){
 token,err:=s.paypalToken(ctx);if err!=nil{return "","",err};amount:=fmt.Sprintf("%.2f",float64(cents)/100)
 payload:=map[string]any{"intent":"CAPTURE","purchase_units":[]any{map[string]any{"reference_id":orderNo,"custom_id":orderNo,"amount":map[string]string{"currency_code":strings.ToUpper(currency),"value":amount}}}}
 if returnURL!=""{payload["application_context"]=map[string]string{"return_url":returnURL,"cancel_url":cancelURL,"user_action":"PAY_NOW"}}
 raw,_:=json.Marshal(payload);req,_:=http.NewRequestWithContext(ctx,http.MethodPost,s.paypalBase()+"/v2/checkout/orders",bytes.NewReader(raw));req.Header.Set("Authorization","Bearer "+token);req.Header.Set("Content-Type","application/json");req.Header.Set("PayPal-Request-Id",orderNo)
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return "","",err};defer resp.Body.Close();body,_:=io.ReadAll(resp.Body);if resp.StatusCode/100!=2{return "","",fmt.Errorf("paypal: %s",body)}
 var x map[string]any;if err=json.Unmarshal(body,&x);err!=nil{return "","",err};id,_:=x["id"].(string);approve:="";if ls,ok:=x["links"].([]any);ok{for _,v:=range ls{if m,ok:=v.(map[string]any);ok{if m["rel"]=="approve"{approve,_=m["href"].(string)}}}};return id,approve,nil
}

func (s *PaymentServer) CapturePayment(ctx context.Context,r *pb.CapturePaymentRequest)(*pb.CapturePaymentResponse,error){
 var orderID int64;var provider,providerID,status string
 if err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT id,provider,provider_order_id,status FROM orders WHERE id=?",r.GetOrderId()).Scan(&orderID,&provider,&providerID,&status);err!=nil{return nil,err}
 if status=="PAID"{return &pb.CapturePaymentResponse{OrderId:orderID,Status:status},nil};if providerID==""{providerID=r.GetProviderOrderId()};if provider!="PAYPAL"{return nil,errors.New("Stripe is completed by webhook; PayPal requires capture")}
 token,err:=s.paypalToken(ctx);if err!=nil{return nil,err};req,_:=http.NewRequestWithContext(ctx,http.MethodPost,s.paypalBase()+"/v2/checkout/orders/"+providerID+"/capture",strings.NewReader("{}"));req.Header.Set("Authorization","Bearer "+token);req.Header.Set("Content-Type","application/json");req.Header.Set("PayPal-Request-Id",fmt.Sprintf("%d",orderID))
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();body,_:=io.ReadAll(resp.Body);if resp.StatusCode/100!=2{return nil,fmt.Errorf("paypal capture: %s",body)}
 var x map[string]any;if err=json.Unmarshal(body,&x);err!=nil{return nil,err};st,_:=x["status"].(string);if st=="COMPLETED"{if err=s.markPaid(ctx,orderID);err!=nil{return nil,err}};return &pb.CapturePaymentResponse{OrderId:orderID,Status:st},nil
}
func (s *PaymentServer) markPaid(ctx context.Context,orderID int64)error{
 tx,err:=s.svcCtx.DB.BeginTx(ctx,nil);if err!=nil{return err};defer tx.Rollback();var userID,dramaID int64
 if err=tx.QueryRowContext(ctx,"SELECT user_id,drama_id FROM orders WHERE id=?",orderID).Scan(&userID,&dramaID);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"UPDATE orders SET status='PAID',paid_at=NOW() WHERE id=? AND status<>'PAID'",orderID);err!=nil{return err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO user_entitlements(user_id,drama_id,order_id) VALUES(?,?,?) ON DUPLICATE KEY UPDATE order_id=VALUES(order_id)",userID,dramaID,orderID);err!=nil{return err};return tx.Commit()
}
func (s *PaymentServer) GetOrder(ctx context.Context,r *pb.GetOrderRequest)(*pb.GetOrderResponse,error){
 var x pb.GetOrderResponse;var cents int64
 err:=s.svcCtx.DB.QueryRowContext(ctx,"SELECT id,user_id,drama_id,order_no,provider,provider_order_id,amount_cents,currency,status FROM orders WHERE id=?",r.GetOrderId()).Scan(&x.OrderId,&x.UserId,&x.DramaId,&x.OrderNo,&x.Provider,&x.ProviderOrderId,&cents,&x.Currency,&x.Status);if err!=nil{return nil,err};x.Amount=fmt.Sprintf("%.2f",float64(cents)/100);return &x,nil
}
func providerName(p pb.Provider)string{if p==pb.Provider_PAYPAL{return "PAYPAL"};return "STRIPE"}
var _ = sql.ErrNoRows
