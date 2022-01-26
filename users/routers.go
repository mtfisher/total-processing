package users

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtfisher/total-processing/core"
)

func AuthedRoutes(router *gin.Engine) {
	router.GET("/", index)
}

func LoginRoutes(router *gin.Engine) {
	router.GET("/login", loginPage)
	router.GET("/register", registerPage)
	router.POST("/register", registerAction)
	router.POST("/login", loginAction)
	router.GET("/logout", logoutAction)
}

func loginPage(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")

	if loggedInInterface.(bool) {
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	} else {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Login",
		})
	}
}

func registerPage(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")

	if loggedInInterface.(bool) {
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())
	} else {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "Register",
		})
	}
}

func registerAction(c *gin.Context) {
	core := c.MustGet(core.CoreKey).(*core.Core)
	username := c.PostForm("username")
	password := c.PostForm("password")
	userRepository := newUserRepository(*core)

	if user, err := userRepository.registerNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		core.SetCacheValue(token, user, 1*time.Hour)

		c.Set(IS_LOGGED_IN_KEY, true)

		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.tmpl", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}

func loginAction(c *gin.Context) {
	core := c.MustGet(core.CoreKey).(*core.Core)
	username := c.PostForm("username")
	password := c.PostForm("password")
	userRepository := newUserRepository(*core)

	user, err := userRepository.getUser(username)

	if err != nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": err.Error()})

		return
	}

	passErr := user.PasswordMatches(password)

	if passErr != nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": err.Error()})

		return
	}

	token := generateSessionToken()
	c.SetCookie("token", token, 3600, "", "", false, true)
	core.SetCacheValue(token, user, 1*time.Hour)

	c.Set(IS_LOGGED_IN_KEY, true)

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func logoutAction(c *gin.Context) {
	core := c.MustGet(core.CoreKey).(*core.Core)
	c.SetCookie("token", "", -1, "", "", false, true)
	core.RemoveCacheKey("token")

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func index(c *gin.Context) {
	loggedInInterface, _ := c.Get("is_logged_in")

	if !loggedInInterface.(bool) {
		location := url.URL{Path: "/login"}
		c.Redirect(http.StatusFound, location.RequestURI())
	} else {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "Main website",
			"is_logged_in": loggedInInterface.(bool),
		})
	}
}

func generateSessionToken() string {
	// This is NOT a secure way of generating session tokens and should only be used for demo purposes
	// DO NOT USE THIS IN PRODUCTION
	return strconv.FormatInt(rand.Int63(), 16)
}
