package feature

import "time"
var Names=[]string{"user_watch_seconds","user_completion","user_pay_rate","user_sessions","drama_popularity","drama_completion","drama_pay_rate","drama_age_days","same_region","same_language","genre_affinity","is_new","is_continue","hour","weekday","country_score","language_score","price_score"}
type User struct{WatchSeconds,Completion,PayRate,Sessions float64}
type Drama struct{Popularity,Completion,PayRate float64; PublishedAt int64}
type Context struct{SameRegion,SameLanguage bool; GenreAffinity float64; IsNew,IsContinue bool; CountryScore,LanguageScore,PriceScore float64; Now time.Time}
func Build(u User,d Drama,c Context)[]float32{now:=c.Now;if now.IsZero(){now=time.Now()};age:=now.Unix()-d.PublishedAt;if age<0{age=0};return []float32{float32(u.WatchSeconds),float32(u.Completion),float32(u.PayRate),float32(u.Sessions),float32(d.Popularity),float32(d.Completion),float32(d.PayRate),float32(age/86400),boolf(c.SameRegion),boolf(c.SameLanguage),float32(c.GenreAffinity),boolf(c.IsNew),boolf(c.IsContinue),float32(now.Hour()),float32(int(now.Weekday())),float32(c.CountryScore),float32(c.LanguageScore),float32(c.PriceScore)}}
func boolf(v bool)float32{if v{return 1};return 0}
