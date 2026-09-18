package model

import "time"

type BehaviorEvent struct {
	UserID         uint64
	DramaID        uint64
	EpisodeID      uint64
	EventType      string
	WatchSeconds   uint32
	DurationSeconds uint32
	Country        string
	Language       string
	Device         string
	EventAt        time.Time
}
