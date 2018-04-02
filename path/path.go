package path

import (
	"errors"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// Expand expands the tilde and %tmp abbreviations
func Expand(w string) (string, error) {
	// cant use root or v small paths
	if len(w) < 2 {
		return "", errors.New("path too short")
	}

	b := strings.Split(w, "/")
	r := ""
	if b[0] == "" {
		r = string(filepath.Separator)
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	hd := usr.HomeDir

	// expand ~
	if b[0] == "~" {
		if b[1] == "" { // disallow "~/"
			return "", errors.New("root of user folder not allowed")
		}
		b[0] = hd
	}
	// replace %tmp with a temp folder
	if b[0] == "%tmp" {
		tmp, err := ioutil.TempDir("", "floe")
		if err != nil {
			return "", err
		}
		b[0] = tmp
	}

	return r + filepath.Join(b...), nil
}