package core

import "sync"

type AccountType int

const (
	User AccountType = iota
	Organization
)

// MetaFollow stands for metadata of a follower/following.
type MetaFollow struct {
	Username string `json:"username"`
	HTMLURL  string `json:"html_url"`
}

// Mutuals gives the list of mutuals between followers and following.
// TC: O(nLogn)
// SC: O(1)
func Mutuals(followers, following []MetaFollow, c chan []MetaFollow, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(followers) == 0 || len(following) == 0 {
		return
	}
	var results []MetaFollow
	for _, follower := range followers {
		low, high := 0, len(following)-1
		for low <= high {
			mid := low + (high-low)/2
			if follower.Username == following[mid].Username {
				results = append(results, follower)
				break
			} else if follower.Username < following[mid].Username {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	c <- results
}

// followers - following
// TC: O(n)
// SC: O(N)
func IDontFollow(followers, followings []MetaFollow, c chan []MetaFollow, wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]MetaFollow)
	for _, following := range followings {
		m[following.Username] = following
	}

	var results []MetaFollow
	for _, follower := range followers {
		if _, ok := m[follower.Username]; !ok {
			results = append(results, follower)
		}
	}
	c <- results
}

// following - followers
// TC: O(n)
// SC: O(N)
func TheyDontFollow(followers, followings []MetaFollow, c chan []MetaFollow, wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]MetaFollow)
	for _, follower := range followers {
		if _, ok := m[follower.Username]; !ok {
			m[follower.Username] = follower
		}
	}

	var results []MetaFollow
	for _, following := range followings {
		if _, ok := m[following.Username]; !ok {
			results = append(results, following)
		}
	}
	c <- results
}
