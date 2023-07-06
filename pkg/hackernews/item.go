package hackernews

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ItemType string

const (
	JobType     ItemType = "job"
	StoryType   ItemType = "story"
	CommentType ItemType = "comment"
	PollType    ItemType = "poll"
	PollOptType ItemType = "pollopt"
)

const itemBaseURL = "https://hacker-news.firebaseio.com/v0/item/"

type Item struct {
	ID          int      `json:"id"`                    // The item's unique id.
	Deleted     bool     `json:"deleted,omitempty"`     // true if the item is deleted.
	Type        ItemType `json:"type"`                  // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By          string   `json:"by,omitempty"`          // The username of the item's author.
	Time        UnixTime `json:"time"`                  // Creation date of the item, in Unix Time.
	Text        string   `json:"text,omitempty"`        // The comment, story or poll text. HTML.
	Dead        bool     `json:"dead,omitempty"`        // true if the item is dead.
	Parent      int      `json:"parent,omitempty"`      // The comment's parent: either another comment or the relevant story.
	Poll        int      `json:"poll,omitempty"`        // The pollopt's associated poll.
	Kids        []int    `json:"kids,omitempty"`        // The ids of the item's comments, in ranked display order.
	URL         string   `json:"url,omitempty"`         // The URL of the story.
	Score       int      `json:"score,omitempty"`       // The story's score, or the votes for a pollopt.
	Title       string   `json:"title,omitempty"`       // The title of the story, poll or job. HTML.
	Parts       []int    `json:"parts,omitempty"`       // A list of related pollopts, in display order.
	Descendants int      `json:"descendants,omitempty"` // In the case of stories or polls, the total comment count.
}

func GetItem(id int) (*Item, error) {
	resp, err := http.Get(itemBaseURL + strconv.Itoa(id) + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

// IterateStoriesByBatch retrieves stories by batch of a specified limit.
func IterateItemsByBatch(batchSize int, cb func(batch []*Item) (bool, error)) error {
	maxItem, err := GetLatestItemId()
	if err != nil {
		return err
	}

	var items []*Item
	for i := maxItem; i > -1; i -= batchSize {
		for j := i; j > i-batchSize && j > -1; j-- {
			item, err := GetItem(j)
			if err != nil {
				fmt.Printf("Error getting item with id %d: %v\n", j, err)
				continue
			}

			items = append(items, item)
		}

		continu, err := cb(items)
		if err != nil {
			return err
		}
		if !continu {
			return nil
		}

		items = items[:0]
	}

	return nil
}

type UnixTime int64

// UnmarshalJSON converts a Unix timestamp to UnixTime during JSON unmarshaling.
func (t *UnixTime) UnmarshalJSON(b []byte) error {
	var n int64
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	*t = UnixTime(n)
	return nil
}

// Time converts UnixTime to time.Time format.
func (t UnixTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}
