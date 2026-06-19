package models

import (
"time"
"gorm.io/gorm"
)

type PasswordReset struct {
ID        uint           `gorm:"primaryKey" json:"id"`
Email     string         `gorm:"not null;index" json:"email"`
Token     string         `gorm:"uniqueIndex;not null" json:"token"`
ExpiresAt time.Time      `json:"expires_at"`
Used      bool           `gorm:"default:false" json:"used"`
CreatedAt time.Time      `json:"created_at"`
UpdatedAt time.Time      `json:"updated_at"`
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
