package issues

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/dzfrias/ghi/pkg/postutil"
)

// Modified during testing
var OpenUrl = "https://api.github.com/repos"
var ErrBadTempFile = errors.New("issue with using temporary file")

type issueFile struct {
	Title string
	Cont  string
}

// Open opens in issue in GitHub after prompting the user to create one
func Open(owner, repo string) error {
	f, err := os.CreateTemp(".", "ghi_temp")
	if err != nil {
		panic("error creating temporary file")
	}
	fname := f.Name()
	defer os.Remove(fname)

	err = openEditor(fname)
	if err != nil {
		return ErrBadTempFile
	}
	iss, err := readIssF(fname)
	if err != nil {
		return ErrBadTempFile
	}

	reqUrl := strings.Join([]string{OpenUrl, owner, repo, "issues"}, "/")
	req, err := openReq(reqUrl, iss)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("open issue failed with status: %s", resp.Status)
	}

	return nil
}

func openReq(reqUrl string, iss issueFile) (*http.Request, error) {
	vals := map[string]string{
		"title": iss.Title,
		"body":  iss.Cont,
	}

	buf, err := postutil.EncodeMap(vals)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	s, err := GetCreds(ConfigPath)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+s)

	return req, nil
}

// openEditor opens the given file in the default editor
func openEditor(fname string) error {
	e := os.Getenv("EDITOR")
	if e == "" {
		e = "vim"
	}
	cmd := exec.Command(e, fname)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// readIssF reads the file in the issue file format: TITLE\nCONTENTS
func readIssF(fname string) (issueFile, error) {
	var iss issueFile
	l, err := readLines(fname)
	if err != nil {
		return iss, err
	}

	if len(l) == 0 {
		return iss, ErrBadTempFile
	}
	iss.Title = l[0]
	if len(l) > 1 {
		iss.Cont = strings.Join(l[1:], "\n")
	} else {
		iss.Cont = ""
	}
	return iss, nil
}

// readLines reads a whole file into memory and returns a slice of its lines.
func readLines(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
