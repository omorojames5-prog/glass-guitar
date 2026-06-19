package services

import (
"errors"
"fmt"
"log"
"os"
"time"

"github.com/omorojames5-prog/glass-guitar/internal/models"
"github.com/omorojames5-prog/glass-guitar/internal/repositories"
"github.com/golang-jwt/jwt/v5"
"golang.org/x/crypto/bcrypt"
"gorm.io/gorm"
)

type UserService struct {
repo *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
return &UserService{
repo: repositories.NewUserRepository(db),
}
}

func hexDump(data string) string {
if data == "" {
return "empty"
}
result := ""
for i, b := range []byte(data) {
if i > 0 && i%4 == 0 {
result += " "
}
result += fmt.Sprintf("%02x", b)
}
return result
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(user *models.User) error {
log.Printf("🔵 [CreateUser] =====================")
log.Printf("🔵 [CreateUser] Creating user with email: %s", user.Email)
log.Printf("🔵 [CreateUser] Original password: '%s'", user.Password)
log.Printf("🔵 [CreateUser] Password hex: %s", hexDump(user.Password))
log.Printf("🔵 [CreateUser] Password length: %d", len(user.Password))

// Hash the password
if user.Password == "" {
log.Printf("🔴 [CreateUser] Password is empty!")
return errors.New("password cannot be empty")
}

hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
if err != nil {
log.Printf("🔴 [CreateUser] Error hashing password: %v", err)
return err
}
user.Password = string(hashedPassword)

log.Printf("🟢 [CreateUser] Password hashed successfully")
log.Printf("🟢 [CreateUser] Hashed password: '%s'", user.Password)
log.Printf("🟢 [CreateUser] Hashed hex: %s", hexDump(user.Password))
log.Printf("🟢 [CreateUser] Hashed length: %d", len(user.Password))

return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
if user.Password != "" {
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
if err != nil {
return err
}
user.Password = string(hashedPassword)
}
return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
return s.repo.Delete(id)
}

func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
log.Printf("🔵 [AuthenticateUser] =====================")
log.Printf("🔵 [AuthenticateUser] Login attempt for email: %s", email)
log.Printf("🔵 [AuthenticateUser] Input password: '%s'", password)
log.Printf("🔵 [AuthenticateUser] Password hex: %s", hexDump(password))
log.Printf("🔵 [AuthenticateUser] Password length: %d", len(password))

users, err := s.repo.GetAll()
if err != nil {
log.Printf("🔴 [AuthenticateUser] Error getting users: %v", err)
return nil, err
}

log.Printf("🟡 [AuthenticateUser] Found %d users in database", len(users))

for _, user := range users {
if user.Email == email {
log.Printf("🟢 [AuthenticateUser] Found matching user: %s", email)
log.Printf("🟢 [AuthenticateUser] Stored password: '%s'", user.Password)
log.Printf("🟢 [AuthenticateUser] Stored hex: %s", hexDump(user.Password))
log.Printf("🟢 [AuthenticateUser] Stored length: %d", len(user.Password))

// Try bcrypt comparison
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
if err == nil {
log.Printf("🟢 [AuthenticateUser] ✅ Password matched!")
return &user, nil
} else {
log.Printf("🔴 [AuthenticateUser] ❌ Password did not match")
log.Printf("🔴 [AuthenticateUser] bcrypt error: %v", err)
return nil, errors.New("invalid credentials")
}
}
}

log.Printf("🔴 [AuthenticateUser] No user found with email: %s", email)
return nil, errors.New("invalid credentials")
}

func (s *UserService) GenerateJWT(user *models.User) (string, error) {
secret := os.Getenv("JWT_SECRET")
if secret == "" {
secret = "your-secret-key"
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
"user_id": user.ID,
"email":   user.Email,
"exp":     time.Now().Add(time.Hour * 24).Unix(),
})

return token.SignedString([]byte(secret))
}
