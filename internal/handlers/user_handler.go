package handlers

import (
"fmt"
"log"
"net/http"
"strconv"
"strings"

"github.com/gin-gonic/gin"
"github.com/omorojames5-prog/glass-guitar/internal/models"
"github.com/omorojames5-prog/glass-guitar/internal/services"
"github.com/omorojames5-prog/glass-guitar/internal/validation"
"gorm.io/gorm"
)

type UserHandler struct {
service *services.UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
return &UserHandler{
service: services.NewUserService(db),
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

func (h *UserHandler) GetUsers(c *gin.Context) {
users, err := h.service.GetAllUsers()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
return
}

user, err := h.service.GetUserByID(uint(id))
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}
c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
log.Println("🔵 [CreateUser] =====================")
log.Println("🔵 [CreateUser] Endpoint called")

var user models.User
if err := c.ShouldBindJSON(&user); err != nil {
log.Printf("🔴 [CreateUser] Failed to bind JSON: %v", err)
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

// Validate input
if !validation.IsValidEmail(user.Email) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
return
}

if !validation.IsValidPassword(user.Password) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters and contain letters and numbers"})
return
}

if !validation.IsValidName(user.Name) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Name must be at least 2 characters"})
return
}

log.Printf("🔵 [CreateUser] Email: '%s'", user.Email)
log.Printf("🔵 [CreateUser] Password: '%s'", user.Password)
log.Printf("🔵 [CreateUser] Password hex: %s", hexDump(user.Password))
log.Printf("🔵 [CreateUser] Password length: %d", len(user.Password))

if err := h.service.CreateUser(&user); err != nil {
log.Printf("🔴 [CreateUser] Failed to create user: %v", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

log.Printf("🟢 [CreateUser] User created successfully with ID: %d", user.ID)
c.JSON(http.StatusCreated, gin.H{
"message": "User created successfully",
"user":    user,
})
}

func (h *UserHandler) Login(c *gin.Context) {
log.Println("🔵 [Login] =====================")
log.Println("🔵 [Login] Endpoint called")

var loginRequest struct {
Email    string `json:"email" binding:"required"`
Password string `json:"password" binding:"required"`
}

if err := c.ShouldBindJSON(&loginRequest); err != nil {
log.Printf("🔴 [Login] Failed to bind JSON: %v", err)
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

email := strings.TrimSpace(loginRequest.Email)
password := strings.TrimSpace(loginRequest.Password)

log.Printf("🔵 [Login] Raw password hex: %x", []byte(password))
log.Printf("🔵 [Login] Email: '%s'", email)
log.Printf("🔵 [Login] Password: '%s'", password)
log.Printf("🔵 [Login] Password hex: %s", hexDump(password))
log.Printf("🔵 [Login] Password length: %d", len(password))
log.Printf("🔵 [Login] About to authenticate with password: '%s'", password)

user, err := h.service.AuthenticateUser(email, password)
if err != nil {
log.Printf("🔴 [Login] Authentication failed: %v", err)
c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
return
}

token, err := h.service.GenerateJWT(user)
if err != nil {
log.Printf("🔴 [Login] Failed to generate token: %v", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
return
}

log.Printf("🟢 [Login] Login successful for user: %s", user.Email)
c.JSON(http.StatusOK, gin.H{
"token": token,
"user":  user,
})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
return
}

var user models.User
if err := c.ShouldBindJSON(&user); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
user.ID = uint(id)

if err := h.service.UpdateUser(&user); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, gin.H{
"message": "User updated successfully",
"user":    user,
})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
return
}

if err := h.service.DeleteUser(uint(id)); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
var request struct {
Email string `json:"email" binding:"required"`
}

if err := c.ShouldBindJSON(&request); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

if !validation.IsValidEmail(request.Email) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
return
}

c.JSON(http.StatusOK, gin.H{
"message": "If the email exists, a reset link has been sent",
})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
var request struct {
Token       string `json:"token" binding:"required"`
NewPassword string `json:"new_password" binding:"required"`
}

if err := c.ShouldBindJSON(&request); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

if !validation.IsValidPassword(request.NewPassword) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters and contain letters and numbers"})
return
}

c.JSON(http.StatusOK, gin.H{
"message": "Password reset successful",
})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
userID, exists := c.Get("user_id")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
return
}

user, err := h.service.GetUserByID(uint(userID.(float64)))
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}

c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
userID, exists := c.Get("user_id")
if !exists {
c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
return
}

var updateData struct {
Name   string `json:"name"`
Bio    string `json:"bio"`
Avatar string `json:"avatar"`
}

if err := c.ShouldBindJSON(&updateData); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

user, err := h.service.GetUserByID(uint(userID.(float64)))
if err != nil {
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
return
}

if updateData.Name != "" {
if !validation.IsValidName(updateData.Name) {
c.JSON(http.StatusBadRequest, gin.H{"error": "Name must be at least 2 characters"})
return
}
user.Name = updateData.Name
}
if updateData.Bio != "" {
user.Bio = updateData.Bio
}
if updateData.Avatar != "" {
user.Avatar = updateData.Avatar
}

if err := h.service.UpdateUser(user); err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{
"message": "Profile updated successfully",
"user":    user,
})
}
