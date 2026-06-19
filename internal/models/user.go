package models

import (
"time"
"gorm.io/gorm"
)

type User struct {
ID        uint           `gorm:"primaryKey" json:"id"`
Name      string         `gorm:"not null" json:"name"`
Email     string         `gorm:"uniqueIndex;not null" json:"email"`
Password  string         `gorm:"not null" json:"password"`
Bio       string         `gorm:"type:text" json:"bio,omitempty"`
Avatar    string         `gorm:"size:255" json:"avatar,omitempty"`
IsActive  bool           `gorm:"default:true" json:"is_active"`
CreatedAt time.Time      `json:"created_at"`
UpdatedAt time.Time      `json:"updated_at"`
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
