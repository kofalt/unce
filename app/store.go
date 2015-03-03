package app


import (
	. "fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

var (
	SeenFile = []string{ ".local", "share", "unce", "seen.db" }
	LogFile = []string{ ".local", "share", "unce", "log.db" }
)

// loltrashycodecomeatmebro
func GetSeenFilename() string {
	home, err := homedir.Dir()
	if err != nil { Println("Could not find homedir", err); os.Exit(1) }

	path := []string{}
	path = append(path, home)
	path = append(path, SeenFile...)

	return filepath.Join(path...)
}
func GetLogFilename() string {
	home, err := homedir.Dir()
	if err != nil { Println("Could not find homedir", err); os.Exit(1) }

	path := []string{}
	path = append(path, home)
	path = append(path, LogFile...)

	return filepath.Join(path...)
}

// Open some BEES
func Bees() (*bolt.DB, *bolt.DB) {
	os.Mkdir(filepath.Dir(GetSeenFilename()), 0777)
	os.Mkdir(filepath.Dir(GetLogFilename()), 0777)

	timeout := &bolt.Options{Timeout: 1 * time.Second}

	seenDB, err := bolt.Open(GetSeenFilename(), 0600, timeout)
	if err != nil { Println("Could not open ", GetSeenFilename(), err); os.Exit(1) }

	logDB, err := bolt.Open(GetLogFilename(), 0600, timeout)
	if err != nil { Println("Could not open ", GetLogFilename(), err); os.Exit(1) }

	return seenDB, logDB
}
