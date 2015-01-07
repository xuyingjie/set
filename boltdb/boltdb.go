package boltdb

import "github.com/boltdb/bolt"

type Blog struct {
	ID      string
	Title   string
	Content string
}

var db *bolt.DB

func Open(path string) {
	db, _ = bolt.Open(path, 0666, nil)
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("set"))
		return nil
	})
}

func Query(i int) []Blog {

	var blog []Blog

	db.View(func(tx *bolt.Tx) error {

		set := tx.Bucket([]byte("set"))
		set.ForEach(func(k, v []byte) error {

			var b Blog
			b.ID = string(k)
			b.Title = string(v)
			b.Content = string(tx.Bucket(k).Get([]byte("content")))

			blog = append(blog, b)

			return nil
		})
		return nil
	})

	var reverseBlog []Blog
	j := 0
	k := len(blog)
	for j < k && j < i {
		reverseBlog = append(reverseBlog, blog[k-1-j])
		j++
	}

	return reverseBlog
}

func Put(blog Blog) {

	db.Update(func(tx *bolt.Tx) error {

		// set
		set := tx.Bucket([]byte("set"))
		set.Put([]byte(blog.ID), []byte(blog.Title))

		// blog
		tx.CreateBucket([]byte(blog.ID))
		bucket := tx.Bucket([]byte(blog.ID))

		bucket.Put([]byte("title"), []byte(blog.Title))
		bucket.Put([]byte("content"), []byte(blog.Content))

		return nil
	})
}

func Get(id []byte) Blog {

	var b Blog

	db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(id)
		t := bucket.Get([]byte("title"))
		c := bucket.Get([]byte("content"))

		b.ID = string(id)
		b.Title = string(t)
		b.Content = string(c)

		return nil
	})
	return b
}

func Delete(id []byte) {

	db.Update(func(tx *bolt.Tx) error {

		tx.DeleteBucket(id)
		tx.Bucket([]byte("set")).Delete(id)
		return nil
	})
}
