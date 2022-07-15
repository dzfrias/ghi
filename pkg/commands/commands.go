// Package commands holds the underlying functions called when the user uses a
// command.
// Ex:
//    `ghi list`
//    Calls commands.List
package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// errNotEnoughArgs is an error that is associated with a command, specifically
// thrown when the user does not provide the appriopriate amount of positional
// arguments to a command.
type errNotEnoughArgs struct {
	command string
}

// Error() returns the error in a specific format
func (e errNotEnoughArgs) Error() string {
	return fmt.Sprintf("not enough args to command `%s`", e.command)
}

var errInvalidRepo = errors.New("invalid repo name, must have a '/'")

// repo represents a GitHub repository
type repo struct {
	Owner string
	Name  string
}

// String joins the repo by a "/"
func (r repo) String() string {
	return strings.Join([]string{r.Owner, r.Name}, "/")
}

// newRepo creates a repo struct from a string, splitting by the "/". If an
// invalid repository is given, errInvalidRepo is returned.
func newRepo(fullRepo string) (repo, error) {
	split := strings.Split(fullRepo, "/")
	if len(split) != 2 {
		return repo{"", ""}, errInvalidRepo
	}

	o := split[0]
	r := split[1]

	return repo{o, r}, nil
}

// currentRepo gets the current repo the user is in in the form of a repo
func currentRepo() repo {
	var r repo

	pwd, err := os.Getwd()
	if err != nil {
		return r
	}

	_, err = exec.Command(
		"sh", "-c", fmt.Sprintf("cd %s; git status", pwd),
	).Output()
	if err != nil {
		return r
	}

	// Get origin url in cwd
	ori, err := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("cd %s; git config --get remote.origin.url", pwd),
	).Output()
	if err != nil {
		return r
	}

	repoStr := strings.TrimPrefix(string(ori), "https://github.com/")
	// Strip the '.git' at the end
	repoStr = repoStr[:len(repoStr)-5]

	r, err = newRepo(repoStr)
	if err != nil {
		return r
	}

	return r
}
