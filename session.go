package session

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/Unknwon/com"
	"github.com/codinl/go-logger"
	"github.com/codinl/martini"
	"github.com/codinl/ttlmap"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	SESSION_USER                = "session_user"
	COOKIE_NAME                 = "sessionid"
	MAX_AGE       time.Duration = 3600
)

var sessions *ttlmap.TtlMap

// Handler is a middleware that maps a Session service into the Martini handler chain.
func Handler(newUser func() User, options ...Option) martini.Handler {
	var option Option
	if len(options) == 0 {
		option = Option{
			CookieName: COOKIE_NAME,
			MaxAge:     MAX_AGE,
		}
	}
	if sessions == nil {
		sessions = ttlmap.NewTtlMap(time.Second * option.MaxAge)
	}

	return func(resp http.ResponseWriter, req *http.Request, ctx martini.Context) {
		session, err := getSession(resp, req, option)
		if err != nil {
			logger.Error(err)
			ctx.Next()
			return
		}

		ctx.MapTo(session, (*Session)(nil))

		u, found := session.Get(SESSION_USER)
		if !found {
			u = newUser()
			session.Set(SESSION_USER, u)
		}

		ctx.MapTo(u, (*User)(nil))

		ctx.Next()
	}
}

func getSession(resp http.ResponseWriter, req *http.Request, option Option) (Session, error) {
	var sid string
	c, err := req.Cookie(option.CookieName)
	if err == nil {
		sid, _ = url.QueryUnescape(c.Value)
		s, found := sessions.Get(sid)
		if found {
			return s.(*session), nil
		}
	} else {
		logger.Error(err)
		sid = genSessionId()
		setCookie(resp, req, sid, option)
	}

	session := NewSession(req, resp, sid)

	sessions.Set(sid, session)

	return session, nil
}

func setCookie(resp http.ResponseWriter, req *http.Request, sid string, option Option) {
	cookie := &http.Cookie{
		Name:     option.CookieName,
		Value:    sid,
		HttpOnly: option.HttpOnly,
		Secure:   option.Secure,
		MaxAge:   int(option.MaxAge),
	}

	http.SetCookie(resp, cookie)
}

func genSessionId() string {
	return hex.EncodeToString(generateRandomKey(16))
}

// generateRandomKey creates a random key with the given strength.
func generateRandomKey(strength int) []byte {
	k := make([]byte, strength)
	if n, err := io.ReadFull(rand.Reader, k); n != strength || err != nil {
		return com.RandomCreateBytes(strength)
	}
	return k
}

// Option --------------------------------------------------------------------

// Option stores configuration for a session or session store.
//
// Fields are a subset of http.Cookie fields.
type Option struct {
	CookiePath string
	Domain     string
	CookieName string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   time.Duration
	Secure   bool
	HttpOnly bool
}

// Session --------------------------------------------------------------------

type Session interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, found bool)
	Delete(key string)
}

// Session stores the values and optional configuration for a session.
type session struct {
	mutex  sync.RWMutex
	Req    *http.Request
	Resp   http.ResponseWriter
	Sid    string
	Option *Option
	data   map[string]interface{}
}

func NewSession(req *http.Request, resp http.ResponseWriter, sid string) Session {
	logger.Debug("--------NewSession()")
	return &session{
		Req:  req,
		Resp: resp,
		Sid:  sid,
		data: make(map[string]interface{}),
	}
}

func (s *session) Set(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}

func (s *session) Get(key string) (value interface{}, found bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if value, ok := s.data[key]; ok {
		return value, ok
	}
	return nil, false
}

func (s *session) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}
