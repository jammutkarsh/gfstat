package core

import (
	"encoding/json"
	"os"
	"reflect"
	"sync"
	"testing"
)

const testInput = "../../testData/test.input.json"

var input struct {
	Username  string       `json:"username"`
	Followers []MetaFollow `json:"followers"`
	Following []MetaFollow `json:"following"`
}

type testFields struct {
	Followers []MetaFollow
	Following []MetaFollow
	chann     chan []MetaFollow
	wg        sync.WaitGroup
}

func init() {
	prepareTests(&input, testInput)
}

func prepareTests(v any, file string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic("error in reading test file:" + err.Error() + "\n")
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		panic("error in converting test file to JSON:" + err.Error() + "\n")
	}
}

func TestMutuals(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	prepareTests(&want, "../../testData/test.output.mutuals.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
		chann:     make(chan []MetaFollow), // Add channel initialization
		wg:        sync.WaitGroup{},
	}
	t.Run("Test 1", func(t *testing.T) {
		f.wg.Add(1)
		go Mutuals(f.Followers, f.Following, f.chann, &f.wg)
		got := <-f.chann // Receive value from channel
		close(f.chann)   // Close the channel after receiving the value
		f.wg.Wait()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\nMutuals() = %v\nwant %v", got, want)
		}
	})
}

func TestFollowersYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	prepareTests(&want, "../../testData/test.output.iDontFollow.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
		chann:     make(chan []MetaFollow),
		wg:        sync.WaitGroup{},
	}
	t.Run("Test 1", func(t *testing.T) {
		f.wg.Add(1)
		go IDontFollow(f.Followers, f.Following, f.chann, &f.wg)
		got := <-f.chann
		close(f.chann)
		f.wg.Wait()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\nFollowersYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}

func TestFollowingYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	prepareTests(&want, "../../testData/test.output.theyDontFollow.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
		chann:     make(chan []MetaFollow),
		wg:        sync.WaitGroup{},
	}
	t.Run("Test 1", func(t *testing.T) {
		f.wg.Add(1)
		go TheyDontFollow(f.Followers, f.Following, f.chann, &f.wg)
		got := <-f.chann
		close(f.chann)
		f.wg.Wait()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\nFollowingYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}
