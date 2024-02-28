package services

import (
	"fmt"
	"manual-chess/constants"
	matchmakingRepository "manual-chess/repository/matchmaking"
	playerRepository "manual-chess/repository/player"
	"math"
	"sync"
	"time"
)

// Need to implement some form of websocket comunication

type queueMember struct {
	id        string
	cancelled bool
	weight    int
}

type MatchMakingService struct {
	playerRepo    playerRepository.PlayerRepository
	mmRepo        matchmakingRepository.MatchMakingRepository
	socketService *SocketService
	queue         []queueMember
	mu            sync.Mutex
}

func NewMatchMakingService(
	s *SocketService,
	p playerRepository.PlayerRepository,
	m matchmakingRepository.MatchMakingRepository,
) *MatchMakingService {
	return &MatchMakingService{
		playerRepo:    p,
		mmRepo:        m,
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
			userData, _ := m.playerRepo.GetPlayerById(m.queue[i].id)
			if userData.State != constants.InQueue {
				// mark for deletion
				m.queue[i].cancelled = true
				m.mu.Unlock()
				continue
			}

			// Find best match for current player
			if m.queue[i].weight < 5 {
				m.queue[i].weight++
			}
			eloMin := userData.MMR - 20*m.queue[i].weight
			eloMax := userData.MMR + 20*m.queue[i].weight
			potentialOpps, _ := m.mmRepo.GetPlayersByEloRange(eloMin, eloMax)
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
				oppData, _ := m.playerRepo.GetPlayerById(oppId)
				oppData.State = constants.InGame
				userData.State = constants.InGame
				m.playerRepo.SetPlayerById(oppId, oppData)
				m.playerRepo.SetPlayerById(m.queue[i].id, userData)

				// Delete user and opponent from sorted set
				m.mmRepo.RemovePlayer(m.queue[i].id)

				// Create new match
				// r := rand.New(rand.NewSource(time.Now().Unix()))
				// randomNumber := r.Float64()
				// actionFirst := ""
				// if randomNumber < 0.5 {
				// 	actionFirst = oppId
				// } else {
				// 	actionFirst = m.queue[i].id
				// }

				// matchId := uuid.New()
				// match := models.Match{
				// 	ID:      matchId.String(),
				// 	State:   constants.Select,
				// 	Player1: m.queue[i].id,
				// 	Player2: oppId,
				// 	Action:  actionFirst,
				// }

				// utils.MarshalAndSet[models.Match](m.redisClient, matchId.String(), &match)
				// m.socketService.GetConnection(m.queue[i].id).WriteJSON(map[string]interface{}{"matchId": matchId})
				// m.socketService.GetConnection(oppId).WriteJSON(map[string]interface{}{"matchId": matchId})
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
	defer m.mu.Unlock()
	// Attempt to find match (TODO)

	// Otherwise place user in queue
	user, _ := m.playerRepo.GetPlayerById(id)
	user.State = constants.InQueue
	m.playerRepo.SetPlayerById(id, user)

	// Add player to redis sorted set and matchmaking queue
	m.mmRepo.AddPlayer(id, user.MMR)
	m.queue = append(m.queue, queueMember{user.ID, false, 0})

	fmt.Println("Added user " + id + " to the queue")
}

func (m *MatchMakingService) RemoveFromMatchMakingQueue(id string) {
	// Remove player from redis sorted set and mark the player as In Lobby
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mmRepo.RemovePlayer(id)

	user, _ := m.playerRepo.GetPlayerById(id)
	user.State = constants.InLobby
	m.playerRepo.SetPlayerById(id, user)
}
