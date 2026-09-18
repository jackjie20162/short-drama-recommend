package model

type Drama struct {
	ID             uint64
	Title          string
	Description    string
	Cover          string
	Country        string
	Language       string
	TotalEpisodes  uint32
	IsPaid         bool
	Status         int8
}
