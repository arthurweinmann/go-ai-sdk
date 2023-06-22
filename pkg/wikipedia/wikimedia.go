package wikipedia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// A Wikimedia API response
type ApiResponse struct {
	Query         ApiQuery         `json:"query"`
	QueryContinue ApiQueryContinue `json:"query-continue"`
	Parse         ApiParse         `json:"parse"`
}

// Strip all HTML tags from the response
func (a *ApiResponse) StripHtml() {
	for k, v := range a.Query.Pages {
		v.Title = stripHtml(v.Title)
		v.Extract = stripHtml(v.Extract)
		a.Query.Pages[k] = v
	}
	for i, v := range a.Query.Search {
		v.Title = stripHtml(v.Title)
		v.Snippet = stripHtml(v.Snippet)
		a.Query.Search[i] = v
	}
}

type ApiQuery struct {
	Pages      map[string]ApiPage `json:"pages"`
	Search     []ApiSearch        `json:"search"`
	SearchInfo ApiSearchInfo      `json:"searchinfo"`
}

type ApiPage struct {
	PageId  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type ApiSearch struct {
	Ns        int       `json:"ns"`
	Title     string    `json:"title"`
	Snippet   string    `json:"snippet"`
	Size      int       `json:"size"`
	WordCount int       `json:"wordcount"`
	Timestamp time.Time `json:"timestamp"`
}

type ApiSearchInfo struct {
	Totalhits int `json:"totalhits"`
}

type ApiQueryContinue struct {
	Search ApiQueryContinueSearch `json:"search"`
}

type ApiQueryContinueSearch struct {
	SrOffset int `json:"sroffset"`
}

type ApiParse struct {
	PageId     int               `json:"pageid"`
	Title      string            `json:"title"`
	Categories []ApiPageCategory `json:"categories"`
	Sections   []ApiPageSection  `json:"sections"`
}

type ApiPageCategory struct {
	SortKey string `json:"sortkey"`
	Hidden  string `json:"hidden"`
	Name    string `json:"*"`
}

type ApiPageSection struct {
	TOCLevel   int    `json:"toclevel"`
	Level      string `json:"level"`
	Line       string `json:"line"`
	Number     string `json:"number"`
	Index      string `json:"index"`
	FromTitle  string `json:"fromtitle"`
	ByteOffset int    `json:"byteoffset"`
	Anchor     string `json:"anchor"`
}

func stripHtml(s string) string {
	var rs []rune
	in := false
	for _, v := range s {
		if in {
			if v == '>' {
				in = false
			}
			continue
		}
		if v == '<' {
			in = true
			continue
		}
		rs = append(rs, v)
	}
	return string(rs)
}

// A Wikimedia API client
type Wikimedia struct {
	// Full URL of the Wikimedia API, e.g. url.Parse("http://en.wikipedia.org/w/api.php")
	Url *url.URL

	// Automatically strip HTML tags from API responses
	StripHtml bool

	// HTTP client to use (defaults to http.DefaultClient)
	Client *http.Client

	// User-Agent header to provide
	UserAgent string

	url string
}

func (w *Wikimedia) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if w.UserAgent != "" {
		req.Header.Add("User-Agent", w.UserAgent)
	}
	if w.Client != nil {
		return w.Client.Do(req)
	}
	return http.DefaultClient.Do(req)
}

// Queries the Wikimedia API using the specified values, and returns an
// ApiResponse. See http://en.wikipedia.org/w/api.php for a reference.
func (w *Wikimedia) Query(vals url.Values) (*ApiResponse, error) {
	vals["format"] = []string{"json"}
	if w.url == "" {
		w.url = w.Url.String()
	}
	u := fmt.Sprintf("%s?%s", w.url, vals.Encode())
	res, err := w.get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var api ApiResponse
	err = json.NewDecoder(res.Body).Decode(&api)
	if w.StripHtml {
		api.StripHtml()
	}
	return &api, nil
}

// Set up a client that queries the specified API, e.g.
// http://en.wikipedia.org/w/api.php or http://da.wiktionary.org/w/api.php.
// Returns an error if the URL is invalid.
func New(apiUrl string) (*Wikimedia, error) {
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	w := &Wikimedia{
		Url: u,
	}
	return w, nil
}
