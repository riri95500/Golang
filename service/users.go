package users

import "fmt"

type User struct {
	ID       string
	Username string
}

type UserRepository interface {
	AddUser(user User) error
	RemoveUser(userID string) error
	GetUserByID(userID string) (User, error)
	GetUserByUsername(username string) (User, error)
	GetAllUsers() ([]User, error)
}

type userRepository struct {
	users map[string]User
}

func (repo *userRepository) AddUser(user User) error {
	repo.users[user.ID] = user
	return nil
}

func (repo *userRepository) RemoveUser(userID string) error {
	delete(repo.users, userID)
	return nil
}

func (repo *userRepository) GetUserByID(userID string) (User, error) {
	user, ok := repo.users[userID]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (repo *userRepository) GetUserByUsername(username string) (User, error) {
	for _, user := range repo.users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}

func (repo *userRepository) GetAllUsers() ([]User, error) {
	users := make([]User, 0, len(repo.users))

	for _, user := range repo.users {
		users = append(users, user)
	}

	return users, nil
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make(map[string]User),
	}
}
