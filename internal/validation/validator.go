package validation

import (
"regexp"
"strings"
)

func IsValidEmail(email string) bool {
// Basic email validation
if strings.TrimSpace(email) == "" {
return false
}

// Email regex pattern
regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
return regex.MatchString(email)
}

func IsValidPassword(password string) bool {
// Password must be at least 6 characters
if len(password) < 6 {
return false
}

// Check for at least one number and one letter
hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

return hasLetter && hasNumber
}

func IsValidName(name string) bool {
return len(strings.TrimSpace(name)) >= 2
}
