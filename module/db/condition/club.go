package condition

import "github.com/go-pg/pg/v9/orm"

type IClub interface {
	Apply(q *orm.Query) *orm.Query
}

type Club struct {
	id int64
}

// NewClubCondition
func NewClubCondition() IClub {
	return &Club{}
}

// SetId
func (c *Club) SetId(id int64) {
	c.id = id
}

// Apply
func (c *Club) Apply(q *orm.Query) *orm.Query {

	if c.id > 0 {
		q.Where("id = ?", c.id)
	}

	return q
}
