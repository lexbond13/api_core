package structure

import "time"

const StatusEnabled = "enabled"
const StatusDisabled = "disabled"
const StatusWaitActivate = "wait_activate"
const StatusNew = "new"

type User struct {
	//nolint
	tableName struct{} `json:"-" pg:"public.user"`

	ID        int64     `json:"id" pg:"id,pk"`
	Name      string    `json:"name" pg:"name"`
	Email     string    `json:"email" pg:"email,unique"`
	ActiveKey string    `json:"active_key" pg:"active_key"`
	Status    string    `json:"status" pg:"status"`
	CreatedAt time.Time `pg:"created_at,notnull,use_zero"`
	UpdatedAt time.Time `pg:"updated_at,notnull,use_zero"`
}
