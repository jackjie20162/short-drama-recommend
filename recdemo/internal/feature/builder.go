package feature

import (
	"hash/fnv"
	"time"
)

var Names = []string{
	"user_watch_seconds", "user_completion", "user_pay_rate", "user_sessions",
	"drama_popularity", "drama_completion", "drama_pay_rate", "drama_age_days",
}

type User struct {
	ID int64
	WatchSeconds, Completion, PayRate, Sessions float64
}

type Drama struct {
	ID int64
	RegionID, LanguageID, GenreID int64
	Popularity, Completion, PayRate float64
	PublishedAt int64
}

type Context struct {
	RegionID, LanguageID int64
	SameRegion, SameLanguage bool
	GenreAffinity float64
	Now time.Time
}

func BuildV1(u User, d Drama, c Context) (sparse []int64, dense []float32) {
	now := c.Now
	if now.IsZero() { now = time.Now() }
	age := now.Unix() - d.PublishedAt
	if age < 0 { age = 0 }
	sparse = []int64{
		normalizeID(u.ID, 1000000),
		normalizeID(d.ID, 1000000),
		normalizeID(c.RegionID, 256),
		normalizeID(c.LanguageID, 256),
		normalizeID(d.GenreID, 512),
	}
	dense = []float32{
		float32(u.WatchSeconds), float32(u.Completion), float32(u.PayRate), float32(u.Sessions),
		float32(d.Popularity), float32(d.Completion), float32(d.PayRate), float32(age / 86400),
	}
	return
}

func normalizeID(v, mod int64) int64 {
	if v < 0 { v = -v }
	if mod <= 0 { return v }
	return v % mod
}

func HashID(s string, mod int64) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return normalizeID(int64(h.Sum64()&0x7fffffffffffffff), mod)
}

// Build is kept as a compatibility helper for the old 18-float model.
func Build(u User, d Drama, c Context) []float32 {
	now := c.Now
	if now.IsZero() { now = time.Now() }
	age := now.Unix() - d.PublishedAt
	if age < 0 { age = 0 }
	return []float32{
		float32(u.WatchSeconds), float32(u.Completion), float32(u.PayRate), float32(u.Sessions),
		float32(d.Popularity), float32(d.Completion), float32(d.PayRate), float32(age / 86400),
		boolf(c.SameRegion), boolf(c.SameLanguage), float32(c.GenreAffinity),
		0, 0, float32(now.Hour()), float32(int(now.Weekday())), 0, 0, 0,
	}
}

func boolf(v bool) float32 { if v { return 1 }; return 0 }
