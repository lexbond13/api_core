package db

import (
	"github.com/lexbond13/api_core/module/db/condition"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/go-pg/pg/v9"
	"time"
)

type IUserRepository interface {
	Create (user *structure.User) error
	Update(user *structure.User) error
	FindByID(ID int64) (*structure.User, error)
	FindOne(condition condition.IUser) (*structure.User, error)
	FindByCondition(condition condition.IUser) ([]*structure.User, error)
}

// NewUserRepository
func NewUserRepository() IUserRepository {
	return &userRepository{}
}

type userRepository struct {
}

// Create
func (ur *userRepository) Create(user *structure.User) error {
	user.CreatedAt = time.Now()
	_, err := connection.Model(user).Insert()
	if err != nil {
		return err
	}

	return nil
}

//Update
func (ur *userRepository) Update(user *structure.User) error {
	user.UpdatedAt = time.Now()
	_, err := connection.Model(user).WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

// FindByID
func (ur *userRepository) FindByID(ID int64) (*structure.User, error) {
	user := &structure.User{ID: ID}
	err := connection.Model(user).WherePK().First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// FindOne
func (ur *userRepository) FindOne(condition condition.IUser) (*structure.User, error) {

	result := &structure.User{}
	query := connection.Model(result)
	condition.Apply(query)

	err := query.First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

// FindByCondition
func (ur *userRepository) FindByCondition(condition condition.IUser) (result []*structure.User, err error) {

	query := connection.Model(&result)
	condition.Apply(query)

	err = query.Select()
	if err != nil {
		return nil, err
	}

	return
}
