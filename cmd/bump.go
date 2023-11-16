package cmd

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(bumpVersionCmd)

	bumpVersionCmd.Flags().StringP("repository", "r", ".", "path to the repository")
	bumpVersionCmd.Flags().StringSliceP("include", "i", []string{"."}, "the subfolders which you want the commits analysed in")         //nolint:lll
	bumpVersionCmd.Flags().StringP("current-version", "c", "", "the current version of the application we want to bump the version on") //nolint:lll

	bumpVersionCmd.MarkFlagRequired("repository")
	bumpVersionCmd.MarkFlagRequired("current-version")
}

var bumpVersionCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bumps the version of the application",
	Run:   bumpVersionCmdRunE,
}

const numVersionParts int = 3

type commitType string

const (
	patch commitType = "fix"
	minor commitType = "feature"
	major commitType = "breaking"
	none  commitType = "none"
)

type commitPrefix string

const (
	fix          commitPrefix = "fix:"
	fixBreaking  commitPrefix = "fix!:"
	feat         commitPrefix = "feat:"
	featBreaking commitPrefix = "feat!:"
)

func bumpVersionCmdRunE(cmd *cobra.Command, args []string) {
	// get cli values
	currentVersion, _ := cmd.Flags().GetString("current-version")
	repository, _ := cmd.Flags().GetString("repository")
	include, _ := cmd.Flags().GetStringSlice("include")

	// Get the last tag, we will use the commits since this tag to determine
	// what kind of version increment we will do
	lastTag, err := getLastTag(repository)
	if err != nil {
		log.Fatal(err)
	}

	// Get the commit messages since the last tag
	commitMessages, err := getCommitMessagesSince(lastTag, repository, include...)
	if err != nil {
		log.Fatal(err)
	}

	// Check the commit messages and determine the semantic version change type
	versionChange := determineVersionChangeType(commitMessages)

	// Increment the semantic version
	newVersion, err := incrementSemanticVersion(currentVersion, versionChange)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newVersion)
}

func getLastTag(repository string) (string, error) {
	cmd := exec.Command("git", "-C", repository, "describe", "--abbrev=0", "--tags")

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func getCommitMessagesSince(tag string, repository string, subfolders ...string) ([]string, error) {
	args := []string{"-C", repository, "log", "--pretty=format:%s", tag + "..HEAD", "--"}
	args = append(args, subfolders...)

	cmd := exec.Command("git", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func determineVersionChangeType(commitMessages []string) commitType {
	var commitTypes []commitType

	if len(commitMessages) == 0 {
		return none
	}

	// go over each commit and check what type of change it was
	for _, commitMessage := range commitMessages {
		var commitType commitType

		switch {
		case isRegexMatch(commitMessage, `^(fix|feat)(?:\([^\)]+\))?!:.*$`):
			// we can return early here since a breaking change is a major
			// release and overrides all other change types
			return major
		case isRegexMatch(commitMessage, `^feat(?:\([^\)]+\))?(!:)?.*$`):
			commitType = minor
		case isRegexMatch(commitMessage, `^fix(?:\([^\)]+\))?(!:)?.*$`):
			commitType = patch
		default:
			commitType = none
		}

		commitTypes = append(commitTypes, commitType)
	}

	// if we got this farthen its not a breaking change so now check if it
	// contains any minor changes
	if slices.Contains(commitTypes, minor) {
		return minor
	}

	// if we got this far there are no breaking or minor changes, so check if it
	// contains any patches
	if slices.Contains(commitTypes, patch) {
		return patch
	}

	// if we got this far then none of the commits told us what type of change
	// it was
	return none
}

func incrementSemanticVersion(currentVersion string, versionChange commitType) (string, error) {
	parts := strings.Split(currentVersion, ".")
	if len(parts) != numVersionParts {
		return "", errors.New("invalid version format")
	}

	// Convert each part to an integer
	majorInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}

	minorInt, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}

	patchInt, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", err
	}

	// Increment the corresponding part based on the version change
	switch versionChange {
	case major:
		majorInt++

		minorInt = 0
		patchInt = 0
	case minor:
		minorInt++

		patchInt = 0
	case patch:
		patchInt++
	case none:
		// dont do anything
	default:
		return "", errors.New("invalid version change type")
	}

	// Format the new version string
	newVersion := fmt.Sprintf("%d.%d.%d", majorInt, minorInt, patchInt)

	return newVersion, nil
}

func isRegexMatch(input, pattern string) bool {
	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false
	}

	return match
}
