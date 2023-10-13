package services

import (
	"context"
	"fmt"
	"manual-chess/models"
	"manual-chess/utils"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Need to implement some form of websocket communication

type queueMember struct {
	ID        string
	Cancelled bool
}

var queue []queueMember

func InitMatchMaker() {
	for {

		// Make matches
		for i := 0; i < len(queue); i++ {
			fmt.Println("Running match making coordinator.")

			// Get user data from redis map
			userData, _ := utils.GetAndUnmarshal[models.User](queue[i].ID)
			if userData.State != models.InQueue {
				queue[i].Cancelled = true
				continue
			}

			// Find suitable matches for user in redis sorted set
			eloMin := strconv.Itoa(userData.MMR - 50)
			eloMax := strconv.Itoa(userData.MMR + 50)
			potentialOpps, _ := utils.RedisClient.ZRangeByScoreWithScores(context.Background(), "mmList", &redis.ZRangeBy{Min: eloMin, Max: eloMax}).Result()

			// Find best match
			bestMatchDiff := math.Inf(1)
			var oppId string
			for j := 0; j < len(potentialOpps); j++ {
				if potentialOpps[j].Member != queue[i].ID && math.Abs(potentialOpps[j].Score-float64(userData.MMR)) < bestMatchDiff {
					oppId = potentialOpps[j].Member.(string)
				}
			}

			if oppId != "" {
				fmt.Printf("Found match between %s and %s\n", oppId, queue[i].ID)
				// Mark in queue that opp and userId are in match
				queue[i].Cancelled = true

				// Set user and opponent status to in match
				oppData, _ := utils.GetAndUnmarshal[models.User](oppId)
				oppData.State = models.InGame
				userData.State = models.InGame
				utils.MarshalAndSet[models.User](oppId, oppData)
				utils.MarshalAndSet[models.User](queue[i].ID, userData)

				// Delete user and opponent from sorted set
				utils.RedisClient.ZRem(context.Background(), "mmList", oppId, queue[i].ID)
			}
		}

		// Clean up queue
		var tempQueue []queueMember
		for _, member := range queue {
			if member.Cancelled == false {
				tempQueue = append(tempQueue, member)
			}
		}
		queue = tempQueue

		// Potential optimization: increase time between match maker runs when
		//    no matches found or no users are in the queue
		time.Sleep(5 * time.Second)
	}
}

func AddToMatchMakingQueue(id string) {
	// Mark player as finding match in redis
	user, _ := utils.GetAndUnmarshal[models.User](id)
	user.State = models.InQueue
	utils.MarshalAndSet[models.User](id, user)

	// Add player to redis sorted set and matchmaking queue
	utils.RedisClient.ZAdd(context.Background(), "mmList", redis.Z{Score: float64(user.MMR), Member: id})
	queue = append(queue, queueMember{user.ID, false})
}

func RemoveFromMatchMakingQueue(id string) {
	// Remove player from redis sorted set and mark the player as In Lobby
	utils.RedisClient.ZRem(context.Background(), "mmList", id)
	user, _ := utils.GetAndUnmarshal[models.User](id)

	user.State = models.InLobby
	utils.MarshalAndSet[models.User](id, user)
}
