package repository

import (
	"github.com/huandu/skiplist"
	"github.com/redis/go-redis/v9"
)

type InMemMatchmakingRepository struct {
	queue skiplist.SkipList
}

func NewInMemMatchmakingRepository() *InMemMatchmakingRepository {
	return &InMemMatchmakingRepository{
		queue: *skiplist.New(skiplist.Int),
	}
}

func (r *InMemMatchmakingRepository) GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error) {
	start := r.queue.FindNext(nil, eloMin)
	res := []redis.Z{}
	if start != nil {
		res = append(res, start.Value.(redis.Z))
		for start.Next() != nil {
			next := start.Next().Value.(redis.Z)
			if int(next.Score) > eloMax {
				break
			}
			res = append(res, start.Next().Value.(redis.Z))
			start = start.Next()
		}
	}
	return res, nil
}

func (r *InMemMatchmakingRepository) AddPlayer(id string, score int) error {
	r.queue.Set(score, redis.Z{Score: float64(score), Member: id})
	return nil
}

func (r *InMemMatchmakingRepository) RemovePlayer(id string) error {
	r.queue.Remove(id)
	return nil
}
