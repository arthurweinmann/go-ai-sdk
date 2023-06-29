package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	// The user's unique username. Case-sensitive.
	Id string `json:"id"`
	// Creation date of the user, in Unix Time.
	Created int64 `json:"created"`
	// The user's karma.
	Karma int `json:"karma"`
	// The user's optional self-description. HTML.
	About string `json:"about"`
	// List of the user's stories, polls and comments.
	Submitted []int `json:"submitted"`
}

func GetUser(userId string) (*User, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/user/%s.json", userId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
