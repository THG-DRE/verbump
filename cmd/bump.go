package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
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
	bumpVersionCmd.Flags().StringSliceP("include", "i", []string{"."}, "the subfolders in which you want the commits analysed") //nolint:lll
	bumpVersionCmd.Flags().StringP("version-file", "c", "", "the file which contains the version to be bumped")
	bumpVersionCmd.Flags().StringP("pre-release", "p", "", "the pre release label we want appended to the version")

	bumpVersionCmd.MarkFlagRequired("repository")
	bumpVersionCmd.MarkFlagRequired("version-file")

	// TODO: i want to be able to bump a beta build
	// bump --pre-release=beta
	// 0.5.0 -> 0.5.1-beta.0
	// bump --pre-release=beta
	// 0.5.1-beta.0 -> 0.5.1-beta.1
	// bump
	// 0.5.1-beta.1 -> 0.5.1
}

var bumpVersionCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bumps the version in the file",
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

const writeMode0644 fs.FileMode = 644

const git string = "git"

func bumpVersionCmdRunE(cmd *cobra.Command, args []string) {
	// get cli values
	versionFile, _ := cmd.Flags().GetString("version-file")
	repository, _ := cmd.Flags().GetString("repository")
	include, _ := cmd.Flags().GetStringSlice("include")
	preReleaseLabel, _ := cmd.Flags().GetString("pre-release")

	currentVersionBytes, err := os.ReadFile(versionFile)
	if err != nil {
		log.Fatal(err)
	}

	currentVersion := strings.Trim(string(currentVersionBytes), "\n")

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
	newVersion, err := incrementSemanticVersion(currentVersion, versionChange, preReleaseLabel)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(currentVersion)
	log.Print(newVersion)

	os.WriteFile(versionFile, []byte(newVersion), writeMode0644)
}

func getLastTag(repository string) (string, error) {
	app := git
	args := []string{"-C", repository, "describe", "--abbrev=0", "--tags"}

	command := exec.Command(app, args...)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	// we dont want check the error immediately, sometimes there is useful info
	// in stdout and stderr so try print them first
	err := command.Run()

	// maybe here we check if stdout matches a tag, if not then print and return
	if stdout.String() != "" {
		fmt.Printf("------------------------------ %s - stdout ------------------------------\n", app)
		fmt.Println(stdout.String())
		fmt.Printf("------------------------------ %s - stdout ------------------------------\n", app)
	}

	if stderr.String() != "" {
		fmt.Printf("------------------------------ %s - stderr ------------------------------\n", app)
		fmt.Println(stderr.String())
		fmt.Printf("------------------------------ %s - stderr ------------------------------\n", app)

		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

func getCommitMessagesSince(tag string, repository string, subfolders ...string) ([]string, error) {
	app := git
	args := []string{"-C", repository, "log", "--pretty=format:%s", tag + "..HEAD", "--"}
	args = append(args, subfolders...)

	command := exec.Command(app, args...)

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	// we dont want check the error immediately, sometimes there is useful info
	// in stdout and stderr so try print them first
	err := command.Run()

	// maybe here we check if stdout matches a version, if not then print and return
	if stdout.String() != "" {
		fmt.Printf("------------------------------ %s - stdout ------------------------------\n", app)
		fmt.Println(stdout.String())
		fmt.Printf("------------------------------ %s - stdout ------------------------------\n", app)
	}

	if stderr.String() != "" {
		fmt.Printf("------------------------------ %s - stderr ------------------------------\n", app)
		fmt.Println(stderr.String())
		fmt.Printf("------------------------------ %s - stderr ------------------------------\n", app)

		return []string{""}, err
	}

	return strings.Split(strings.TrimSpace(stdout.String()), "\n"), nil
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

func incrementSemanticVersion(currentVersion string, versionChange commitType, preReleaseLabel string) (string, error) {
	parts := splitAny(currentVersion, ".-")
	if len(parts) < numVersionParts {
		return "", errors.New("invalid version format")
	}

	var currentPreReleaseLabel string
	var preRelease string

	preReleaseLabelPresent := len(parts) >= 4
	if preReleaseLabelPresent {
		// if this version already contains a preReleaseLabel then
		// dont increment anything other than the preReleaseLabel
		versionChange = none
		currentPreReleaseLabel = parts[3]
		preRelease = parts[4]
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

	// if there is any sort of pre-release label applied then increment it if
	// it was already present or add it
	if preReleaseLabel != "" {
		var preReleaseInt int

		// if the version already contains the label then increment it
		if preReleaseLabelPresent && currentPreReleaseLabel == preReleaseLabel {
			preReleaseInt, err = strconv.Atoi(preRelease)
			if err != nil {
				return "", err
			}

			preReleaseInt++
		}

		newVersion = fmt.Sprintf("%s-%s.%d", newVersion, preReleaseLabel, preReleaseInt)
	}

	return newVersion, nil
}

func isRegexMatch(input, pattern string) bool {
	match, err := regexp.MatchString(pattern, input)
	if err != nil {
		return false
	}

	return match
}

func splitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
