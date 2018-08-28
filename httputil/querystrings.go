package httputil

import (
	"net/http"
	"net/url"
	"time"

	"github.com/ghetzel/go-stockutil/stringutil"
)

// Parses the named query string from a request as an integer.
func QInt(req *http.Request, key string, fallbacks ...int64) int64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToInteger(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a float.
func QFloat(req *http.Request, key string, fallbacks ...float64) float64 {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToFloat(v); err == nil {
			return i
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return 0
	}
}

// Parses the named query string from a request as a date/time value.
func QTime(req *http.Request, key string) time.Time {
	if v := Q(req, key); v != `` {
		if i, err := stringutil.ConvertToTime(v); err == nil {
			return i
		}
	}

	return time.Time{}
}

// Parses the named query string from a request as a boolean value.
func QBool(req *http.Request, key string) bool {
	if v, err := stringutil.ConvertToBool(Q(req, key)); err == nil {
		return v
	}

	return false
}

// Parses the named query string from a request as a string.
func Q(req *http.Request, key string, fallbacks ...string) string {
	if v := req.URL.Query().Get(key); v != `` {
		if vS, err := url.QueryUnescape(v); err == nil {
			return vS
		}
	}

	if len(fallbacks) > 0 {
		return fallbacks[0]
	} else {
		return ``
	}
}

// Sets a query string to the given value in the given url.URL
func SetQ(u *url.URL, key string, value interface{}) {
	qs := u.Query()
	qs.Set(key, stringutil.MustString(value))
	u.RawQuery = qs.Encode()
}

// Appends a query string from then given url.URL
func AddQ(u *url.URL, key string, value interface{}) {
	qs := u.Query()
	qs.Add(key, stringutil.MustString(value))
	u.RawQuery = qs.Encode()
}

// Deletes a query string from then given url.URL
func DelQ(u *url.URL, key string) {
	qs := u.Query()
	qs.Del(key)
	u.RawQuery = qs.Encode()
}
