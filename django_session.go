package django_session

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type (
	SessionMiddleWare interface {
		Filter(next http.Handler) http.Handler
	}

	BaseSession struct {
		UserID      string `json:"_auth_user_id"`
		UserBackend string `json:"_auth_user_backend"`
		UserHash    string `json:"_auth_user_hash"`
	}

	Store interface {
		Fetch(ctx context.Context, key string, dest interface{}) error
	}

	DjangoSession struct {
		Store            Store
		OnError          http.HandlerFunc
		OnInvalidSession http.HandlerFunc
	}
)

var (
	ErrSessionInvalid = errors.New("invalid session")
)

const (
	SessionContextKey = "github.com/fdelbos/django-session.context_key"
)

func (ds DjangoSession) Filter(cookieName string, dest interface{}) func(http.Handler) http.Handler {
	destType := reflect.TypeOf(dest)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			session := reflect.New(destType).Interface()
			err = ds.Store.Fetch(r.Context(), cookie.Value, session)

			if errors.Is(err, ErrSessionInvalid) {
				ds.OnError(w, r)

			} else if err != nil {
				ds.OnInvalidSession(w, r)

			} else {
				ctx := context.WithValue(r.Context(), SessionContextKey, session)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func GetSession(r *http.Request) interface{} {
	return r.Context().Value(SessionContextKey)
}

func decodeString(str string) ([]byte, error) {
	// discard timestamp and signature if any
	str = strings.Split(str, ":")[0]

	// the cookie can be compressed with zlib, in that case a '.' prefix will be present
	// see https://github.com/django/django/blob/master/django/core/signing.py
	compressed := false
	if str[0] == '.' {
		compressed = true
		str = str[1:]
	}

	res, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	if compressed {
		b := bytes.NewBuffer([]byte(res))
		dest := bytes.Buffer{}
		if r, err := zlib.NewReader(b); err != nil {
			return nil, err

		} else if _, err := io.Copy(&dest, r); err != nil {
			return nil, err

		}
		res = dest.Bytes()
	}

	return res, nil
}
