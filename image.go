package image_uploader

import "time"

type Image struct {
	Hash      string    `json:"hash" gorm:"primary_key;type:char(32)"`
	Format    string    `json:"Format" gorm:"not null"`
	Size      int64     `json:"Size" gorm:"NOT NULL"`
	Title     string    `json:"title" gorm:"not null"`
	Width     uint      `json:"Width" gorm:"type:MEDIUMINT UNSIGNED;not null"`
	Height    uint      `json:"Height" gorm:"type:MEDIUMINT UNSIGNED;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
