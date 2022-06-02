package ierr

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Error struct {
	msg   string
	code  int
	props map[string]interface{}

	next *Error
}

const (
	ID     = "id"
	Amount = "amount"
	Field  = "field"
	Fields = "fields"
	Value  = "value"
)

var (
	ErrIncorrectPassword  = New("incorrect password").BadRequest()
	ErrUserNotFound       = New("user not found").NotFound()
	ErrUnauthorized       = New("user unauthorized").Unauthorized()
	ErrInvalidAccessToken = New("invalid access token").Forbidden()
)

func Internal(err error) *Error {
	return &Error{msg: err.Error()}
}

func New(msg string) *Error {
	return &Error{msg: msg}
}

func Get(err error) *Error {
	if ierr, ok := err.(*Error); ok {
		return ierr
	}

	return Internal(err)
}

func Wrap(err error, msg string) *Error {
	return Get(err).Wrap(msg)
}

func (e *Error) copy() *Error {
	return &Error{
		msg:  e.msg,
		code: e.code,
	}
}

func (e *Error) error() string {
	msg := e.msg
	if e.props != nil {
		props := make([]string, 0, len(e.props))
		for prop, val := range e.props {
			switch v := val.(type) {
			case []string:
				sort.Strings(v)
			}
			props = append(props, fmt.Sprintf("%s: %s", prop, val))
		}
		sort.Strings(props)
		msg += fmt.Sprintf(" { %s }", strings.Join(props, ", "))
	}

	return msg
}

func (e *Error) Error() string {
	msgs := []string{e.error()}
	for ierr := e; ierr.next != nil; {
		ierr = ierr.next
		msgs = append(msgs, ierr.error())
	}
	sort.Strings(msgs)

	return strings.Join(msgs, "; ")
}

func (e *Error) Code() int {
	if e.code == 0 {
		return http.StatusInternalServerError
	}
	return e.code
}

func (e *Error) Is(err error) bool {
	return e.msg == err.Error()
}

func (e *Error) Wrap(msg string) *Error {
	err := e.copy()
	err.msg = msg + ": " + e.Error()
	return err
}

func (e *Error) Add(err error) *Error {
	if e == nil {
		return Get(err)
	}

	var ierr *Error
	for ierr = e; ierr.next != nil; ierr = ierr.next {
	}

	ierr.next = Get(err)
	return e
}

func (e Error) WithProperty(prop string, val interface{}) *Error {
	if e.props == nil {
		e.props = make(map[string]interface{})
	}
	e.props[prop] = val
	return &e
}

func (e Error) WithProperties(props ...interface{}) *Error {
	if e.props == nil {
		e.props = make(map[string]interface{}, len(props)/2)
	}

	if len(props)%2 == 1 {
		props = append(props, "")
	}

	for i := 0; i < len(props); i += 2 {
		e.props[props[i].(string)] = props[i+1]
	}

	return &e
}

func (e Error) BadRequest() *Error {
	return e.updateCode(http.StatusBadRequest)
}
func (e Error) NotFound() *Error {
	return e.updateCode(http.StatusNotFound)
}
func (e Error) UnprocessableEntity() *Error {
	return e.updateCode(http.StatusUnprocessableEntity)
}
func (e Error) Unauthorized() *Error {
	return e.updateCode(http.StatusUnauthorized)
}
func (e Error) Forbidden() *Error {
	return e.updateCode(http.StatusForbidden)
}

func (e *Error) isInternal() bool {
	return e.code == 0
}

func (e *Error) updateCode(code int) *Error {
	if e.isInternal() {
		e.code = code
	}
	return e
}
