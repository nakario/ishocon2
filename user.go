package main

import (
	"database/sql"
	"sync"
)

// User Model
type User struct {
	ID       int
	Name     string
	Address  string
	MyNumber string
	Votes    int
	Voted    int
	sync.Mutex
}

func getUser(name string, address string, myNumber string) (*User, error) {
	user, ok := usersMap[myNumber]
	if !ok || user.Name != name || user.Address != address {
		return nil, sql.ErrNoRows
	}
	return user, nil
}
