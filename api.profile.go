package main

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
func Mutuals(followers, following []MetaFollow) []MetaFollow {
	if len(followers) == 0 || len(following) == 0 {
		return nil
	}
	var mutuals []MetaFollow
	for _, follower := range followers {
		low, high := 0, len(following)-1
		for low <= high {
			mid := low + (high-low)/2
			if follower.Username == following[mid].Username {
				mutuals = append(mutuals, follower)
				break
			} else if follower.Username < following[mid].Username {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return mutuals
}

// followers - following
// TC: O(n)
// SC: O(N)
func IDontFollow(followers, following []MetaFollow) []MetaFollow {
	m := make(map[string]MetaFollow)
	for _, following := range following {
		m[following.Username] = following
	}

	var iDontFollow []MetaFollow
	for _, follower := range followers {
		if _, ok := m[follower.Username]; !ok {
			iDontFollow = append(iDontFollow, follower)
		}
	}
	return iDontFollow
}

// following - followers
// TC: O(n)
// SC: O(N)
func TheyDontFollow(followers, following []MetaFollow) []MetaFollow {
	m := make(map[string]MetaFollow)
	for _, follower := range followers {
		if _, ok := m[follower.Username]; !ok {
			m[follower.Username] = follower
		}
	}

	var iDontFollow []MetaFollow
	for _, following := range following {
		if _, ok := m[following.Username]; !ok {
			iDontFollow = append(iDontFollow, following)
		}
	}

	return iDontFollow
}
