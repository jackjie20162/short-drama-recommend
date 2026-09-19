package model

type Episode struct {
	ID              uint64
	DramaID         uint64
	EpisodeNo       uint32
	Title           string
	Description     string
	DurationSeconds uint32
	VideoURL        string
	PosterURL       string
	IsPaid          bool
	Status          int8
}
