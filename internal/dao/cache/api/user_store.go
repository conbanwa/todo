package api

import "github.com/conbanwa/todo/internal/model"

// UserStore defines the interface for user storage operations
type UserStore interface {
	Create(user *model.User) (int64, error)
	GetByID(id int64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
}
