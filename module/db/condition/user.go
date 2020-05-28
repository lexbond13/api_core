package condition

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type IUser interface {
	Apply(q *orm.Query) *orm.Query
}

func NewUserCondition() *User {
	return &User{}
}

type User struct {
	id int64
	ids []int64
	email string
	activeKey string
}

func(u *User) SetID(id int64) {
	u.id = id
}

func(u *User) SetIDs(ids []int64) {
	u.ids = ids
}

func(u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetActiveKey(activeKey string) {
	u.activeKey = activeKey
}

func (u *User) Apply(q *orm.Query) *orm.Query {

	if u.id > 0 {
		q.Where("id = ?", u.id)
	}

	if len(u.ids) > 0 {
		q.Where("id IN (?)", pg.In(u.ids))
	}

	if u.email != "" {
		q.Where("email = ?", u.email)
	}

	if u.activeKey != "" {
		q.Where("active_key = ?", u.activeKey)
	}

	return q
}

