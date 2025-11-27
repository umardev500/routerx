package routerx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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

// Locals sets or gets context values.
// Setter: c.Locals("key", "value")
// Getter: c.Locals("key")
func (c *Ctx) Locals(key any, val ...any) any {
	// Setter
	if len(val) > 0 {
		ctx := context.WithValue(c.Request.Context(), key, val[0])
		c.Request = c.Request.WithContext(ctx)
		return nil
	}

	// Getter
	return c.Request.Context().Value(key)
}

// QueryParser parses the query parameters.
// Make sure the out parameter is a pointer.
func (c *Ctx) QueryParser(out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("out must be a pointer")
	}

	if v.IsNil() {
		return fmt.Errorf("out must not be nil")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("out must be a struct")
	}

	t := v.Type()
	q := c.Request.URL.Query()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("query")
		if !ok {
			continue
		}

		val := q.Get(tag)

		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		if err := setValue(f, val); err != nil {
			return err
		}
	}

	return nil
}

// setValue sets the value of a struct field from a string.
// It supports basic type (string, int, bool) and pointer fields.
// If the field is a pointer, it allocates a new value and set it recursively.
// Returns an error if the value cannot be set.
func setValue(f reflect.Value, val string) error {
	if val == "" {
		return nil
	}

	// Handle pointers
	if f.Kind() == reflect.Pointer {
		elemType := f.Type().Elem()
		ptr := reflect.New(elemType)
		if err := setValue(ptr.Elem(), val); err != nil {
			return err
		}
		f.Set(ptr)
		return nil
	}

	// Handle basic types
	switch f.Kind() {
	case reflect.String:
		f.SetString(val)
	case reflect.Int:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("invalid int value: %v", err)
		}
		f.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid bool value: %v", err)
		}
		f.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type %s", f.Kind())
	}
	return nil
}

// BodyParser parses the JSON request body.
// Make sure the body parameter is a pointer.
func (c *Ctx) BodyParser(body any) error {
	return json.NewDecoder(c.Request.Body).Decode(body)
}

// Context returns the request context.
func (c *Ctx) Context() context.Context {
	return c.Request.Context()
}

// WithContext sets the request context.
func (c *Ctx) WithContext(ctx context.Context) {
	c.Request = c.Request.WithContext(ctx)
}

// Next calls the next handler.
func (c *Ctx) Next() error {
	if c.index >= len(c.handlers) {
		return nil
	}

	handler := c.handlers[c.index]
	c.index++

	return handler(c)
}

// Status sets the response status code.
func (c *Ctx) Status(code int) *Ctx {
	c.status = code
	return c
}

// SendStatus sends a response with only a status code.
func (c *Ctx) SendStatus(code int) error {
	c.Writer.WriteHeader(code)
	return nil
}

// JSON sends a JSON response.
func (c *Ctx) JSON(data any) error {
	code := http.StatusOK
	if c.status != 0 {
		code = c.status
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(data)
}

// Query returns the query paramter value for the given key.
// If there are not values associated with the key, Query returns the default value.
func (c *Ctx) Query(key string, defaultVal ...string) string {
	val := c.Request.URL.Query().Get(key)

	if len(defaultVal) > 0 && val == "" {
		return defaultVal[0]
	}

	return val
}

// QueryInt returns the query paramter value as int for the given key.
// If there are no values associated with the key, QueryInt returns the default value.
// If the default value is not provided, it returns the zero value of int.
func (c *Ctx) QueryInt(key string, defaultVal ...int) int {
	defaultValue := "0"
	if len(defaultVal) > 0 {
		defaultValue = strconv.Itoa(defaultVal[0])
	}

	valStr := c.Query(key, defaultValue)

	val, _ := strconv.Atoi(valStr)

	return val
}

// QueryBool returns the query paramter value as bool for the given key.
// If there are no values associated with the key, QueryBool returns the default value.
// If the default value is not provided, it returns the zero value of bool.
func (c *Ctx) QueryBool(key string, defaulVal ...bool) bool {
	defaultValue := "false"
	if len(defaulVal) > 0 {
		defaultValue = strconv.FormatBool(defaulVal[0])
	}

	valStr := c.Query(key, defaultValue)

	val, _ := strconv.ParseBool(valStr)

	return val
}

// Param returns the URL parameter value for the given key.
// It returns the empty string if there are no values associated with the key.
func (c *Ctx) Param(key string) string {
	return c.Request.PathValue(key)
}
