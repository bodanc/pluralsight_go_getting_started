// the model layer typically exposes some sort of behavior to allow us to manipulate the application state;

// probably the most fundamental OOP concept or OOP construct is the ability to add behavior and associate that to a
// custom data type (object/struct);
// in other words, we need to be able to create methods;

package models

import (
	"errors"
	"fmt"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
}

var (
	// a slice that holds pointers to User objects;
	users  []*User
	nextID = 1
)

// return a slice of pointers to User objects;
func GetUsers() []*User {
	return users
}

// we accept a User coming in, we assign an Id to it, append it to our user collection and return that collection back;
func AddUser(u User) (User, error) {
	if u.ID != 0 {
		return User{}, errors.New("new User must not include an ID or it must be set to 0")
	}
	u.ID = nextID
	nextID++
	users = append(users, &u)
	return u, nil
}

func GetUserByID(id int) (User, error) {
	for _, u := range users {
		if u.ID == id {
			// we need to dereference the pointer because we're expecting a user value to come out of this function;
			return *u, nil
		}
	}
	return User{}, fmt.Errorf("user with ID '%v' not found", id)
}

func UpdateUser(u User) (User, error) {
	for i, candidate := range users {
		if candidate.ID == u.ID {
			users[i] = &u
			return u, nil
		}
	}
	return User{}, fmt.Errorf("user with ID '%v' not found", u.ID)
}

func RemoveUserById(id int) error {
	for i, u := range users {
		if u.ID == id {
			// we are performing a splice operation on the slice;
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user with ID '%v' not found", id)
}
