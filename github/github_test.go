package github

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var SearchQueryNoDate = SearchQuery{
	Repo:        "CapstoneLabs/slick",
	Labels:      []string{"bug"},
	ClosedSince: "",
}

var SearchQueryWithDate = SearchQuery{
	Repo:        "CapstoneLabs/slick",
	Labels:      []string{"bug"},
	ClosedSince: "2018-10-03T18:16:21Z",
}

var user = GHUser{Login: "3sky"}

var IssueEventClose = IssueEvent{
	Actor: user,
	Event: "closed",
}
var IssueItemNoNumber = IssueItem{
	Assignee: user,
	Events:   []IssueEvent{IssueEventClose},
}

var IssueEventOpen = IssueEvent{
	Actor: user,
	Event: "open",
}

var IssueItemWithNumber = IssueItem{
	Number:   1,
	Assignee: user,
	Events:   []IssueEvent{IssueEventOpen},
}

var config = Conf{
	Repos: []string{"CapstoneLabs/slick"},
}

var TestClient = Client{Conf: config}

func TestUrl(t *testing.T) {
	test1 := SearchQueryNoDate.Url()
	test2 := SearchQueryWithDate.Url()
	assert.Equal(t, test1, "https://api.github.com/search/issues?q=+repo:CapstoneLabs/slick+label:bug")
	assert.Equal(t, test2, "https://api.github.com/search/issues?q=+repo:CapstoneLabs/slick+label:bug+closed:>2018-10-03T18:16:21Z")
}

func TestLastClosedBy(t *testing.T) {
	test1 := IssueItemNoNumber.LastClosedBy()
	test2 := IssueItemWithNumber.LastClosedBy()
	assert.Equal(t, test1, "3sky")
	assert.Equal(t, test2, "")
}

func TestGet(t *testing.T) {
	r, _ := TestClient.Get("https://api.github.com/repos/CapstoneLabs/slick/issues/1")
	//366003717 ID of first issue
	assert.Contains(t, string(r), "366003717")
}

func TestDoSearchQuery(t *testing.T) {
	QueryResult, _ := TestClient.DoSearchQuery(SearchQueryWithDate)
	//Get first close issues since "2018-10-03T18:16:21Z" - #13
	FirstClosedIssue := fmt.Sprintf("%+v", QueryResult[len(QueryResult)-1])
	assert.Contains(t, FirstClosedIssue, "Fix typos")
}

func TestDoEventQuery(t *testing.T) {
	ch := make(chan IssueItem)
	go TestClient.DoEventQuery([]IssueItem{IssueItemWithNumber}, "CapstoneLabs/slick", ch)
	x := <-ch
	EventID := x.Events[0].ID
	//ID of First Event on #1 issue in CapstoneLabs/slick
	assert.Equal(t, EventID, 1879968249)
}
