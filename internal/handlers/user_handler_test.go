package handlers

import (
"testing"
"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
router := gin.Default()
assert.NotNil(t, router)
}
