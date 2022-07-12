package issues

import (
	"os"
	"testing"
)

const credsFile = "./TestCreds.txt"

func setupCreds(t *testing.T) {
	f, err := os.Create(credsFile)
	if err != nil {
		t.Error("setupCreds() failed")
	}
	defer f.Close()
	_, err = f.WriteString("testing")
	if err != nil {
		t.Error("setupCreds() failed")
	}
}

func TestGetCreds(t *testing.T) {
	setupCreds(t)
	creds, err := GetCreds(credsFile)
	if creds != "testing" {
		t.Error(`GetCreds(credsFile) != "testing"`)
	}
	if err != nil {
		t.Error(`GetCreds(credsFile) returns error`)
	}
	os.Remove(credsFile)
}

func TestGetCredsNoneExist(t *testing.T) {
	const target = "no credentials. Run `ghi auth` to access this feature"
	_, err := GetCreds(credsFile)
	if err.Error() != target {
		t.Errorf(`GetCreds(credsFile) with no credsFile does not throw error: %q != %q`, err, target)
	}
}

func TestStoreCredsNoConfigDir(t *testing.T) {
	const newCreds = "./.config/ghi/TestCreds.txt"
	err := StoreCreds("testing", newCreds)
	if err != nil {
		t.Errorf(`StoreCreds("testing", %s) returns %v`, newCreds, err)
	}
	os.Remove(newCreds)
}

func TestStoreCreds(t *testing.T) {
	err := StoreCreds("testing", credsFile)
	if err != nil {
		t.Errorf(`StoreCreds("testing", credsFile) returns %v`, err)
	}
	os.Remove(credsFile)
}
