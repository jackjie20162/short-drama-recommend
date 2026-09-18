package logic
import("short-drama-recommend/internal/model";"short-drama-recommend/rpc/drama-rpc/pb")
func toPBDrama(d *model.Drama)*pb.Drama{return &pb.Drama{Id:int64(d.ID),Title:d.Title,Description:d.Description,Cover:d.Cover,Country:d.Country,Language:d.Language,TotalEpisodes:int32(d.TotalEpisodes),IsPaid:d.IsPaid}}
func toPBEpisode(e *model.Episode)*pb.Episode{return &pb.Episode{Id:int64(e.ID),DramaId:int64(e.DramaID),EpisodeNo:int32(e.EpisodeNo),Title:e.Title,DurationSeconds:int32(e.DurationSeconds),VideoUrl:e.VideoURL,PosterUrl:e.PosterURL,IsPaid:e.IsPaid,Status:int32(e.Status)}}
