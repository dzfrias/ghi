package issues

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const credsFile = "./TestCreds.txt"

func TestGetCreds(t *testing.T) {
	creds, err := GetCreds("./testdata/TestCreds.txt")
	assert.Equal(t, creds, "testing", "does not read creds correctly")
	assert.Nil(t, err, "throws an error")
}

func TestGetCredsNoneExist(t *testing.T) {
	const target = "no credentials. Run `ghi auth` to access this feature"
	_, err := GetCreds(credsFile)
	assert.Equal(t, err.Error(), target, "no credsFile does not throw error")
}

func TestStoreCredsNoConfigDir(t *testing.T) {
	const newCreds = "./.config/ghi/TestCreds.txt"
	err := StoreCreds("testing", newCreds)
	assert.Nil(t, err, "StoreCreds with no config dir throws error")
	assert.Equal(t, readfile(newCreds), "testing\n")
	os.Remove(newCreds)
}

func TestStoreCreds(t *testing.T) {
	err := StoreCreds("testing", credsFile)
	assert.Nil(t, err, "StoreCreds throws error")
	assert.Equal(t, readfile(credsFile), "testing\n")
	os.Remove(credsFile)
}

func readfile(fname string) string {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
