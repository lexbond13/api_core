package structure

import "time"

type File struct {
	//nolint
	tableName struct{} `json:"-" pg:"public.files"`

	ID           int64     `pg:"id"`
	UserID       int64     `pg:"user_id"`
	OriginalName string    `pg:"original_name"`
	Type         string    `pg:"type"`
	ExternalURL  string    `pg:"external_url"`
	Size         int64     `pg:"size"`
	CreatedAt    time.Time `pg:"created_at,notnull,use_zero"`
	UpdatedAt    time.Time `pg:"updated_at,notnull,use_zero"`
}
