package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Transact runs a function encapsulated by a single atomic transaction.
// If any error is returned by tf, or if tf panics, the database transaction will be rolled back.
// Otherwise the transaction will be committed.
//
// Source: https://stackoverflow.com/a/23502629/1320648
//
// Usage:
// 	Transact(db, func(tx *gorm.DB) error {
// 		err := tx.
// 			Model(&User{ID: 123}).
// 			Update(&User{ID: 123}).
// 			Error
// 		if err != nil {
// 			return fmt.Errorf("could not update user. %s", err)
// 		}
// 	})
func Transact(db *gorm.DB, tf func(tx *gorm.DB) error) error {
	var err error
	if commonDB, ok := db.CommonDB().(sqlTx); ok && commonDB != nil {
		// If the db is already in a transaction, just execute tf
		// and let the outer transaction handle Rollback and Commit.
		err = tf(db)
		return err
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("could not start transaction. %w", err)
	}
	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	err = tf(tx)
	return err
}

// sqlTx is a helper interface to check if a gorm.DB.CommonDB() is already in a transaction.
type sqlTx interface {
	Commit() error
	Rollback() error
}
