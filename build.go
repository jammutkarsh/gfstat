//go:build exclude

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	githubPublicID := os.Getenv("GH_BASIC_CLIENT_ID")     // like public key
	githubServerSecret := os.Getenv("GH_BASIC_SECRET_ID") // like private key

	if githubPublicID == "" || githubServerSecret == "" {
		panic("Environment variable GH_SECRET is not set")
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	s1 := client.SetSecret("GH_BASIC_CLIENT_ID", githubPublicID)
	s2 := client.SetSecret("GH_BASIC_SECRET_ID", githubServerSecret)
	buildPath := "/go/bin/gfstat"
	srcPath := "/go/src/gfstat"
	src := client.Host().Directory(".")
	golang := client.Container().From("golang:alpine").
		WithDirectory(srcPath, src).
		WithWorkdir(srcPath).
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"go", "get", "-d", "-v", "./..."}).
		WithExec([]string{"go", "build", "-o", buildPath, "-v", "./..."})
	
		prodImg := client.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "ca-certificates"}).
		WithSecretVariable("GH_BASIC_CLIENT_ID", s1).
		WithSecretVariable("GH_BASIC_SECRET_ID", s2).
		WithFile("/gfstat", golang.File(buildPath)).
		WithDirectory("/views", golang.Directory(srcPath+"/views")).
		WithExposedPort(3639).WithDefaultArgs(dagger.ContainerWithDefaultArgsOpts{
		Args: []string{"./gfstat"},
	})

	exportedTar := "./gfstat.tar"
	val, err := prodImg.Export(ctx, exportedTar)
	if err != nil {
		panic(err)
	}

	// print result
	fmt.Println("Exported image: ", val)
	fmt.Println(exportedTar + " is ready to be loaded in Docker")
	opt := "N"
	fmt.Println("Load it in Docker? Enter Y")
	fmt.Scanf("%s", &opt)
	if opt == "Y" {
		cmd := exec.Command("docker", "load", "-i", exportedTar)
		b, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		tarImage := strings.Split(string(b), ":")[2]
		// taggedImage := "gfstat:latest"
		fmt.Println()
		fmt.Println(tarImage)
	}
}
