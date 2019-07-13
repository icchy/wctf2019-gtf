package main

import (
	"bytes"
	"encoding/gob"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/crypto/bcrypt"
)

type DBConn struct {
	DB *leveldb.DB
}

func (dbconn *DBConn) UserExists(username string) (bool, error) {
	ret, err := dbconn.DB.Has([]byte("user-"+username), nil)
	if err != nil {
		return false, err
	}
	return ret, nil
}

func (dbconn *DBConn) VerifyUser(username, password string) (bool, error) {
	snapshot, err := dbconn.DB.GetSnapshot()
	defer snapshot.Release()
	if err != nil {
		return false, err
	}

	data, err := snapshot.Get([]byte("user-"+username), nil)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(data, []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (dbconn *DBConn) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tx, err := dbconn.DB.OpenTransaction()
	if err != nil {
		return err
	}
	tx.Put([]byte("user-"+username), []byte(hash), nil)
	return tx.Commit()
}

func (dbconn *DBConn) GetGyotakuList(username string) ([]string, error) {
	ret, err := dbconn.DB.Has([]byte("gyotaku-"+username), nil)
	if err != nil {
		return nil, err
	}
	if !ret {
		return []string{}, nil
	}

	data, err := dbconn.DB.Get([]byte("gyotaku-"+username), nil)
	if err != nil {
		return nil, err
	}
	var gyotaku []string
	buf := bytes.NewBuffer(data)
	err = gob.NewDecoder(buf).Decode(&gyotaku)
	if err != nil {
		return nil, err
	}
	return gyotaku, nil
}

func (dbconn *DBConn) AddGyotakuList(username, gid string) error {
	tx, err := dbconn.DB.OpenTransaction()
	if err != nil {
		return err
	}

	ret, err := tx.Has([]byte("gyotaku-"+username), nil)
	if err != nil {
		return err
	}

	var gyotaku []string

	if ret {
		data, err := tx.Get([]byte("gyotaku-"+username), nil)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(data)
		err = gob.NewDecoder(buf).Decode(&gyotaku)
		if err != nil {
			return err
		}
	}

	gyotaku = append(gyotaku, gid)
	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(&gyotaku)
	if err != nil {
		return err
	}

	err = tx.Put([]byte("gyotaku-"+username), buf.Bytes(), nil)
	if err != nil {
		return err
	}

	return tx.Commit()
}
