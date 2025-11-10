package routerx

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Ctx struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	handlers []Handler
	index    int
	status   int
}

func NewCtx(w http.ResponseWriter, r *http.Request, hs []Handler) *Ctx {
	return &Ctx{
		Writer:   w,
		Request:  r,
		handlers: hs,
	}
}

// BodyParser parse the json request body.
// Make sure that body paramter is a pointer.
func (c *Ctx) BodyParser(body any) error {
	return json.NewDecoder(c.Request.Body).Decode(body)
}

// Context return the request context.
func (c *Ctx) Context() context.Context {
	return c.Request.Context()
}

// WithContext set the request context.
func (c *Ctx) WithContext(ctx context.Context) {
	c.Request = c.Request.WithContext(ctx)
}

// Next call the next handler
func (c *Ctx) Next() error {
	if c.index >= len(c.handlers) {
		return nil
	}

	handler := c.handlers[c.index]
	c.index++

	return handler(c)
}

// Status set response status code
func (c *Ctx) Status(code int) *Ctx {
	c.status = code
	return c
}

// SendStatus send response just only status code
func (c *Ctx) SendStatus(code int) error {
	c.Writer.WriteHeader(code)
	return nil
}

// Json send json response
func (c *Ctx) Json(data any) error {
	code := http.StatusOK
	if c.status != 0 {
		code = c.status
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(data)
}

// Query get the query parameter value by the given key.
// If there no values associated with the key, Query returns the default value.
func (c *Ctx) Query(key string, defaultVal ...string) string {
	val := c.Request.URL.Query().Get(key)

	if len(defaultVal) > 0 && val == "" {
		return defaultVal[0]
	}

	return val
}

// QueryInt returns the query paramter value int by the given key.
// If there no values associated with the key, QueryInt returns the default value.
// If the default falue is not provided then it returns the zero value of int.
func (c *Ctx) QueryInt(key string, defaultVal ...int) int {
	defaultValue := "0"
	if len(defaultVal) > 0 {
		defaultValue = strconv.Itoa(defaultVal[0])
	}

	valStr := c.Query(key, defaultValue)

	val, _ := strconv.Atoi(valStr)

	return val
}

// QueryBool returns the query parameter value bool by the given key.
// If there no values associated with the key, QueryBool returns the default value.
// If the default value is not provided then it returns the zero value of bool.
func (c *Ctx) QueryBool(key string, defaulVal ...bool) bool {
	defaultValue := "false"
	if len(defaulVal) > 0 {
		defaultValue = strconv.FormatBool(defaulVal[0])
	}

	valStr := c.Query(key, defaultValue)

	val, _ := strconv.ParseBool(valStr)

	return val
}
