package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

const testInput = "./testData/test.input.json"

var input CurrentUser

func init() {
	testDataReader(&input, testInput)
}

type fields struct {
	Username  string
	Followers []MetaFollow
	Following []MetaFollow
}

func testDataReader(v any, file string) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		panic(err)
	}
}

func TestCurrentUser_Mutuals(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	testDataReader(&want, "./testData/test.output.mutuals.json")
	f := fields{
		Username:  input.Username,
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		c := CurrentUser{
			Followers: f.Followers,
			Following: f.Following,
		}
		if got := c.Mutuals(); !reflect.DeepEqual(got, want) {
			t.Errorf("\nCurrentUser.Mutuals() = %v\nwant %v", got, want)
		}
	})
}

func TestCurrentUser_FollowersYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	testDataReader(&want, "./testData/test.output.iDontFollow.json")
	f := fields{
		Username:  input.Username,
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		c := CurrentUser{
			Followers: f.Followers,
			Following: f.Following,
		}
		if got := c.FollowersYouDontFollow(); !reflect.DeepEqual(got, want) {
			t.Errorf("\nCurrentUser.FollowersYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}

func TestCurrentUser_FollowingYouDontFollow(t *testing.T) {
	t.Parallel()
	var want []MetaFollow
	testDataReader(&want, "./testData/test.output.theyDontFollow.json")
	f := fields{
		Username:  input.Username,
		Followers: input.Followers,
		Following: input.Following,
	}
	t.Run("Test 1", func(t *testing.T) {
		c := CurrentUser{
			Followers: f.Followers,
			Following: f.Following,
		}
		if got := c.FollowingYouDontFollow(); !reflect.DeepEqual(got, want) {
			t.Errorf("\nCurrentUser.FollowingYouDontFollow() = %v\nwant %v", got, want)
		}
	})
}
