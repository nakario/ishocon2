package main

import "database/sql"

// User Model
type User struct {
	ID       int
	Name     string
	Address  string
	MyNumber string
	Votes    int
}

func getUser(name string, address string, myNumber string) (user User, err error) {
	user, ok := usersByMyNumber[myNumber]
	if !ok || user.Name != name || user.Address != address {
		return User{}, sql.ErrNoRows
	}
	return user, nil
}
