package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Claims struct { UserID int64 `json:"user_id"`; Email string `json:"email"`; Exp int64 `json:"exp"` }

func secret() []byte { s:=os.Getenv("AUTH_JWT_SECRET"); if s=="" { s="change-me-in-production" }; return []byte(s) }
func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
func sign(input string) string { h:=hmac.New(sha256.New,secret()); h.Write([]byte(input)); return b64(h.Sum(nil)) }

func Issue(userID int64,email string,ttl time.Duration)(string,error){
 if userID<=0{return "",errors.New("invalid user id")}; now:=time.Now().Unix(); p,_:=json.Marshal(Claims{UserID:userID,Email:email,Exp:now+int64(ttl.Seconds())}); head:=b64([]byte(`{"alg":"HS256","typ":"JWT"}`)); body:=b64(p); return head+"."+body+"."+sign(head+"."+body),nil
}
func Parse(token string)(Claims,error){
 var c Claims; parts:=strings.Split(token,".");if len(parts)!=3||!hmac.Equal([]byte(parts[2]),[]byte(sign(parts[0]+"."+parts[1]))){return c,errors.New("invalid token")}
 b,err:=base64.RawURLEncoding.DecodeString(parts[1]);if err!=nil{return c,errors.New("invalid token")};if err=json.Unmarshal(b,&c);err!=nil{return c,err};if c.UserID<=0||c.Exp<=time.Now().Unix(){return c,errors.New("token expired")};return c,nil
}
func Bearer(header string)(int64,error){if !strings.HasPrefix(header,"Bearer "){return 0,errors.New("missing bearer token")};c,e:=Parse(strings.TrimSpace(strings.TrimPrefix(header,"Bearer ")));return c.UserID,e}
func UserIDString(id int64) string{return strconv.FormatInt(id,10)}
