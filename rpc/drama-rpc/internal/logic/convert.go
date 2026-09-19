package logic

import (
	"short-drama-recommend/internal/model"
	"short-drama-recommend/rpc/drama-rpc/pb"
)

func toPBDrama(d *model.Drama) *pb.Drama {
	return &pb.Drama{
		Id: int64(d.ID), Title: d.Title, Subtitle: d.Subtitle, Description: d.Description,
		Cover: d.Cover, Country: d.Country, Language: d.Language, Genres: d.Genres, Tags: d.Tags,
		TotalEpisodes: int32(d.TotalEpisodes), IsPaid: d.IsPaid, PriceCents: int64(d.PriceCents),
		Currency: d.Currency, Popularity: d.Popularity, CompletionRate: d.CompletionRate,
		PayRate: d.PayRate, PublishedAt: d.PublishedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPBEpisode(e *model.Episode) *pb.Episode {
	return &pb.Episode{Id:int64(e.ID), DramaId:int64(e.DramaID), EpisodeNo:int32(e.EpisodeNo), Title:e.Title,
		DurationSeconds:int32(e.DurationSeconds), VideoUrl:e.VideoURL, PosterUrl:e.PosterURL, IsPaid:e.IsPaid, Status:int32(e.Status)}
}

func toPBPublicEpisode(e *model.Episode, unlocked bool) *pb.DramaEpisode {
	return &pb.DramaEpisode{Id:int64(e.ID), DramaId:int64(e.DramaID), EpisodeNo:int32(e.EpisodeNo), Title:e.Title,
		DurationSeconds:int32(e.DurationSeconds), VideoUrl:e.VideoURL, PosterUrl:e.PosterURL, IsPaid:e.IsPaid, Status:int32(e.Status), Unlocked:unlocked,
		Description:e.Description}
}
