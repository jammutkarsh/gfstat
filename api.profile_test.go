package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

const testInput = "./testData/test.input.json"

var input struct {
	Username  string `json:"username"`
	Followers []MetaFollow `json:"followers"`
	Following []MetaFollow `json:"following"`
}

type testFields struct {
	Followers []MetaFollow
	Following []MetaFollow
}

func init() {
	testDataReader(&input, testInput)
}

func testDataReader(v any, file string) {
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
	testDataReader(&want, "./testData/test.output.mutuals.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		if got := Mutuals(f.Followers, f.Following); !reflect.DeepEqual(got, want) {
			t.Errorf("\nMutuals() = %v\nwant %v", got, want)
		}
	})
}

func TestFollowersYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	testDataReader(&want, "./testData/test.output.iDontFollow.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		if got := FollowersYouDontFollow(f.Followers, f.Following); !reflect.DeepEqual(got, want) {
			t.Errorf("\nFollowersYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}

func TestFollowingYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	testDataReader(&want, "./testData/test.output.theyDontFollow.json")
	f := testFields{
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		if got := FollowingYouDontFollow(f.Followers, f.Following); !reflect.DeepEqual(got, want) {
			t.Errorf("\nFollowingYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}
