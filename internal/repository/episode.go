package repository

import (
 "context"
 "database/sql"
 "errors"
 "short-drama-recommend/internal/model"
)

type EpisodeRepository interface {
 Create(context.Context,*model.Episode) error
 ListByDrama(context.Context,uint64,*int8)([]*model.Episode,error)
 SetStatus(context.Context,uint64,int8) error
}

type MySQLEpisodeRepository struct{ db *sql.DB }
func NewMySQLEpisodeRepository(db *sql.DB) EpisodeRepository { return &MySQLEpisodeRepository{db} }

func (r *MySQLEpisodeRepository) Create(ctx context.Context,e *model.Episode) error {
 res,err:=r.db.ExecContext(ctx,"INSERT INTO episodes (drama_id,episode_no,title,duration_seconds,video_url,poster_url,is_paid,status) VALUES (?,?,?,?,?,?,?,0) ON DUPLICATE KEY UPDATE title=VALUES(title),duration_seconds=VALUES(duration_seconds),video_url=VALUES(video_url),poster_url=VALUES(poster_url),is_paid=VALUES(is_paid)",e.DramaID,e.EpisodeNo,e.Title,e.DurationSeconds,e.VideoURL,e.PosterURL,e.IsPaid)
 if err!=nil{return err}
 id,err:=res.LastInsertId(); if err==nil&&id>0{e.ID=uint64(id)}
 if e.ID==0 { err=r.db.QueryRowContext(ctx,"SELECT id FROM episodes WHERE drama_id=? AND episode_no=?",e.DramaID,e.EpisodeNo).Scan(&e.ID) }
 return err
}
func (r *MySQLEpisodeRepository) ListByDrama(ctx context.Context,dramaID uint64,status *int8)([]*model.Episode,error){
 q:="SELECT id,drama_id,episode_no,title,duration_seconds,video_url,poster_url,is_paid,status FROM episodes WHERE drama_id=?"; args:=[]any{dramaID}
 if status!=nil{q+=" AND status=?";args=append(args,*status)};q+=" ORDER BY episode_no ASC"
 rows,err:=r.db.QueryContext(ctx,q,args...);if err!=nil{return nil,err};defer rows.Close()
 var out []*model.Episode
 for rows.Next(){e:=&model.Episode{};if err:=rows.Scan(&e.ID,&e.DramaID,&e.EpisodeNo,&e.Title,&e.DurationSeconds,&e.VideoURL,&e.PosterURL,&e.IsPaid,&e.Status);err!=nil{return nil,err};out=append(out,e)}
 return out,rows.Err()
}
func (r *MySQLEpisodeRepository) SetStatus(ctx context.Context,id uint64,status int8) error {
 if id==0{return errors.New("episode id is required")};_,err:=r.db.ExecContext(ctx,"UPDATE episodes SET status=? WHERE id=?",status,id);return err
}