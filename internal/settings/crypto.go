package settings

import (
 "crypto/aes"
 "crypto/cipher"
 "crypto/rand"
 "encoding/base64"
 "errors"
 "io"
 "os"
)

func key() ([]byte,error) {
 raw:=os.Getenv("SETTINGS_ENCRYPTION_KEY")
 if raw=="" { return nil, errors.New("SETTINGS_ENCRYPTION_KEY is not configured") }
 b:=[]byte(raw)
 if len(b)!=32 { return nil, errors.New("SETTINGS_ENCRYPTION_KEY must be exactly 32 bytes") }
 return b,nil
}
func Encrypt(plain string)(string,error){
 k,e:=key();if e!=nil{return "",e}; block,e:=aes.NewCipher(k);if e!=nil{return "",e}
 g,e:=cipher.NewGCM(block);if e!=nil{return "",e}; nonce:=make([]byte,g.NonceSize());if _,e=io.ReadFull(rand.Reader,nonce);e!=nil{return "",e}
 return base64.StdEncoding.EncodeToString(g.Seal(nonce,nonce,[]byte(plain),nil)),nil
}
func Decrypt(encoded string)(string,error){
 k,e:=key();if e!=nil{return "",e}; raw,e:=base64.StdEncoding.DecodeString(encoded);if e!=nil{return "",e}
 block,e:=aes.NewCipher(k);if e!=nil{return "",e};g,e:=cipher.NewGCM(block);if e!=nil{return "",e}
 n:=g.NonceSize();if len(raw)<n{return "",errors.New("invalid encrypted setting")};plain,e:=g.Open(nil,raw[:n],raw[n:],nil);if e!=nil{return "",e};return string(plain),nil
}
