package main

import (
	"runtime"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

var (
	appgroup  string
	app       string
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

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	// set build status
	githubAPIClient := newGithubAPIClient(*accessToken)

	// get milestone by version
	err := githubAPIClient.RequestPageBuild(*gitRepoOwner, *gitRepoName)
	if err != nil {
		log.Fatal().Err(err).Msg("Requesting page build failed")
	}

	log.Info().Msg("Finished estafette-extension-github-pages...")
}
