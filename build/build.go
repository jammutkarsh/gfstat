package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"dagger.io/dagger"
)

var (
	imageTag  = os.Getenv("IMAGE_TAG") // TODO: Get this from git tag
	imageName = "registry.digitalocean.com/jammutkarsh/gfstat:" + imageTag
	ctx       = context.Background()
)

func main() {
	daggerClient, err := dagger.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer daggerClient.Close()
	if err := generateDockerImage(daggerClient); err != nil {
		log.Fatal(err)
	}
}

// generateDockerImage generates a container image similar to Dockerfile in root.
func generateDockerImage(client *dagger.Client) error {
	var (
		projectName  = "gfstat"
		binPath      = "/go/bin/" + projectName
		workDir      = "/go/src/" + projectName
		builderImage = "golang:1.22-alpine"
		production   = "alpine:3.19"
	)
	srcCode := client.Host().Directory(".")

	// Docker file instructions
	golang := client.Container().From(builderImage).
		WithDirectory(workDir, srcCode).
		WithWorkdir(workDir).
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"go", "get", "-d", "./..."}).
		WithExec([]string{"go", "build", "-o", binPath, "."})

	prodImg := client.Container().From(production).
		WithExec([]string{"apk", "add", "--no-cache", "ca-certificates"}).
		WithFile("/gfstat", golang.File(binPath)).
		WithDirectory("/views", golang.Directory(workDir+"/views")).
		WithExposedPort(3639).WithDefaultArgs([]string{"./gfstat"})
	log.Println("Image built successfully")

	tarFile := projectName + ".tar"
	// Export blog as tar file
	if _, err := prodImg.Export(ctx, tarFile); err != nil {
		return err
	} else {
		log.Println("Image exported as tar file")
	}
	defer removeFile(tarFile)

	imageID, err := loadContainerImage(tarFile)
	if err != nil {
		return err
	}

	// Tag image
	if err := exec.Command("docker", "image", "tag", imageID, imageName).Run(); err != nil {
		return err
	} else {
		log.Println("Image tagged as: ", projectName)
	}
	return nil
}

func removeFile(fileName string) {
	if err := os.Remove(fileName); err != nil {
		log.Printf("Error removing file %s: %s", fileName, err)
	} else {
		log.Printf("%s removed\n", fileName)
	}
}

func loadContainerImage(projectName string) (string, error) {
	cmdOutput, err := exec.Command("docker", "image", "load", "-i", projectName).Output()
	if err != nil {
		return "", err
	} else {
		log.Println("Image loaded in Docker")
	}
	return strings.Split(strings.Split(string(cmdOutput), ":")[2], "\n")[0], nil
}
