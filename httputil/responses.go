package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/ghetzel/go-stockutil/maputil"
	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
)

var FormUnmarshalStructTag = `json`

type RequestParseFunc func(*http.Request, interface{}) error

// specifies a mapping between Content-Type (first by exact match, then by just the Media Type component)
// The "empty string" key is the default if no other keys match first.
var parsers = map[string]RequestParseFunc{
	`application/json`: ParseJSONRequest,
	``:                 ParseFormRequest,
}

// Sets a parser implementation for the given HTTP content type.
func SetContentTypeParser(contentType string, parser RequestParseFunc) {
	if parser != nil {
		parsers[contentType] = parser
	}
}

// Marshal the given data as a JSON document and write the output to the given ResponseWriter. If
// a status is given, that will be used as the HTTP response status.  If data is an error, and no
// status is given, the status will be "500 Internal Server Error"; if data is nil, the status will
// be "204 No Content".  The Content-Type of the response is "application/json".
//
func RespondJSON(w http.ResponseWriter, data interface{}, status ...int) {
	w.Header().Set(`Content-Type`, `application/json`)
	headerSent := false

	if err, ok := data.(error); ok && err != nil {
		data = map[string]interface{}{
			`success`: false,
			`error`:   err.Error(),
		}

		if len(status) == 0 {
			status = []int{http.StatusInternalServerError}
		}
	}

	if len(status) > 0 {
		w.WriteHeader(status[0])
		headerSent = true
	}

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			Logger.Warningf("Failed to encode response body: %v", err)
		}
	} else if !headerSent {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Parses the Request as JSON and unmarshals into the given value.
func ParseJSONRequest(req *http.Request, into interface{}) error {
	return ParseJSON(req.Body, into)
}

// Parses a given reader as a JSON document and unmarshals into the given value.
func ParseJSON(r io.Reader, into interface{}) error {
	return json.NewDecoder(r).Decode(into)
}

// Parses a set of values received from an HTML form (usually the value of the
// http.Request.Form property) and unmarshals into the given value.
func ParseFormValues(formValues url.Values, into interface{}) error {
	data := make(map[string]interface{})

	for key, values := range formValues {
		values = sliceutil.CompactString(values)
		data[key] = nil

		switch len(values) {
		case 1:
			data[key] = stringutil.Autotype(values[0])
		default:
			data[key] = sliceutil.Autotype(values)
		}
	}

	return maputil.TaggedStructFromMap(data, into, FormUnmarshalStructTag)
}

// Parses the form values for a Request and unmarshals into the given value.
func ParseFormRequest(req *http.Request, into interface{}) error {
	if err := req.ParseForm(); err == nil {
		if req.Method == `POST` {
			return ParseFormValues(req.PostForm, into)
		}else{
			return ParseFormValues(req.Form, into)
		}
	} else {
		return err
	}
}

// Autodetect the Content-Type of the given request and unmarshals into the
// given value.
func ParseRequest(req *http.Request, into interface{}) error {
	contentType := req.Header.Get(`Content-Type`)

	// look for an exact match first based on the entire Content-Type value (mediatype + params)
	if parser, ok := parsers[contentType]; ok {
		return parser(req, into)
	} else if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
		if parser, ok := parsers[mediaType]; ok {
			return parser(req, into)
		} else if parser, ok := parsers[``]; ok {
			return parser(req, into)
		} else {
			return fmt.Errorf("No parser could be found for Content-Type %q", contentType)
		}
	} else {
		return err
	}
}
