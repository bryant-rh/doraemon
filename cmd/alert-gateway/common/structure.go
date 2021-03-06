package common

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

var ErrHttpRequest = errors.New("create HTTP request failed")
var Maintain map[string]bool
var RuleCount map[[2]int64]int64
var Recover2Send = map[string]map[[2]int64]*Ready2Send{
	"LANXIN": map[[2]int64]*Ready2Send{},
	//"HOOK":   map[[2]int64]*Ready2Send{},
}

var Lock sync.Mutex
var Rw sync.RWMutex

// AuthModel holds information used to authenticate.
type AuthModel struct {
	Username string
	Password string
}

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type BrokenList struct {
	Hosts []struct {
		Hostname string `json:"hostname"`
	} `json:"hosts"`
	Error interface{} `json:"error"`
}

type Msg struct {
	Content string   `json:"content"`
	From    string   `json:"from"`
	Title   string   `json:"title"`
	To      []string `json:"to"`
}

type SingleAlert struct {
	Id       int64             `json:"id"`
	Count    int               `json:"count"`
	Value    float64           `json:"value"`
	Summary  string            `json:"summary"`
	Hostname string            `json:"hostname"`
	Labels   map[string]string `json:"labels"`
}

type Ready2Send struct {
	RuleId int64
	Start  int64
	User   []string
	Alerts []SingleAlert
}

type UserGroup struct {
	Id                    int64
	StartTime             string
	EndTime               string
	Start                 int
	Period                int
	ReversePolishNotation string
	User                  string
	Group                 string
	DutyGroup             string
	Method                string
}

type Alert []struct {
	ActiveAt    time.Time `json:"active_at"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
		RuleId      string `json:"rule_id"`
	} `json:"annotations"`
	FiredAt    time.Time         `json:"fired_at"`
	Labels     map[string]string `json:"labels"`
	LastSentAt time.Time         `json:"last_sent_at"`
	ResolvedAt time.Time         `json:"resolved_at"`
	State      int               `json:"state"`
	ValidUntil time.Time         `json:"valid_until"`
	Value      float64           `json:"value"`
}

type AlertForShow struct {
	Id              int64             `json:"id,omitempty"`
	RuleId          int64             `json:"rule_id"`
	Labels          map[string]string `json:"labels"`
	Value           float64           `json:"value"`
	Count           int               `json:"count"`
	Status          int8              `json:"status"`
	Summary         string            `json:"summary"`
	Description     string            `json:"description"`
	ConfirmedBy     string            `json:"confirmed_by"`
	FiredAt         *time.Time        `json:"fired_at"`
	ConfirmedAt     *time.Time        `json:"confirmed_at"`
	ConfirmedBefore *time.Time        `json:"confirmed_before"`
	ResolvedAt      *time.Time        `json:"resolved_at"`
}

type Confirm struct {
	Duration int
	User     string
	Ids      []int
}

type ValidUserGroup struct {
	User      string
	Group     string
	DutyGroup string
}

func GenerateJsonHeader() map[string]string {
	return map[string]string {
		"Content-Type": "application/json",
	}
}

func HttpPost(url string, params map[string]string, headers map[string]string, body []byte) (*http.Response, error) {
	//new request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, ErrHttpRequest
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{Timeout: 5 * time.Second} //Add the timeout,the reason is that the default client has no timeout set; if the remote server is unresponsive, you're going to have a bad day.
	return client.Do(req)
}

func HttpGet(url string, params map[string]string, headers map[string]string) (*http.Response, error) {
	//new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, ErrHttpRequest
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{Timeout: 5 * time.Second} //Add the timeout,the reason is that the default client has no timeout set; if the remote server is unresponsive, you're going to have a bad day.
	return client.Do(req)
}
