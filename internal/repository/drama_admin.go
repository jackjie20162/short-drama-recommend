package repository

import (
 "context"
 "database/sql"
 "strings"
 "short-drama-recommend/internal/model"
)

type DramaAdminRepository interface {
 Create(context.Context,*model.Drama) error
 Update(context.Context,*model.Drama) error
 SetStatus(context.Context,uint64,int8) error
 List(context.Context,string,string,string,*int8,int,int)([]*model.Drama,int64,error)
}
type MySQLDramaAdminRepository struct{db *sql.DB}
func NewMySQLDramaAdminRepository(db *sql.DB) DramaAdminRepository{return &MySQLDramaAdminRepository{db}}

func(r *MySQLDramaAdminRepository)Create(ctx context.Context,d *model.Drama)error{
 res,err:=r.db.ExecContext(ctx,"INSERT INTO dramas (title,description,cover,country,language,total_episodes,is_paid,status) VALUES (?,?,?,?,?,?,?,0)",d.Title,d.Description,d.Cover,d.Country,d.Language,d.TotalEpisodes,d.IsPaid);if err!=nil{return err};id,err:=res.LastInsertId();if err!=nil{return err};d.ID=uint64(id);return nil
}
func(r *MySQLDramaAdminRepository)Update(ctx context.Context,d *model.Drama)error{
 _,err:=r.db.ExecContext(ctx,"UPDATE dramas SET title=?,description=?,cover=?,country=?,language=?,total_episodes=?,is_paid=? WHERE id=?",d.Title,d.Description,d.Cover,d.Country,d.Language,d.TotalEpisodes,d.IsPaid,d.ID);return err
}
func(r *MySQLDramaAdminRepository)SetStatus(ctx context.Context,id uint64,status int8)error{
 _,err:=r.db.ExecContext(ctx,"UPDATE dramas SET status=?,published_at=CASE WHEN ?=1 THEN COALESCE(published_at,UTC_TIMESTAMP()) ELSE published_at END WHERE id=?",status,status,id);return err
}
func(r *MySQLDramaAdminRepository)List(ctx context.Context,keyword,country,language string,status *int8,limit,offset int)([]*model.Drama,int64,error){
 cond:=[]string{"1=1"};args:=[]any{}
 if strings.TrimSpace(keyword)!=""{k:="%"+strings.TrimSpace(keyword)+"%";cond=append(cond,"(title LIKE ? OR description LIKE ?)");args=append(args,k,k)}
 if country!=""{cond=append(cond,"country=?");args=append(args,country)}
 if language!=""{cond=append(cond,"language=?");args=append(args,language)}
 if status!=nil{cond=append(cond,"status=?");args=append(args,*status)}
 where:=" WHERE "+strings.Join(cond," AND ");var total int64
 if err:=r.db.QueryRowContext(ctx,"SELECT COUNT(*) FROM dramas"+where,args...).Scan(&total);err!=nil{return nil,0,err}
 rows,err:=r.db.QueryContext(ctx,"SELECT id,title,description,cover,country,language,total_episodes,is_paid,status FROM dramas"+where+" ORDER BY id DESC LIMIT ? OFFSET ?",append(args,limit,offset)...);if err!=nil{return nil,0,err};defer rows.Close()
 var out []*model.Drama
 for rows.Next(){d:=&model.Drama{};if err:=rows.Scan(&d.ID,&d.Title,&d.Description,&d.Cover,&d.Country,&d.Language,&d.TotalEpisodes,&d.IsPaid,&d.Status);err!=nil{return nil,0,err};out=append(out,d)}
 return out,total,rows.Err()
}