package reaper

import (
	"log"
	"time"

	"github.com/admpub/boltstore/shared"
	"github.com/boltdb/bolt"
)

//##############//
//### Public ###//
//##############//

// Run invokes a reap function as a goroutine.
func Run(db *bolt.DB, options Options) (chan<- struct{}, <-chan struct{}) {
	options.setDefault()
	quitC, doneC := make(chan struct{}), make(chan struct{})
	go reap(db, options, quitC, doneC)
	return quitC, doneC
}

// Quit terminates the reap goroutine.
func Quit(quitC chan<- struct{}, doneC <-chan struct{}) {
	quitC <- struct{}{}
	<-doneC
}

//###############//
//### Private ###//
//###############//

func reap(db *bolt.DB, options Options, quitC <-chan struct{}, doneC chan<- struct{}) {
	// Create a new ticker
	ticker := time.NewTicker(options.CheckInterval)

	defer func() {
		// Stop the ticker
		ticker.Stop()
	}()

	var prevKey []byte

	for {
		select {
		case <-quitC: // Check if a quit signal is sent.
			doneC <- struct{}{}
			return
		case <-ticker.C: // Check if the ticker fires a signal.
			// This slice is a buffer to save all expired session keys.
			expiredSessionKeys := make([][]byte, 0)

			// Start a bolt read transaction.
			err := db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket(options.BucketName)
				if bucket == nil {
					return nil
				}

				c := bucket.Cursor()

				var i int
				var isExpired bool

				for k, v := c.Seek(prevKey); ; k, v = c.Next() {
					// If we hit the end of our sessions then
					// exit and start over next time.
					if k == nil {
						prevKey = nil
						return nil
					}

					i++

					// The flag if the session is expired
					isExpired = false

					session, err := shared.Session(v)
					if err != nil {
						// Just remove the session with the invalid session data.
						// Log the error first.
						log.Printf("boltstore: removing session from database with invalid value: %v", err)
						isExpired = true
					} else if shared.Expired(session) {
						isExpired = true
					}

					if isExpired {
						// Copy the byte slice key, because this data is
						// not safe outside of this transaction.
						temp := make([]byte, len(k))
						copy(temp, k)

						// Add it to the expired sessios keys slice
						expiredSessionKeys = append(expiredSessionKeys, temp)
					}

					if options.BatchSize == i {
						// Store the current key to the previous key.
						// Copy the byte slice key, because this data is
						// not safe outside of this transaction.
						prevKey = make([]byte, len(k))
						copy(prevKey, k)
						return nil
					}
				}
			})

			if err != nil {
				log.Printf("boltstore: obtain expired sessions error: %v", err)
			}

			if len(expiredSessionKeys) > 0 {
				// Remove the expired sessions from the database
				err = db.Update(func(txu *bolt.Tx) error {
					// Get the bucket
					b := txu.Bucket(options.BucketName)
					if b == nil {
						return nil
					}

					// Remove all expired sessions in the slice
					for _, key := range expiredSessionKeys {
						err = b.Delete(key)
						if err != nil {
							return err
						}
					}

					return nil
				})

				if err != nil {
					log.Printf("boltstore: remove expired sessions error: %v", err)
				}
			}
		}
	}
}
