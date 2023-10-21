package main

import (
	"sort"

	"github.com/google/go-github/github"
)

func GETFollowers(c *github.Client, u github.User) (followers []MetaFollow, err error) {
	pageCount := 1
	opts := &github.ListOptions{Page: pageCount, PerPage: 100}
	follow, res, err := c.Users.ListFollowers(internalGitHubCtx, *u.Name, opts)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	for _, v := range follow {
		followers = append(followers, MetaFollow{
			Username: *v.Login,
			HTMLURL:  *v.HTMLURL,
		})
	}

	sort.Slice(followers, func(i, j int) bool {
		return followers[i].Username < followers[j].Username
	})

	return followers, nil
}

func GETFollowing(c *github.Client, u github.User) (followers []MetaFollow, err error) {
	pageCount := 1
	opts := &github.ListOptions{Page: pageCount, PerPage: 100}
	follow, res, err := c.Users.ListFollowing(internalGitHubCtx, *u.Name, opts)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	for _, v := range follow {
		followers = append(followers, MetaFollow{
			Username: *v.Login,
			HTMLURL:  *v.HTMLURL,
		})
	}

	sort.Slice(followers, func(i, j int) bool {
		return followers[i].Username < followers[j].Username
	})

	return followers, nil

}
