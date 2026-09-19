package model

import "time"

type DramaContent struct {
	ID             uint64
	Title          string
	Subtitle       string
	Description    string
	Cover          string
	Country        string
	Language       string
	Genres         []string
	Tags           []string
	TotalEpisodes  uint32
	Popularity     float64
	CompletionRate float64
	PayRate        float64
	Status         int8
	PublishedAt    time.Time
}

type RecommendationFeature struct {
	Sparse   SparseFeature
	Dense    DenseFeature
	Semantic SemanticFeature
}

type SparseFeature struct {
	UserID     uint64
	DramaID    uint64
	CountryID  int64
	LanguageID int64
	GenreIDs   []int64
	TagIDs     []int64
}

type DenseFeature struct {
	UserWatchSeconds float64
	UserCompletion   float64
	UserPayRate      float64
	UserSessions     float64
	DramaPopularity  float64
	DramaCompletion  float64
	DramaPayRate     float64
	DramaAgeDays     float64
}

type SemanticFeature struct {
	TitleEmbedding       []float32
	DescriptionEmbedding []float32
	GenreEmbedding       []float32
	TagEmbedding         []float32
}
