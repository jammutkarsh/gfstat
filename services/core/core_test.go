package core

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
)

type outputFields struct {
	Followers []MetaFollow
	Following []MetaFollow
	chann     chan []MetaFollow
	wg        sync.WaitGroup
}

var inputField struct {
	Username  string       `json:"username"`
	Followers []MetaFollow `json:"followers"`
	Following []MetaFollow `json:"following"`
}

func init() {
	if err := json.Unmarshal([]byte(testInput), &inputField); err != nil {
		panic(err)
	}
}

func TestMutuals(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	if err := json.Unmarshal([]byte(testMutuals), &want); err != nil {
		t.Fatalf("Error in unmarshalling output: %v", err)
	}
	f := outputFields{inputField.Followers, inputField.Following, make(chan []MetaFollow), sync.WaitGroup{}}
	f.wg.Add(1)
	defer f.wg.Wait()
	go Mutuals(f.Followers, f.Following, f.chann, &f.wg)
	if got := <-f.chann; !reflect.DeepEqual(got, want) {
		t.Errorf("\nMutuals() = %v\nwant %v", got, want)
	}
}

func TestFollowersYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	if err := json.Unmarshal([]byte(testIDontFollow), &want); err != nil {
		t.Fatalf("Error in unmarshalling output: %v", err)
	}
	f := outputFields{inputField.Followers, inputField.Following, make(chan []MetaFollow), sync.WaitGroup{}}
	f.wg.Add(1)
	defer f.wg.Wait()
	go IDontFollow(f.Followers, f.Following, f.chann, &f.wg)
	if got := <-f.chann; !reflect.DeepEqual(got, want) {
		t.Errorf("\nFollowersYouDontFollow() = %v\nwant %v", got, want)
	}
}

func TestFollowingYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	f := outputFields{inputField.Followers, inputField.Following, make(chan []MetaFollow), sync.WaitGroup{}}
	if err := json.Unmarshal([]byte(testTheyDontFollow), &want); err != nil {
		t.Fatalf("Error in unmarshalling output: %v", err)
	}
	f.wg.Add(1)
	defer f.wg.Wait()
	go TheyDontFollow(f.Followers, f.Following, f.chann, &f.wg)
	if got := <-f.chann; !reflect.DeepEqual(got, want) {
		t.Errorf("\nFollowingYouDontFollow() = %v\nwant %v", got, want)
	}
}

const (
	// Testcases
	testInput = `{"username":"soumyyyaaa","followers": [{"username":"JammUtkarsh","html_url":"https://github.com/JammUtkarsh"},{"username":"MadhaviGupta","html_url":"https://github.com/MadhaviGupta"},
	{"username":"Mishank24","html_url":"https://github.com/Mishank24"},{"username":"SparshGarg1","html_url":"https://github.com/SparshGarg1"},
	{"username":"SujalSamai","html_url":"https://github.com/SujalSamai"},{"username":"TheGameisYash","html_url":"https://github.com/TheGameisYash"},
	{"username":"dcdeepesh","html_url":"https://github.com/dcdeepesh"},{"username":"golemvincible","html_url":"https://github.com/golemvincible"},
	{"username":"shreyash2002","html_url":"https://github.com/shreyash2002"},{"username":"shristigupta12","html_url":"https://github.com/shristigupta12"},
	{"username":"sushantsharma08","html_url":"https://github.com/sushantsharma08"},{"username":"tanishjain158","html_url":"https://github.com/tanishjain158"}],
	"following":[{"username":"JammUtkarsh","html_url":"https://github.com/JammUtkarsh"},{"username":"MadhaviGupta","html_url":"https://github.com/MadhaviGupta"},
	{"username":"Mishank24","html_url":"https://github.com/Mishank24"},{"username":"SparshGarg1","html_url":"https://github.com/SparshGarg1"},
	{"username":"SujalSamai","html_url":"https://github.com/SujalSamai"},{"username":"TheGameisYash","html_url":"https://github.com/TheGameisYash"},
	{"username":"dcdeepesh","html_url":"https://github.com/dcdeepesh"},{"username":"dxaman","html_url":"https://github.com/dxaman"},
	{"username":"golemvincible","html_url":"https://github.com/golemvincible"},{"username":"shreyash2002","html_url":"https://github.com/shreyash2002"},
	{"username":"shristigupta12","html_url":"https://github.com/shristigupta12"},{"username":"sushantsharma08","html_url":"https://github.com/sushantsharma08"}]}`
	testMutuals = `[{"username":"JammUtkarsh","html_url":"https://github.com/JammUtkarsh"},{"username":"MadhaviGupta","html_url":"https://github.com/MadhaviGupta"},
	{"username":"Mishank24","html_url":"https://github.com/Mishank24"},{"username":"SparshGarg1","html_url":"https://github.com/SparshGarg1"},
	{"username":"SujalSamai","html_url":"https://github.com/SujalSamai"},{"username":"TheGameisYash","html_url":"https://github.com/TheGameisYash"},
	{"username":"dcdeepesh","html_url":"https://github.com/dcdeepesh"},{"username":"golemvincible","html_url":"https://github.com/golemvincible"},
	{"username":"shreyash2002","html_url":"https://github.com/shreyash2002"},{"username":"shristigupta12","html_url":"https://github.com/shristigupta12"},
	{"username":"sushantsharma08","html_url":"https://github.com/sushantsharma08"}]`
	testIDontFollow    = `[{"username":"tanishjain158","html_url":"https://github.com/tanishjain158"}]`
	testTheyDontFollow = `[{"username":"dxaman","html_url":"https://github.com/dxaman"}]`
)
