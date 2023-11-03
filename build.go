//go:build ignore

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
	var (
		binPath     string          = "/go/bin/gfstat"
		workDir     string          = "/go/src/gfstat"
		exportedTar string          = "./gfstat.tar"
		imageName   string          = "gfstat:" + os.Getenv("GFV")
		ctx         context.Context = context.Background()
		err         error
	)
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	srcCode := client.Host().Directory(".")

	golang := client.Container().From("golang:alpine").
		WithDirectory(workDir, srcCode).
		WithWorkdir(workDir).
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"go", "get", "-d", "-v", "./..."}).
		WithExec([]string{"go", "build", "-o", binPath, "-v", "./..."})

	prodImg := client.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "ca-certificates"}).
		WithFile("/gfstat", golang.File(binPath)).
		WithDirectory("/views", golang.Directory(workDir+"/views")).
		WithExposedPort(3639).WithDefaultArgs(dagger.ContainerWithDefaultArgsOpts{
		Args: []string{"./gfstat"},
	})

	tarImage, err := prodImg.Export(ctx, exportedTar)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading image: ", tarImage)
	cmd := exec.Command("docker", "load", "-i", exportedTar)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	imageID := strings.Split(string(output), ":")[2]
	imageID = strings.Split(imageID, "\n")[0]

	fmt.Println("Tagging image: ", imageName)
	if err = exec.Command("docker", "tag", imageID, imageName).Run(); err != nil {
		fmt.Println("Error tagging image: ", err)
	}

	if err = os.Remove(exportedTar); err != nil {
		fmt.Println("Error removing tar file: ", err)
	}

	fmt.Println("docker run --rm -it -p 3639:3639 --env-file .env ", imageName)
}
