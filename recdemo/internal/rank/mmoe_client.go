package rank

import (
 "bytes"
 "context"
 "encoding/json"
 "fmt"
 "math"
 "net/http"
 "strings"
 "time"
)

type MMoEClient struct { address, model string; http *http.Client }
type Prediction struct { CTR float64; WatchProbs []float64; Completion float64; Pay float64 }

func NewMMoEClient(address, model string) *MMoEClient {
 return &MMoEClient{address:strings.TrimRight(address,"/"),model:model,http:&http.Client{Timeout:120*time.Millisecond}}
}

func (c *MMoEClient) PredictV1(ctx context.Context, sparse [][]int64, dense [][]float32) ([]Prediction,error) {
 if c.address==""||c.model=="" { return nil,fmt.Errorf("mmoe endpoint/model is empty") }
 if len(sparse)!=5||len(dense)==0 { return nil,fmt.Errorf("invalid mmoe batch shape") }
 n:=len(dense)
 for _,v:=range sparse { if len(v)!=n{return nil,fmt.Errorf("sparse batch length mismatch")} }
 inputs:=[]any{
  map[string]any{"name":"user_id","shape":[]int{n,1},"datatype":"INT64","data":sparse[0]},
  map[string]any{"name":"drama_id","shape":[]int{n,1},"datatype":"INT64","data":sparse[1]},
  map[string]any{"name":"region_id","shape":[]int{n,1},"datatype":"INT64","data":sparse[2]},
  map[string]any{"name":"language_id","shape":[]int{n,1},"datatype":"INT64","data":sparse[3]},
  map[string]any{"name":"genre_id","shape":[]int{n,1},"datatype":"INT64","data":sparse[4]},
  map[string]any{"name":"dense_features","shape":[]int{n,8},"datatype":"FP32","data":flattenFloat32(dense)},
 }
 body:=map[string]any{"inputs":inputs,"outputs":[]any{
  map[string]any{"name":"p_ctr"},map[string]any{"name":"watch_logits"},
  map[string]any{"name":"p_completion"},map[string]any{"name":"p_pay"},
 }}
 raw,err:=json.Marshal(body);if err!=nil{return nil,err}
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,c.address+"/v2/models/"+c.model+"/infer",bytes.NewReader(raw));if err!=nil{return nil,err}
 req.Header.Set("Content-Type","application/json")
 resp,err:=c.http.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 if resp.StatusCode<200||resp.StatusCode>=300{return nil,fmt.Errorf("triton status=%d",resp.StatusCode)}
 var decoded struct{Outputs []struct{Name string `json:"name"`;Data []float64 `json:"data"`}`json:"outputs"`}
 if err:=json.NewDecoder(resp.Body).Decode(&decoded);err!=nil{return nil,err}
 pred:=make([]Prediction,n)
 for _,out:=range decoded.Outputs {
  switch out.Name {
  case "p_ctr": for i:=0;i<n&&i<len(out.Data);i++{pred[i].CTR=sigmoid(out.Data[i])}
  case "watch_logits": for i:=0;i<n;i++{s:=i*7;e:=s+7;if e<=len(out.Data){pred[i].WatchProbs=softmax(out.Data[s:e])}}
  case "p_completion": for i:=0;i<n&&i<len(out.Data);i++{pred[i].Completion=sigmoid(out.Data[i])}
  case "p_pay": for i:=0;i<n&&i<len(out.Data);i++{pred[i].Pay=sigmoid(out.Data[i])}
  }
 }
 for i:=range pred{if len(pred[i].WatchProbs)!=7{return nil,fmt.Errorf("triton missing watch_logits for row %d",i)}}
 return pred,nil
}

func flattenFloat32(v [][]float32)[]float32{out:=make([]float32,0,len(v)*8);for _,row:=range v{out=append(out,row...)};return out}
func sigmoid(x float64)float64{return 1/(1+math.Exp(-x))}
func softmax(v []float64)[]float64{if len(v)==0{return nil};m:=v[0];for _,x:=range v[1:]{if x>m{m=x}};sum:=0.0;out:=make([]float64,len(v));for i,x:=range v{out[i]=math.Exp(x-m);sum+=out[i]};for i:=range out{out[i]/=sum};return out}
