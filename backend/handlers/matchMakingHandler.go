package handlers

import (
	"context"
	"fmt"
	"manual-chess/models"
	"manual-chess/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var cancelChan = make(chan string)
var matchChan = make(chan string)
var resultChan = make(chan int)
var searchLock sync.Mutex

// @POST Add user to redis match making list
func FindMatch(c *gin.Context) {

	// TODO: find elo ranges based on user's MMR and other factors
	eloMin := 0
	eloMax := 3000

	// Add user to redis store of active users looking for match
	var user models.User

	// Bind the JSON data from the request body to the User struct
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.RedisClient.ZAdd(ctx, "mmList", redis.Z{Score: float64(user.MMR), Member: user.ID})

	// Try finding a match every x seconds
	go findMatchLoop(user.ID, eloMin, eloMax)

	// Use result channel to send response as required
	result := <-resultChan
	// 200 - Cancelled
	// 201 - Match found
	// Other codes are errors
	c.IndentedJSON(result, user)
}

// Helper function to find match
func findMatchLoop(id string, eloMin int, eloMax int) {
	for {
		select {
		case cancelId := <-cancelChan:
			if cancelId == id {
				fmt.Printf("Match making cancelled by %s\n", cancelId)
				resultChan <- http.StatusOK
				return
			}
		case matchedId := <-matchChan:
			if matchedId == id {
				fmt.Printf("Match making cancelled because of match found by %s\n", matchedId)
				resultChan <- http.StatusCreated
				return
			}
		default:
			searchLock.Lock()

			fmt.Printf("%s is finding match...\n", id)
			users, err := utils.RedisClient.ZRangeByScore(ctx, "mmList", &redis.ZRangeBy{Min: strconv.Itoa(eloMin), Max: strconv.Itoa(eloMax)}).Result()
			if err != nil {
				resultChan <- http.StatusInternalServerError
				return
			}

			var opp string

			// Attempt to find a match
			for i := 0; i < len(users); i++ {
				if users[i] != id {
					opp = users[i]
					break
				}
			}

			if opp != "" {
				defer searchLock.Unlock()
				// Match was found
				// Create lobby with redis hashes
				utils.RedisClient.ZRem(ctx, "mmList", id)
				utils.RedisClient.ZRem(ctx, "mmList", opp)
				matchChan <- opp
				resultChan <- http.StatusCreated

				fmt.Printf("Match found between %s and %s\n", opp, id)
				return
			}

			searchLock.Unlock()
			time.Sleep(4 * time.Second)
		}
	}
}

// @DELETE Remove user from redis match making list
func CancelMatch(c *gin.Context) {
	// Use go channels to signal
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := utils.RedisClient.ZRem(ctx, "mmList", user.ID).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if s == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with id " + user.ID + " is not finding a match."})
		return
	}

	cancelChan <- user.ID
}
