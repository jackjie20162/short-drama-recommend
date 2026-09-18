package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(addr, password string) *Redis {
	return &Redis{Client: redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})}
}

func (r *Redis) Close() error { return r.Client.Close() }

func (r *Redis) IncrDramaEvent(ctx context.Context, dramaID int64, eventType string) error {
	return r.Client.Incr(ctx, fmt.Sprintf("drama:event:%d:%s", dramaID, eventType)).Err()
}

func (r *Redis) SetUserRecentDrama(ctx context.Context, userID, dramaID int64) error {
	key := fmt.Sprintf("user:recent:%d", userID)
	if err := r.Client.LPush(ctx, key, dramaID).Err(); err != nil { return err }
	if err := r.Client.LTrim(ctx, key, 0, 49).Err(); err != nil { return err }
	return r.Client.Expire(ctx, key, 7*24*60*60*1e9).Err()
}
