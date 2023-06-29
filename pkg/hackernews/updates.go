package hackernews

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Updates represents the updated items and profiles.
type Updates struct {
	Items    []int    `json:"items"`    // List of updated item IDs
	Profiles []string `json:"profiles"` // List of updated profile usernames
}

// GetUpdates retrieves the updates.
func GetUpdates() (Updates, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/updates.json")
	if err != nil {
		return Updates{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Updates{}, err
	}

	var updates Updates
	err = json.Unmarshal(body, &updates)
	if err != nil {
		return Updates{}, err
	}

	return updates, nil
}
