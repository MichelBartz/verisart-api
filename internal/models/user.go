package models

import (
	"fmt"
	"verisart-api/internal"
)

// UserModel a user model object
type UserModel struct {
	store      map[string]*User
	emailIndex map[string]*User
}

// User represents a Verisart user
type User struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Our "database", would be injected in an ideal world
var store = make(map[string]User)

// NewUserModel returns a new model to interact with
func NewUserModel() *UserModel {
	model := UserModel{
		store:      make(map[string]*User),
		emailIndex: make(map[string]*User),
	}

	return &model
}

// GetUserByID returns the user associated with the id
func (u *UserModel) GetUserByID(id string) (user *User, err error) {
	user, ok := u.store[id]

	if !ok {
		err = fmt.Errorf("User not found")
	}

	return
}

// GetUserByEmail returns the user matching the email
func (u *UserModel) GetUserByEmail(email string) (user *User, err error) {
	user, ok := u.emailIndex[email]

	if !ok {
		err = fmt.Errorf("No user with email %s found", email)
	}

	return
}

// GetUsers returns all existing users
func (u *UserModel) GetUsers() (users []*User) {
	for _, value := range u.store {
		users = append(users, value)
	}

	if users == nil {
		users = make([]*User, 0)
	}

	return
}

// CreateUser creates a user with the given email
func (u *UserModel) CreateUser(user *User) (err error) {
	id := internal.GenerateIDFromString(user.Email)

	if _, ok := u.store[id]; ok {
		err = fmt.Errorf("User already exists")
	} else {
		user.ID = id
		u.store[id] = user
		u.emailIndex[user.Email] = user
	}

	return
}
