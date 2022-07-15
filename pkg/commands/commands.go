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

type errNotEnoughArgs struct {
	command string
}

func (e errNotEnoughArgs) Error() string {
	return fmt.Sprintf("not enough args to command `%s`", e.command)
}

var errInvalidRepo = errors.New("invalid repo name, must have a '/'")

type repo struct {
	Owner string
	Name  string
}

func (r repo) String() string {
	return strings.Join([]string{r.Owner, r.Name}, "/")
}

func newRepo(fullRepo string) (repo, error) {
	split := strings.Split(fullRepo, "/")
	if len(split) != 2 {
		return repo{"", ""}, errInvalidRepo
	}

	o := split[0]
	r := split[1]

	return repo{o, r}, nil
}

// currentRepo gets the current repo the user is in
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
