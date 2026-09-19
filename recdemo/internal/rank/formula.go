package rank

import "math"

type Config struct {
	PopularityWeight float64
	CompletionWeight float64
	PayWeight float64
	FreshnessWeight float64
	SeenPenalty float64
}

type Features struct {
	Popularity float64
	Completion float64
	PayRate float64
	PublishedAt int64
	Seen bool
	RegionMatch bool
	LangMatch bool
}

type Ranker struct { cfg Config }

func New(c Config) *Ranker { return &Ranker{cfg:c} }

func (r *Ranker) Score(f Features, now int64) float64 {
	age := float64(now-f.PublishedAt)
	if age < 0 { age = 0 }
	fresh := math.Exp(-age / (30*86400))
	score := r.cfg.PopularityWeight*f.Popularity +
		r.cfg.CompletionWeight*f.Completion +
		r.cfg.PayWeight*f.PayRate +
		r.cfg.FreshnessWeight*fresh
	if f.RegionMatch { score += 0.08 }
	if f.LangMatch { score += 0.12 }
	if f.Seen { score -= r.cfg.SeenPenalty }
	return score
}
