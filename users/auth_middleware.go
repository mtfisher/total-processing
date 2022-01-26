package users

import "github.com/gin-gonic/gin"

const IS_LOGGED_IN_KEY = "is_logged_in"

func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set(IS_LOGGED_IN_KEY, true)
		} else {
			c.Set(IS_LOGGED_IN_KEY, false)
		}
	}
}
