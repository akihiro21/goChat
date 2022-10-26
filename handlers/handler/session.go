package handler

/* https://devcenter.heroku.com/ja/articles/go-sessionsを改変*/

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

const (
	//SessionName to store session under
	SessionName = "go-sessions-demo"
)

var (
	sessionStore sessions.Store
)

func SessionInit() {
	ek, err := determineEncryptionKey()
	if err != nil {
		log.Println("err", err)
	}

	sessionStore = sessions.NewCookieStore(
		[]byte(os.Getenv("SESSION_AUTHENTICATION_KEY")),
		ek,
	)
}

func determineEncryptionKey() ([]byte, error) {
	sek := os.Getenv("SESSION_ENCRYPTION_KEY")
	lek := len(sek)
	switch {
	case lek >= 0 && lek < 16, lek > 16 && lek < 24, lek > 24 && lek < 32:
		return nil, errors.Errorf("SESSION_ENCRYPTION_KEY needs to be either 16, 24 or 32 characters long or longer, was: %d", lek)
	case lek == 16, lek == 24, lek == 32:
		return []byte(sek), nil
	case lek > 32:
		return []byte(sek[0:32]), nil
	default:
		return nil, errors.New("invalid SESSION_ENCRYPTION_KEY: " + sek)
	}

}

func handleSessionError(w http.ResponseWriter, err error) {
	log.Println("err", err)
	http.Error(w, "Application Error", http.StatusInternalServerError)
}

func nowLoginBool(w http.ResponseWriter, r *http.Request) bool {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
	}

	name, found := session.Values["name"].(string)
	if found || name != "" {
		return true
	} else {
		return false
	}
}

func loginSession(name string, w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
		return
	}
	session.Values["name"] = name
	if err := session.Save(r, w); err != nil {
		handleSessionError(w, err)
		return
	}
}

func token(w http.ResponseWriter, r *http.Request) string {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
		return ""
	}

	token, found := session.Values["token"].(string)

	if !found {
		return ""
	}

	return token
}

func sessionName(w http.ResponseWriter, r *http.Request) string {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
		return ""
	}

	name, found := session.Values["name"].(string)

	if !found {
		return ""
	}

	return name
}

func tokenSession(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		handleSessionError(w, err)
		return
	}
	h := md5.New()
	salt := "goUser%^7&8888"
	if _, err = io.WriteString(h, salt+time.Now().String()); err != nil {
		log.Println("err", err)
	}

	token := fmt.Sprintf("%x", h.Sum(nil))
	session.Values["token"] = token
	if err := session.Save(r, w); err != nil {
		handleSessionError(w, err)
		return
	}
}

func contain(tokens []string, token string) bool {
	for _, to := range tokens {
		if to == token {
			return true
		}
	}
	return false
}

func removeToken(tokens []string, token string) (newTokens []string) {
	for _, to := range tokens {
		if to != token {
			newTokens = append(newTokens, to)
		}
	}
	return
}

func deleteSession(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, SessionName)
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		handleSessionError(w, err)
		return
	}
}

func tokenCheck(w http.ResponseWriter, r *http.Request) {
	if token(w, r) == "" {
		tokenSession(w, r)
		tokens = append(tokens, token(w, r))
	} else if contain(tokens, token(w, r)) {
		tokens = removeToken(tokens, token(w, r))
		tokenSession(w, r)
		tokens = append(tokens, token(w, r))
	} else {
		deleteSession(w, r)
		tokenSession(w, r)
		tokens = append(tokens, token(w, r))
		if login := nowLoginBool(w, r); login == false {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	}
}
