package db

import (
	"github.com/lexbond13/api_core/module/db/condition"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/go-pg/pg/v9"
	"time"
)

type IClubRepository interface {
	FindByID(ID int64) (*structure.Club, error)
	FindByCondition(condition *condition.Club) ([]*structure.Club, error)
	Create(club *structure.Club) error
	Update(club *structure.Club) error
	Delete(ID int64) error
}

type clubRepository struct {
}

func NewClubRepository() IClubRepository {
	return &clubRepository{}
}

// FindByID
func (c *clubRepository) FindByID(ID int64) (*structure.Club, error) {
	club := &structure.Club{ID: ID}
	err := connection.Model(club).WherePK().First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return club, nil
}

// FindByCondition
func (c *clubRepository) FindByCondition(condition *condition.Club) (result []*structure.Club, err error) {
	query := connection.Model(&result)
	condition.Apply(query)

	err = query.Select()
	if err != nil {
		return nil, err
	}

	return
}

// Create
func (c *clubRepository) Create(club *structure.Club) error {
	club.CreatedAt = time.Now()
	_, err := connection.Model(club).Insert()
	if err != nil {
		return err
	}

	return nil
}

// Update
func (c *clubRepository) Update(club *structure.Club) error {
	club.UpdatedAt = time.Now()
	_, err := connection.Model(club).WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

// Delete
func (c *clubRepository) Delete(ID int64) error {
	club := &structure.Club{ID: ID}
	_, err := connection.Model(club).WherePK().Delete()
	if err != nil {
		return err
	}

	return nil
}

