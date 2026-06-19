package services

import (
"crypto/rand"
"encoding/hex"
"errors"
"time"

"github.com/omorojames5-prog/glass-guitar/internal/models"
"github.com/omorojames5-prog/glass-guitar/pkg/database"
"gorm.io/gorm"
)

type PasswordResetService struct {
db *gorm.DB
}

func NewPasswordResetService() *PasswordResetService {
return &PasswordResetService{
db: database.DB,
}
}

func (s *PasswordResetService) CreateResetToken(email string) (string, error) {
// Check if user exists
var user models.User
if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
return "", errors.New("user not found")
}

// Generate random token
tokenBytes := make([]byte, 32)
if _, err := rand.Read(tokenBytes); err != nil {
return "", err
}
token := hex.EncodeToString(tokenBytes)

// Delete existing tokens for this email
s.db.Where("email = ? AND used = ?", email, false).Delete(&models.PasswordReset{})

// Create new reset token
reset := models.PasswordReset{
Email:     email,
Token:     token,
ExpiresAt: time.Now().Add(time.Hour * 24), // 24 hours expiration
Used:      false,
}

if err := s.db.Create(&reset).Error; err != nil {
return "", err
}

return token, nil
}

func (s *PasswordResetService) ValidateResetToken(token string) (string, error) {
var reset models.PasswordReset
if err := s.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&reset).Error; err != nil {
return "", errors.New("invalid or expired token")
}

return reset.Email, nil
}

func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
var reset models.PasswordReset
if err := s.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&reset).Error; err != nil {
return errors.New("invalid or expired token")
}

// Find user
var user models.User
if err := s.db.Where("email = ?", reset.Email).First(&user).Error; err != nil {
return errors.New("user not found")
}

// Update password (will be hashed by BeforeCreate hook or service)
user.Password = newPassword

// Mark token as used
reset.Used = true

// Start transaction
tx := s.db.Begin()
if err := tx.Save(&user).Error; err != nil {
tx.Rollback()
return err
}
if err := tx.Save(&reset).Error; err != nil {
tx.Rollback()
return err
}
tx.Commit()

return nil
}
