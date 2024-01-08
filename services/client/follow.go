package client

import (
	"sort"

	"github.com/JammUtkarsh/gfstat/services/core"
	"github.com/google/go-github/github"
)

func GETFollowers(c *github.Client, u github.User) (followers []core.MetaFollow, err error) {
	if err := followOverflow(u); err != nil {
		return nil, err
	}
	for pageCount := 1; pageCount <= numberOfPages(*u.Followers); pageCount++ {
		opts := &github.ListOptions{Page: pageCount, PerPage: 100}

		follow, res, err := c.Users.ListFollowers(internalGitHubCtx, *u.Login, opts)
		if err != nil {
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

func GETFollowing(c *github.Client, u github.User) (followers []core.MetaFollow, err error) {
	if err := followOverflow(u); err != nil {
		return nil, err
	}
	for pageCount := 1; pageCount <= numberOfPages(*u.Following); pageCount++ {
		opts := &github.ListOptions{Page: pageCount, PerPage: 100}

		follow, res, err := c.Users.ListFollowing(internalGitHubCtx, *u.Login, opts)
		if err != nil {
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

func numberOfPages(followCount int) int {
	if followCount < 100 {
		return 1
	}
	return (followCount / 100) + 1
}
