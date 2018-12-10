package models

import "testing"

var user = &User{
	Email: "test.user@verisart.com",
	Name:  "Test User",
}

func TestCreateUser(t *testing.T) {
	model := NewUserModel()

	model.CreateUser(user)
	if user.ID == "" {
		t.Error("Expected a user ID got ", user.ID)
	}
}

func TestDuplicateUser(t *testing.T) {
	model := NewUserModel()

	model.CreateUser(user)
	err := model.CreateUser(user)
	if err == nil {
		t.Error("Expected an error when creating duplicate users, got none.")
	}
}

func TestGetUser(t *testing.T) {
	model := NewUserModel()

	model.CreateUser(user)

	_, err := model.GetUserByID(user.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestGetUsers(t *testing.T) {
	model := NewUserModel()

	users := model.GetUsers()
	if len(users) != 0 {
		t.Errorf("Expected 0 user, got %d", len(users))
	}

	model.CreateUser(user)
	users = model.GetUsers()
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

func TestNonExistingUser(t *testing.T) {
	model := NewUserModel()
	_, err := model.GetUserByID("abcdef")
	if err == nil {
		t.Error("Expected an error when fetching a non-existing user")
	}
}
