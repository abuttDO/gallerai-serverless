package main

import (
	"errors"
	"fmt"
)

type Request struct {
	Username string `json:"username"`
	ImageID  int    `json:"image_id"`
	Reaction string `json:"reaction"`
	Type     string `json:"type"`
}

var (
	Reactions            = []string{"like", "love", "haha", "wow", "sad", "angry", "report"}
	InvalidReactionError = errors.New("invalid reaction")
	InvalidReactionType  = errors.New("invalid reaction type")
)

const (
	RequestTypeGet  = "GET"
	RequestTypePost = "POST"
)

func Main(in Request) (*Response, error) {
	switch in.Type {
	case RequestTypeGet:
		reactions, err := getReactionsByImageID(in.ImageID)
		if err != nil {
			return makeResponse(500, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		return makeResponse(200, reactions, nil), nil
	case RequestTypePost:
		valid := validateReaction(in.Reaction)
		if !valid {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, InvalidReactionError), InvalidReactionError), InvalidReactionError
		}
		repo.db = initDatabase()
		var reaction Reaction
		user, err := getUserByUsername(in.Username)
		if err != nil {
			return makeResponse(500, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		reaction.UserID = user.ID
		reaction.ImageID = in.ImageID
		reaction.Reaction = in.Reaction
		saveReaction(reaction)
		return makeResponse(200, "", nil), nil
	default:
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, InvalidReactionType), InvalidReactionType), InvalidReactionType
	}
}

func saveReaction(reaction Reaction) {
	repo.db = initDatabase()
	repo.db.Create(&reaction)
}

func getUserByUsername(username string) (User, error) {
	repo.db = initDatabase()
	var user User
	repo.db.Where("username = ?", username).First(&user)
	return user, nil
}

func validateReaction(reaction string) bool {
	for _, r := range Reactions {
		if r == reaction {
			return true
		}
	}
	return false
}

func getReactionsByImageID(imageID int) ([]Reaction, error) {
	repo.db = initDatabase()
	var reactions []Reaction
	repo.db.Where("image_id = ?", imageID).Find(&reactions)
	return reactions, nil
}
