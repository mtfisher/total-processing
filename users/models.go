package users

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/mtfisher/total-processing/core"
	"golang.org/x/crypto/bcrypt"
)

const UserBucket = "Users"

type User struct {
	Username     string
	PasswordHash string
}

type UserRepository struct {
	c core.Core
}

func (u *User) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *User) PasswordMatches(password string) error {
	if len(password) == 0 {
		return errors.New("The password can't be empty")
	}
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)

	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (ur UserRepository) registerNewUser(username, password string) (*User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !ur.isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := User{Username: username}
	u.setPassword(password)

	encoded, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	ur.c.InsertDB([]byte(UserBucket), []byte(u.Username), []byte(encoded))

	return &u, nil
}

func (ur UserRepository) getUser(username string) (*User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, errors.New("Details could not be found")
	}

	val := ur.c.QueryDB([]byte(UserBucket), []byte(username))
	if len(val) == 0 {
		return nil, errors.New("Details could not be found")
	}

	user := User{}
	err := json.Unmarshal(val, &user)
	if err != nil {
		return nil, errors.New("Details could not be found")
	}

	return &user, nil
}

func newUserRepository(c core.Core) UserRepository {
	return UserRepository{c: c}
}

// Check if the supplied username is available
func (ur UserRepository) isUsernameAvailable(username string) bool {
	val := ur.c.QueryDB([]byte(UserBucket), []byte(username))

	return len(val) == 0
}
