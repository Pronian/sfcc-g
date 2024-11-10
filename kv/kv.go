package kv

import (
	"fmt"
	"sfcc/g/util"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB
var bucketName = []byte("main")
var expiresSuffix = "_expires"
var timeFormat = time.RFC3339

func Init(constantPath bool) {
	var err error
	var path string
	if constantPath {
		path = util.GetFilePathInExecutableDirectory("kv.db")
	} else {
		path = "./kv.db"
	}
	db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		return nil
	})
}

func Set(key, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put([]byte(key), []byte(value))
		if err != nil {
			fmt.Errorf("Error setting key \"%s\" to value \"%s\"", key, value)
		}
		return nil
	})
}

func SetTemporary(key, value string, duration time.Duration) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put([]byte(key), []byte(value))
		if err != nil {
			fmt.Errorf("Error setting temp key \"%s\" to value \"%s\"", key, value)
		}
		expDate := time.Now().Add(duration).Format(timeFormat)
		fmt.Printf("Key \"%s\" will expire at %s\n", key, expDate)
		err = b.Put([]byte(key+expiresSuffix), []byte(expDate))
		if err != nil {
			fmt.Errorf("Error setting temp key \"%s\" to value \"%s\"", key, value)
		}
		return nil
	})
}

func Get(key string) string {
	var value string
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		value = string(b.Get([]byte(key)))
		return nil
	})
	return value
}

func GetTemporary(key string) string {
	var value string
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		expires := string(b.Get([]byte(key + expiresSuffix)))
		if expires == "" {
			return nil
		}

		expireTime, err := time.Parse(timeFormat, expires)
		if err != nil {
			b.Delete([]byte(key))
			b.Delete([]byte(key + expiresSuffix))
			return nil
		}

		if time.Now().After(expireTime) {
			b.Delete([]byte(key))
			b.Delete([]byte(key + expiresSuffix))
			return nil
		}

		value = string(b.Get([]byte(key)))
		return nil
	})
	return value
}

func ClearExpired() {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		var deleted int

		for k, v := c.First(); k != nil; k, v = c.Next() {
			key := string(k)
			if strings.HasSuffix(key, expiresSuffix) == false {
				continue
			}

			expireTime, err := time.Parse(timeFormat, string(v))
			if err != nil {
				b.Delete([]byte(k))
				b.Delete([]byte(key + expiresSuffix))
				deleted++
				continue
			}

			if time.Now().After(expireTime) {
				b.Delete([]byte(k))
				b.Delete([]byte(key + expiresSuffix))
				deleted++
			}
		}

		if deleted > 0 {
			fmt.Printf("Deleted %d expired keys\n", deleted)
		} else {
			fmt.Println("No expired keys found")
		}
		return nil
	})
}

func Close() {
	db.Close()
}
