package fsrepo

import (
	"os"
	"runtime"
	"github.com/mitchellh/go-homedir"
)

// BestKnownPath returns the best known fsrepo path. If the ENV override is
// present, this function returns that value. Otherwise, it returns the default
// repo path.
func BestKnownPath() (string, error) {
	if runtime.GOOS == "darwin" { //TODO: Leave only ./btfs path by defining properly the $HOME dir for iOS in path.go file
		btfsPath := "~/Documents/.btfs" //iOS path
	}
	if runtime.GOOS == "android" {
		btfsPath := "~/.btfs" //Android path
	}
	if os.Getenv("BTFS_PATH") != "" {
		btfsPath = os.Getenv("BTFS_PATH")
	}
	curPath, err := homedir.Expand(btfsPath)
	if err != nil {
		return "", err
	}
	return curPath, nil
}
