package main

import (
	"log"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"
)

var (
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	accessToken  = kingpin.Flag("token", "Github personal access token.").Envar("ESTAFETTE_EXTENSION_TOKEN").Required().String()
	gitRepoOwner = kingpin.Flag("git-repo-owner", "The owner of the Github repository.").Envar("ESTAFETTE_GIT_OWNER").Required().String()
	gitRepoName  = kingpin.Flag("git-repo-name", "The name of the Github repository.").Envar("ESTAFETTE_GIT_NAME").Required().String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// log to stdout and hide timestamp
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// log startup message
	log.Printf("Starting estafette-extension-github-pages version %v...", version)

	// set build status
	githubAPIClient := newGithubAPIClient(*accessToken)

	// get milestone by version
	err := githubAPIClient.RequestPageBuild(*gitRepoOwner, *gitRepoName)
	if err != nil {
		log.Fatalf("Requesting page build failed. %v", err)
	}

	log.Println("\nFinished estafette-extension-github-pages...")
}
