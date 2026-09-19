package model

import "time"

type Drama struct {
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
	IsPaid         bool
	Status         int8
	PriceCents     uint64
	Currency       string
	Popularity     float64
	CompletionRate float64
	PayRate        float64
	PublishedAt    time.Time
}
