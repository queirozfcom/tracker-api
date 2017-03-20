package main

import (
	"fmt"
	"os"
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

func main() {

	tok,exists := os.LookupEnv("GITHUB_TOKEN")

	if ! exists {
		panic("need a github token")
	}


	ctx := context.Background()

	ts := oauth2.StaticTokenSource(

		// need to specifically set TokenType: Basic
		&oauth2.Token{AccessToken: tok, TokenType: "Basic"},

		// this generates a 401 Bad Credentials error
		//&oauth2.Token{AccessToken: tok},

	)

	tc := oauth2.NewClient(ctx,ts)

	githubClient := github.NewClient(tc)


	repos, _, err := githubClient.Repositories.List(ctx, "", nil)

	fmt.Println(tok)
	fmt.Println(repos)
	fmt.Println(err)

}
