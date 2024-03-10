package runner

import (
	game "GoCats/Game"
	user "GoCats/User"
	"log"
	"sort"
	"sync"
)

/*

	|| Get the request from the HTTP server and communicate with the modules to satisfy the request needs. ||

	Functionalities of this Package:
		1. Manage a pool of users objects
		2. Communicate with different modules
*/

// Features:
// 1. Gameplay of the user must be stored, i.e., if the user leaves the game and joins back his/her game should should start from where he/she left.
// 2. There should be a leader board of points that updates in real-time(web-socket connection)

type UserPool map[string]*user.UserStruct

type ByField1 []user.UserStruct

func (a ByField1) Len() int           { return len(a) }
func (a ByField1) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByField1) Less(i, j int) bool { return a[i].Points > a[j].Points }

type RunnerStruct struct {
	Game     game.GameStruct
	User     user.UserModule
	userPool UserPool
	sync.RWMutex
	wsChan chan<- []map[string]interface{}
}

func NewRunner(user user.UserModule, ws chan<- []map[string]interface{}) *RunnerStruct {
	return &RunnerStruct{User: user, userPool: make(UserPool), wsChan: ws}
}

func (R *RunnerStruct) LoginUser(userName, password string) *user.UserStruct {
	var user *user.UserStruct
	// Find in the Pool first, if not then in Redis
	user = R.userPool[userName]

	if user == nil {
		// Find User in RedisDB

		bufUser := R.User.GetUserWithUsername(userName)
		user = &bufUser
		R.AddUser(user)
		log.Println(R.userPool)
	}

	return user
}

func (R *RunnerStruct) AddUser(U *user.UserStruct) {
	R.Lock()
	defer R.Unlock()
	R.userPool[U.Username] = U
}

func (R *RunnerStruct) StartGame(username string) (game.GameStruct, error) {
	user := R.userPool[username]
	if user == nil {
		userFromRedis := R.User.GetUserWithUsername(username)
		user = &userFromRedis
		R.AddUser(user)
	}
	var gg *game.GameStruct
	log.Println(user)
	if user.GameStatus != nil && user.GameStatus.Status == "playing" {
		return *user.GameStatus, nil
	}
	gg = game.StartNewGame()
	go func() {
		R.Lock()
		defer R.Unlock()

		user.GameStatus = gg
		R.User.UpdateUserInRedis(user)
	}()
	log.Println("\nkokokokokokokoko\n")
	return *gg, nil

}

func (R *RunnerStruct) UserDrawsCard(username string) game.GameStruct {
	// find user in
	user := R.userPool[username]
	if user == nil {
		return game.GameStruct{}
	}
	if user.GameStatus != nil {
		user.GameStatus.DrawCard()
		// log.Println(user.GameStatus)

		if user.GameStatus.Status == "Winner" {
			user.Points++
			log.Println("\n\n\n", user, "\n\n\n")
		}

		R.User.UpdateUserInRedis(user)
		if user.GameStatus.Status == "Winner" {
			go func() {
				leaderBoard := R.GetLeaderBoard()
				R.wsChan <- leaderBoard
			}()
		}

		return *user.GameStatus
	}
	return game.GameStruct{}
}

func (R *RunnerStruct) GetLeaderBoard() []map[string]interface{} {
	usersArr := R.User.GetAllUsers()
	sort.Sort(ByField1(usersArr))
	var resultArray []map[string]interface{}
	for _, item := range usersArr {
		resultArray = append(resultArray, map[string]interface{}{
			"username": item.Username,
			"points":   item.Points,
		})
	}
	return resultArray
}
