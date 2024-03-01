package infrastructure

import (
	"github.com/huandu/skiplist"
	"github.com/redis/go-redis/v9"
)

type InMemMatchmakingQueue struct {
	queue   skiplist.SkipList
	userMap map[string]*skiplist.Element
}

func NewInMemMatchmakingQueue() *InMemMatchmakingQueue {
	return &InMemMatchmakingQueue{
		queue:   *skiplist.New(skiplist.Int),
		userMap: make(map[string]*skiplist.Element),
	}
}

func (r *InMemMatchmakingQueue) GetPlayersByEloRange(eloMin int, eloMax int) ([]redis.Z, error) {
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

func (r *InMemMatchmakingQueue) AddPlayer(id string, score int) error {
	r.userMap[id] = r.queue.Set(score, redis.Z{Score: float64(score), Member: id})
	return nil
}

func (r *InMemMatchmakingQueue) RemovePlayer(id string) error {
	r.queue.RemoveElement(r.userMap[id])
	r.userMap[id] = nil
	return nil
}
