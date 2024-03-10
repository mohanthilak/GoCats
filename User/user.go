package user

import (
	game "GoCats/Game"
	redispkg "GoCats/Redis"
	"encoding/json"
	"log"
)

type UserStruct struct {
	Username    string           `json:"username"`
	Password    string           `json:"password"`
	AccessToken string           `json:"accessToken"`
	Points      int              `json:"points"`
	GameStatus  *game.GameStruct `json:"gameStatus"`
}

type UserModule struct {
	redis *redispkg.RedisStruct
}

func CreateUserModule(redis *redispkg.RedisStruct) *UserModule {
	return &UserModule{redis: redis}
}

func (U *UserModule) UpdateUserInRedis(user *UserStruct) {
	// Update the DB
	U.redis.SetObj(user.Username, user)
}

// Looks for user in Redis. If not found returns a new user with the given Username
func (U *UserModule) GetUserWithUsername(username string) UserStruct {
	userVal := U.redis.GetObj(username)
	if userVal == nil {
		user := UserStruct{Username: username}
		U.redis.SetObj(username, user)
		return user
	}
	var user UserStruct
	json.Unmarshal([]byte(userVal), &user)
	if user.GameStatus != nil && user.GameStatus.Status == "playing" {
		game.HydrateActionsOnCards(user.GameStatus.Cards)
	}
	return user
}

func (U *UserModule) GetAllUsers() []UserStruct {
	var users []UserStruct
	JsonArray := U.redis.GetUsers()

	for _, K := range JsonArray {
		var User UserStruct
		if err := json.Unmarshal([]byte(K), &User); err != nil {
			log.Printf("Error decoding JSON for key %s: %v\n", K, err)
			continue
		}
		users = append(users, User)
	}

	return users
}
