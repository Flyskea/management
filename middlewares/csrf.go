package middlewares

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"io"
	"manage/logger"
	"manage/utils"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Options the options of csrf middleware
type Options struct {
	// maximum age for this token to live
	MaxAge int
	// name of the cookie to keep csrf token
	CookieName string
	// name of the header which the csrf token is sending back
	HeaderName string
	// for setting the cookie
	Secure bool
	// name for keeping csrf token in session
	SessionName string
	// name for keeping issued time in session
	IssuedName string
	// Length of csrf token
	ByteLenth int
	// path which the cookie is valid
	Path string
	// Http methods considered as safe and pass validation
	SafeMethods []string
}

// DefaultOptions get default options
func DefaultOptions() *Options {
	return &Options{
		MaxAge:      60 * 60 * 24 * 180,
		CookieName:  "csrf_token",
		HeaderName:  "X-CSRF-Token",
		Secure:      false,
		SessionName: "csrf_token_session",
		IssuedName:  "csrf_token_issued",
		ByteLenth:   32,
		Path:        "/",
		SafeMethods: []string{"GET", "HEAD", "OPTIONS"},
	}
}

// Csrf ...
func Csrf(options *Options) gin.HandlerFunc {
	if options == nil {
		options = DefaultOptions()
	}

	return func(c *gin.Context) {
		var (
			csrfSession string
			issued      int64
		)

		if utils.InArray(options.SafeMethods, c.Request.Method) {
			c.Next()
			return
		}

		if c.Request.URL.Scheme == "https" {
			referer, err := url.Parse(c.Request.Header.Get("Referer"))
			if err != nil || referer == nil {
				handleError(c, http.StatusBadRequest, nil, "csrf验证错误")
				return
			}
			if !sameOrigin(c.Request.URL, referer) {
				handleError(c, http.StatusBadRequest, nil, "csrf验证错误")
				return
			}
		}
		session := sessions.Default(c)
		csrfCookie, _ := c.Cookie(options.CookieName)

		if csrfCookie == "" {

			logger.Info("csrf_token not found in cookie")
			generateNewCsrfAndHandle(c, session, options)
			return
		}

		if csrfSess := session.Get(options.SessionName); csrfSess != nil {
			csrfSession = csrfSess.(string)
		}

		if csrfSession == "" {
			logger.Info("csrf_token not found in session")
			generateNewCsrfAndHandle(c, session, options)
			return
		}

		// max usage generate new token
		now := time.Now()
		if is := session.Get(options.IssuedName); is != nil {
			issued = is.(int64)
		}
		if now.Unix() > (issued + int64(options.MaxAge)) {
			logger.Info("csrf_token max age. New token required")
			generateNewCsrfAndHandle(c, session, options)
			return
		}

		// compare session with header
		csrfHeader := c.Request.Header.Get(options.HeaderName)
		//log.Println("sess", csrfSession, "cookie", csrfCookie, "csrfHeader", csrfHeader, counter, options.MaxUsage)
		if !isTokenValid(csrfSession, csrfHeader) {
			logger.Info("csrf_token diff. New token required")
			generateNewCsrfAndHandle(c, session, options)
			return
		}
		defer saveSession(session, options, csrfSession, false)
		c.Next()
	}
}

func isTokenValid(csrfSession, csrfHeader string) bool {
	return subtle.ConstantTimeCompare([]byte(csrfSession), []byte(csrfHeader)) == 1
}

func saveSession(session sessions.Session, options *Options, csrfSession string, newCsrfSession bool) {
	if newCsrfSession {
		session.Set(options.SessionName, csrfSession)
		session.Set(options.IssuedName, time.Now().Unix())
		session.Save()
	}
}

func generateNewCsrfAndHandle(c *gin.Context, session sessions.Session, options *Options) {
	csrfSession := newCsrf(c, options.CookieName, options.Path, options.MaxAge, options.ByteLenth, options.Secure)
	saveSession(session, options, csrfSession, true)
	//log.Println("generate new token", csrfSession)
	handleError(c, http.StatusBadRequest, nil, options.CookieName)
}

// GenerateNewCsrf 生成csrftoken
func GenerateNewCsrf(c *gin.Context, session sessions.Session) {
	options := DefaultOptions()
	csrfSession := newCsrf(c, options.CookieName, options.Path, options.MaxAge, options.ByteLenth, options.Secure)
	saveSession(session, options, csrfSession, true)
}

func handleError(c *gin.Context, statusCode int, data gin.H, message string) {
	c.Abort()
	utils.Response(c, statusCode, data, message)
}

func newCsrf(c *gin.Context, cookieName, path string, maxAge, byteLenth int, secure bool) string {
	csrfCookie := generateToken(byteLenth)
	c.SetCookie(cookieName, csrfCookie, maxAge, path, "", secure, false)
	return csrfCookie
}

func generateToken(byteLenth int) string {
	result := make([]byte, byteLenth)
	if _, err := io.ReadFull(rand.Reader, result); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(result)
}

func originOK(u *url.URL) bool {
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func sameOrigin(a, b *url.URL) bool {
	if !originOK(b) {
		return false
	}
	return (a.Scheme == b.Scheme && a.Host == b.Host)
}
