package structure

import "time"

const ClubStatusNew  = "new"
const ClubStatusDeleted  = "deleted"

type Club struct {
	//nolint
	tableName struct{} `json:"-" pg:"public.club"`

	ID          int64     `json:"id" pg:"id,pk"`
	Name        string    `json:"name" pg:"name"`
	Tagline     string    `json:"tagline" pg:"tagline"`
	Logo        string    `json:"name" pg:"logo"`
	CoverImage  string    `json:"cover_image" pg:"cover_image"`
	Description string    `json:"description" pg:"description"`
	Address     string    `json:"address" pg:"address"`
	Email       string    `json:"email" pg:"email"`
	Phone       string    `json:"phone" pg:"phone"`
	Rating      int64     `json:"rating" pg:"rating"`
	Status      string    `json:"status" pg:"status"`
	OwnerID     int64     `json:"owner_id" pg:"owner_id"`
	CreatedAt   time.Time `json:"created_at" pg:"created_at,notnull,use_zero"`
	UpdatedAt   time.Time `json:"updated_at" pg:"updated_at,notnull,use_zero"`
}
