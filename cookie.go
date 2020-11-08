package django_session

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"errors"
// 	"fmt"
// 	"hash"
// 	"log"
// 	"net/http"
// 	"strings"
// )

// type (
// 	SigningAlgorithm int

// 	CookieSession struct {
// 		SecretKey  string
// 		CookieName string
// 		Salt       string
// 	}

// 	ErrInvalidCookieSession struct {
// 		Reason string
// 	}
// )

// var (
// 	ErrInvalidSignature = errors.New("invalid signature")
// )

// func NewCookieSession(secretKey, cookieName string) *CookieSession {
// 	return &CookieSession{
// 		SecretKey:  secretKey,
// 		CookieName: cookieName,
// 		Salt:       "django.core.signing.get_cookie_signer",
// 	}
// }

// func (cs *CookieSession) Filter(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Do stuff here
// 		log.Println(r.RequestURI)
// 		cookie, err := r.Cookie(cs.CookieName)
// 		if err != nil {
// 			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 			return
// 		}
// 		log.Print(cookie.Raw)
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})

// }

// // func DecodeString(str string) ([]byte, error) {
// // 	// the cookie can be compressed with zlib, in that case a '.' prefix will be present
// // 	// see https://github.com/django/django/blob/master/django/core/signing.py
// // 	compressed := false
// // 	if str[0] == '.' {
// // 		compressed = true
// // 		str = str[1:]
// // 	}

// // 	res, err := base64.RawURLEncoding.DecodeString(str)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	if compressed {
// // 		b := bytes.NewBuffer([]byte(res))
// // 		dest := bytes.Buffer{}
// // 		if r, err := zlib.NewReader(b); err != nil {
// // 			return nil, err

// // 		} else if _, err := io.Copy(&dest, r); err != nil {
// // 			return nil, err

// // 		}
// // 		res = dest.Bytes()
// // 	}

// // 	return res, nil
// // }

// func (cs *CookieSession) ValidateSigning(cookie string) error {
// 	fields := strings.Split(cookie, ":")
// 	if len(fields) != 3 {
// 		return errors.New("invalid cookie format, unable to validate")
// 	}

// 	cookieSignature := fields[2]
// 	value := fmt.Sprintf("%s:%s", fields[0], fields[1])
// 	key := "django.http.cookies" + cs.SecretKey
// 	computedSignature, err := SaltedHMAC(cs.Salt+"signer", value, key, sha256.New())
// 	if err != nil {
// 		return err

// 	} else if cookieSignature != computedSignature {
// 		log.Print(cookieSignature)
// 		log.Print(computedSignature)
// 		return ErrInvalidSignature

// 	}
// 	return nil
// }

// // func (cs *CookieSession) Sign(msg string) (string, error) {
// // 	key := []byte("salt" + cs.SecretKey)
// // 	var h hash.Hash

// // 	switch cs.Algo {
// // 	case SHA1:
// // 		h = hmac.New(sha1.New, key)
// // 	case SHA256:
// // 		hash := sha256.New()
// // 		hash.Write(key)
// // 		h = hmac.New(sha256.New, hash.Sum(nil))
// // 	}
// // 	if _, err := h.Write([]byte(msg)); err != nil {
// // 		return "", err
// // 	}
// // 	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
// // }

// // func (cs *CookieSession) Decrypt(cookie string) error {

// // 	fields := strings.Split(cookie, ":")
// // 	if len(fields) != 3 {
// // 		return errors.New("invalid cookie format, unable to extract value")
// // 	}

// // 	value := fields[0]

// // 	b, err := base64.RawURLEncoding.DecodeString(value)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }

// func SaltedHMAC(salt, value, secret string, algorithm hash.Hash) (string, error) {
// 	if _, err := algorithm.Write([]byte(salt + secret)); err != nil {
// 		return "", err
// 	}
// 	key := algorithm.Sum(nil)
// 	hm := hmac.New(sha256.New, key)
// 	if _, err := hm.Write([]byte(value)); err != nil {
// 		return "", err
// 	}
// 	return base64.RawURLEncoding.EncodeToString(hm.Sum(nil)), nil
// }
