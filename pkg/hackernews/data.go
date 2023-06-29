package hackernews

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetLatestItemId returns the current largest item id from the Hacker News API.
func GetLatestItemId() (int, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/maxitem.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var latestItemId int
	err = json.Unmarshal(body, &latestItemId)
	if err != nil {
		return 0, err
	}

	return latestItemId, nil
}

// GetStoryIds retrieves story ids from the given url.
func GetStoryIds(url string) ([]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var storyIds []int
	err = json.Unmarshal(body, &storyIds)
	if err != nil {
		return nil, err
	}

	return storyIds, nil
}

// GetNewStories retrieves new story ids.
func GetNewStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/newstories.json")
}

// GetTopStories retrieves top story ids.
func GetTopStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/topstories.json")
}

// GetBestStories retrieves best story ids.
func GetBestStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/beststories.json")
}

// GetJobStories retrieves job story ids.
func GetJobStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/jobstories.json")
}

// GetAskStories retrieves ask story ids.
func GetAskStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/askstories.json")
}

// GetShowStories retrieves show story ids.
func GetShowStories() ([]int, error) {
	return GetStoryIds("https://hacker-news.firebaseio.com/v0/showstories.json")
}
