package main

import (
	"encoding/json"
	"os"
)

type User struct {
	ID   string
	Name string
}

func GetUsers() ([]User, error) {
	file, err := os.Open("data.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
