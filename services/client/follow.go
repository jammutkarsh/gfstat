package client

import (
	"fmt"
	"log"
	"sort"

	"github.com/JammUtkarsh/gfstat/services/core"
	"github.com/google/go-github/github"
)

const (
	pageLimit = 100
	reqLimit  = 5000 - 2
)

func GETFollowers(c *github.Client, u github.User, resp github.Response) (followers []core.MetaFollow, err error) {
	if count := *u.Followers + *u.Following; count/pageLimit > reqLimit {
		return nil, fmt.Errorf("you have too many followers and following: %d", count)
	}
	for i := 0; i <= resp.LastPage; i++ {
		follow, res, err := c.Users.ListFollowers(ctx, *u.Login, &github.ListOptions{Page: i, PerPage: pageLimit})
		if err != nil {
			log.Println(err)
			return nil, err
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, err
		}

		for _, v := range follow {
			followers = append(followers, core.MetaFollow{
				Username: *v.Login,
				HTMLURL:  *v.HTMLURL,
			})
		}

	}

	sort.Slice(followers, func(i, j int) bool {
		return followers[i].Username < followers[j].Username
	})

	return followers, nil
}

func GETFollowing(c *github.Client, u github.User, resp github.Response) (followings []core.MetaFollow, err error) {
	if count := *u.Followers + *u.Following; count/pageLimit > reqLimit {
		return nil, fmt.Errorf("you have too many followers and following: %d", count)
	}
	for i := 0; i <= resp.LastPage; i++ {
		follow, res, err := c.Users.ListFollowing(ctx, *u.Login, &github.ListOptions{Page: i, PerPage: pageLimit})
		if err != nil {
			log.Println(err)
			return nil, err
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, err
		}

		for _, v := range follow {
			followings = append(followings, core.MetaFollow{
				Username: *v.Login,
				HTMLURL:  *v.HTMLURL,
			})
		}

	}

	sort.Slice(followings, func(i, j int) bool {
		return followings[i].Username < followings[j].Username
	})

	return followings, nil

}
