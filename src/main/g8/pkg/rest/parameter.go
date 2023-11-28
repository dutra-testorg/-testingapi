package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Gympass/gcore/v3/gerror"
	"github.com/gorilla/mux"

	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

var (
	// ErrMissingParameter is returned when parameter key is not found on Query Parameters.
	ErrMissingParameter = gerror.NewBadRequest(errors.New("param is missing"))
)

// GetInt gets a query or path parameter as integer.
func GetInt(r *http.Request, param string) (int, error) {
	p, err := getParam(r, param)
	if err != nil {
		return 0, err
	}

	v, err := atoi(param, p)
	if err != nil {
		return 0, badRequestErr(err, "%s must be integer", param)
	}

	return v, nil
}

// GetIntOrDefault gets a query or path parameter as integer or default.
func GetIntOrDefault(r *http.Request, param string, defaultValue int) (int, error) {
	p, err := getParam(r, param)
	if err != nil {
		return defaultValue, err
	}

	v, err := atoi(param, p)
	if err != nil {
		return defaultValue, err
	}

	return v, nil
}

// GetInt64 gets a query or path parameter as int64.
func GetInt64(r *http.Request, param string) (int64, error) {
	p, err := getParam(r, param)
	if err != nil {
		return 0, err
	}

	return atoi64(param, p)
}

// GetUUID gets a query or path parameter as UUID.
func GetUUID(r *http.Request, param string) (uuid.UUID, error) {
	p, err := getParam(r, param)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.FromString(p)
	if err != nil {
		return uuid.Nil, badRequestErr(err, "%s must be an UUIDv4", param)
	}

	return id, nil
}

// GetString gets a query or path parameter as string.
func GetString(r *http.Request, param string) (string, error) {
	return getParam(r, param)
}

// GetStringList gets a query parameter as string list.
func GetStringList(r *http.Request, param string) ([]string, error) {
	var empty = []string{}

	p, err := getParam(r, param)
	if err != nil {
		return empty, err
	}

	return strings.Split(p, ","), nil
}

// GetBool gets a query parameter as boolean.
func GetBool(r *http.Request, param string) (bool, error) {
	p, err := getParam(r, param)
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(p)
	if err != nil {
		return false, badRequestErr(err, "%s must be a bool", param)
	}

	return b, nil
}

// GetDateOrDefault gets a query parameter as date.
func GetDateOrDefault(r *http.Request, param string) (time.Time, error) {
	var empty = time.Time{}

	p, err := getParam(r, param)
	if err != nil {
		return empty, err
	}

	seconds, err := atoi64(param, p)
	if err != nil {
		return empty, err
	}

	return time.Unix(seconds, 0).UTC(), nil
}

// GetUUIDList gets a query parameter as UUID list.
func GetUUIDList(r *http.Request, param string) ([]uuid.UUID, error) {
	var empty = []uuid.UUID{}

	p, err := getParam(r, param)
	if err != nil {
		return empty, err
	}

	paramIDs := strings.Split(p, ",")
	ids := make([]uuid.UUID, len(paramIDs))
	setIDs := make(map[uuid.UUID]interface{})

	for i := range paramIDs {
		id, err := uuid.FromString(paramIDs[i])
		if err != nil {
			return empty, badRequestErr(err, "%s must be be an UUIDv4 list", param)
		}

		if _, found := setIDs[id]; found {
			continue
		}

		setIDs[id] = nil
		ids[i] = id
	}

	return ids[:len(setIDs)], nil
}

// GetTime gets a query or path parameter as RFC3339 Time.
func GetTime(r *http.Request, param string) (time.Time, error) {
	return getTime(r, param, time.RFC3339)
}

// GetDate gets a query parameter as date with no default value.
func GetDate(r *http.Request, param string) (time.Time, error) {
	return getTime(r, param, "2006-01-02")
}

func getTime(r *http.Request, param, timeLayout string) (time.Time, error) {
	var empty time.Time
	p, err := getParam(r, param)
	if err != nil {
		return empty, err
	}

	t, err := time.Parse(timeLayout, p)
	if err != nil {
		return empty, badRequestErr(err, "%s must be a valid data format (%s)", param, timeLayout)
	}

	return t, nil
}

func atoi(param, value string) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, badRequestErr(err, "%s must be an int", param)
	}

	return v, nil
}

func atoi64(param, value string) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, badRequestErr(err, "%s must be an int64", param)
	}

	return v, nil
}

func missingParameterErr(param string) error {
	return ErrMissingParameter.WithMessage(fmt.Sprintf("missing %s parameter", param))
}

func badRequestErr(err error, format string, any ...interface{}) error {
	return gerror.NewBadRequest(err).WithMessage(fmt.Sprintf(format, any...))
}

func getParam(r *http.Request, param string) (string, error) {
	v := r.URL.Query().Get(param)
	if v == "" {
		v, ok := mux.Vars(r)[param]
		if !ok {
			return "", missingParameterErr(param)
		}
		return v, nil
	}

	return v, nil
}
