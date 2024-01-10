package services

import (
	"context"
	"fmt"
	"manual-chess/models"
	"manual-chess/utils"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Need to implement some form of websocket comunication

type queueMember struct {
	id        string
	cancelled bool
	weight    int
}

type MatchMakingService struct {
	RedisClient   *redis.Client
	socketService *SocketService
	queue         []queueMember
	mu            sync.Mutex
}

func NewMatchMakingService(s *SocketService, r *redis.Client) *MatchMakingService {
	return &MatchMakingService{
		RedisClient:   r,
		socketService: s,
	}
}

// Two pass polling approach
//   - First pass matches players and marks players who have found a match
//   - Second pass deletes matched players from the queue
func (m *MatchMakingService) RunMatchMaker() {

	for {
		fmt.Println("Running match making coordinator.")
		for i := 0; i < len(m.queue); i++ {
			m.mu.Lock()
			// Get user data from redis map
			userData, _ := utils.GetAndUnmarshal[models.User](m.RedisClient, m.queue[i].id)
			if userData.State != models.InQueue {
				// mark for deletion
				m.queue[i].cancelled = true
				m.mu.Unlock()
				continue
			}

			// Find best match for current player
			if m.queue[i].weight < 5 {
				m.queue[i].weight++
			}
			eloMin := strconv.Itoa(userData.MMR - 20*m.queue[i].weight)
			eloMax := strconv.Itoa(userData.MMR + 20*m.queue[i].weight)
			potentialOpps, _ := m.RedisClient.ZRangeByScoreWithScores(context.Background(), "mmList", &redis.ZRangeBy{Min: eloMin, Max: eloMax}).Result()
			bestMatchDiff := math.Inf(1)
			var oppId string
			for j := 0; j < len(potentialOpps); j++ {
				mmrDiff := math.Abs(potentialOpps[j].Score - float64(userData.MMR))
				if potentialOpps[j].Member != m.queue[i].id && mmrDiff < bestMatchDiff {
					bestMatchDiff, oppId = mmrDiff, potentialOpps[j].Member.(string)
				}
			}

			if bestMatchDiff < math.Inf(1) {
				fmt.Printf("Found match between %s and %s\n", oppId, m.queue[i].id)

				// Set user and opponent status to in match
				oppData, _ := utils.GetAndUnmarshal[models.User](m.RedisClient, oppId)
				oppData.State = models.InGame
				userData.State = models.InGame
				utils.MarshalAndSet[models.User](m.RedisClient, oppId, oppData)
				utils.MarshalAndSet[models.User](m.RedisClient, m.queue[i].id, userData)

				// Delete user and opponent from sorted set
				m.RedisClient.ZRem(context.Background(), "mmList", oppId, m.queue[i].id)

				// Create new match
				matchId := uuid.New()
				m.socketService.GetConnection(m.queue[i].id).WriteJSON(map[string]interface{}{"matchId": matchId})
				m.socketService.GetConnection(oppId).WriteJSON(map[string]interface{}{"matchId": matchId})
			}
			m.mu.Unlock()
		}
		// sweep cancelled requests
		tempQueue := make([]queueMember, 0, len(m.queue))
		for _, member := range m.queue {
			if member.cancelled == false {
				tempQueue = append(tempQueue, member)
			}
		}
		m.queue = tempQueue

		// Potential optimization: increase time between match maker runs when
		//    no matches found or no users are in the queue
		time.Sleep(5 * time.Second)
	}
}

func (m *MatchMakingService) AddToMatchMakingQueue(id string) {
	// Mark player as finding match in redis
	m.mu.Lock()

	// Attempt to find match (TODO)

	// Otherwise place user in queue
	user, _ := utils.GetAndUnmarshal[models.User](m.RedisClient, id)
	user.State = models.InQueue
	utils.MarshalAndSet[models.User](m.RedisClient, id, user)

	// Add player to redis sorted set and matchmaking queue
	m.RedisClient.ZAdd(context.Background(), "mmList", redis.Z{Score: float64(user.MMR), Member: id})
	m.queue = append(m.queue, queueMember{user.ID, false, 0})

	m.mu.Unlock()
	fmt.Println("Added user " + id + " to the queue")
}

func (m *MatchMakingService) RemoveFromMatchMakingQueue(id string) {
	// Remove player from redis sorted set and mark the player as In Lobby
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RedisClient.ZRem(context.Background(), "mmList", id)
	user, _ := utils.GetAndUnmarshal[models.User](m.RedisClient, id)

	user.State = models.InLobby
	utils.MarshalAndSet[models.User](m.RedisClient, id, user)

}
