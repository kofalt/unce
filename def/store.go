package def


import (
	"encoding/json"
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

	// zero byte to store in db
	SeenPlaceholder = []byte{byte(0)}
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

	seen, err := bolt.Open(GetSeenFilename(), 0600, timeout)
	if err != nil { Println("Could not open ", GetSeenFilename(), err); os.Exit(1) }

	log, err := bolt.Open(GetLogFilename(), 0600, timeout)
	if err != nil { Println("Could not open ", GetLogFilename(), err); os.Exit(1) }

	err = seen.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("github"))
		return nil
	})
	if err !=nil { Println("Error updating seen.db:", err); os.Exit(1) }

	err = log.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("github"))
		return nil
	})
	if err !=nil { Println("Error updating log.db:", err); os.Exit(1) }


	return seen, log
}

func IsSeen(db *bolt.DB, bucket, key string) bool {
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(key))
		return nil
	})
	if err !=nil { Println("Error checking bees:", err); os.Exit(1) }

	return v != nil
}

func MarkSeen(db *bolt.DB, bucket, key string) {
	Println("Marking seen", bucket, key)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.Put([]byte(key), SeenPlaceholder)
		return nil
	})
	if err !=nil { Println("Error updating bees:", err); os.Exit(1) }
}

func Store(db *bolt.DB, bucket, key string, val []byte) {
	Println("Storing val", bucket, key)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.Put([]byte(key), val)
		return nil
	})
	if err !=nil { Println("Error updating bees:", err); os.Exit(1) }
}

func StoreJSON(db *bolt.DB, bucket, key string, val interface{}) {
	b, err := json.Marshal(val)
	if err != nil { Println("Error json:", err); os.Exit(1) }

	Store(db, bucket, key, b)
}
