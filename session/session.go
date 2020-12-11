package session

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

type Manager struct {
	cookieName  string     // private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxLifeTime int64
}

var provides = make(map[string]Provider)

// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

func (manager *Manager) sessionId() string {
	return uuid.NewV4().String()
}

func (manager *Manager) SessionStart(c *gin.Context) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := c.Cookie(manager.cookieName)
	if err != nil || cookie == "" {
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		c.SetCookie(manager.cookieName, url.QueryEscape(sid), int(manager.maxLifeTime), "/", "localhost", false, true)
	} else {
		sid, _ := url.QueryUnescape(cookie)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

//Destroy sessionid
func (manager *Manager) SessionDestroy(c *gin.Context) {
	cookie, err := c.Cookie(manager.cookieName)
	if err != nil || cookie == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie)
		c.SetCookie(manager.cookieName, "", -1, "/", "localhost", false, true)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
