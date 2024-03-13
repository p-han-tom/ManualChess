package matchmakingservice

import (
	"fmt"
	"manual-chess/constants"
	mathcmaking "manual-chess/infrastructure/matchmaking"
	matchRepository "manual-chess/repository/match"
	userRepository "manual-chess/repository/user"
	gameServices "manual-chess/services/game"
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
	userRepo    userRepository.IUserRepository
	matchRepo   matchRepository.IMatchRepository
	mmQueue     mathcmaking.IMatchMakingQueue
	gameService *gameServices.GameService
	queue       []queueMember
	mu          sync.Mutex
}

func NewMatchMakingService(
	gameService *gameServices.GameService,
	userRepo userRepository.IUserRepository,
	matchRepo matchRepository.IMatchRepository,
	mmQueue mathcmaking.IMatchMakingQueue) *MatchMakingService {
	return &MatchMakingService{
		userRepo:    userRepo,
		matchRepo:   matchRepo,
		mmQueue:     mmQueue,
		gameService: gameService,
	}
}

// Two pass polling approach
//   - First pass matches players and marks players who have found a match
//   - Second pass deletes matched players from the queue
func (m *MatchMakingService) RunMatchMaker() {

	for {
		// fmt.Println("Running match making coordinator.")
		for i := 0; i < len(m.queue); i++ {
			currId := m.queue[i].id

			m.mu.Lock()
			// Get user data from redis map
			userData, _ := m.userRepo.GetUserById(currId)
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
			potentialOpps, _ := m.mmQueue.GetPlayersByEloRange(eloMin, eloMax)
			bestMatchDiff := math.Inf(1)
			var oppId string
			for j := 0; j < len(potentialOpps); j++ {
				mmrDiff := math.Abs(potentialOpps[j].Score - float64(userData.MMR))
				if potentialOpps[j].Member != currId && mmrDiff < bestMatchDiff {
					bestMatchDiff, oppId = mmrDiff, potentialOpps[j].Member.(string)
				}
			}

			if bestMatchDiff < math.Inf(1) {
				fmt.Printf("Found match between %s and %s\n", oppId, currId)

				// Set user and opponent status to in match
				oppData, _ := m.userRepo.GetUserById(oppId)
				oppData.State = constants.InGame
				userData.State = constants.InGame
				m.userRepo.SetUserById(oppId, oppData)
				m.userRepo.SetUserById(currId, userData)

				// Delete user and opponent from sorted set
				m.mmQueue.RemovePlayer(currId)
				m.mmQueue.RemovePlayer(oppId)

				// Create new match
				m.gameService.SetupMatch(currId, oppId)
			}
			m.mu.Unlock()
		}

		// sweep cancelled requests
		tempQueue := make([]queueMember, 0, len(m.queue))
		for _, member := range m.queue {
			if !member.cancelled {
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
	user, _ := m.userRepo.GetUserById(id)
	user.State = constants.InQueue
	m.userRepo.SetUserById(id, user)

	// Add player to redis sorted set and matchmaking queue
	m.mmQueue.AddPlayer(id, user.MMR)
	m.queue = append(m.queue, queueMember{user.ID, false, 0})

	fmt.Println("Added user " + id + " to the queue")
}

func (m *MatchMakingService) RemoveFromMatchMakingQueue(id string) {
	// Remove player from redis sorted set and mark the player as In Lobby
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mmQueue.RemovePlayer(id)

	user, _ := m.userRepo.GetUserById(id)
	user.State = constants.InLobby
	m.userRepo.SetUserById(id, user)
}
