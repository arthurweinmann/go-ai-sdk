package hackernews

type FullStory struct {
	Story    *Item   `json:"story"`
	Comments []*Item `json:"comments"`
}

func FetchFullStory(story *Item) (*FullStory, error) {
	processedComments := make(map[int]bool)
	comments, err := fetchComments(story.Kids, processedComments)
	if err != nil {
		return nil, err
	}

	swc := &FullStory{
		Story:    story,
		Comments: comments,
	}

	return swc, nil
}

func fetchComments(ids []int, processedComments map[int]bool) ([]*Item, error) {
	comments := make([]*Item, 0, len(ids))
	for _, id := range ids {
		if _, ok := processedComments[id]; ok {
			continue
		}
		comment, err := GetItem(id)
		if err != nil {
			return nil, err
		}
		processedComments[id] = true
		comments = append(comments, comment)
		if len(comment.Kids) > 0 {
			childComments, err := fetchComments(comment.Kids, processedComments)
			if err != nil {
				return nil, err
			}
			comments = append(comments, childComments...)
		}
	}
	return comments, nil
}
