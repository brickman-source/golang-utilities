package errors

import (
	goerrors "errors"
	"fmt"
	"io"
)

func New(message string) error {
	return goerrors.New(message)
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error {
	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

func WithCode(cause error, code string, message string) error {
	return &WithCodeError{
		cause:   cause,
		code:    code,
		message: message,
	}
}

type WithCodeError struct {
	cause   error
	code    string
	message string
}

func (w *WithCodeError) Error() string   { return w.code + ": " + w.message }
func (w *WithCodeError) Cause() error    { return w.cause }
func (w *WithCodeError) Code() string    { return w.code }
func (w *WithCodeError) Message() string { return w.message }
func (w *WithCodeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.code+":"+w.message)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
func Code(err error) string {
	type withCodeError interface {
		Code() string
	}

	withCode, ok := err.(withCodeError)
	if !ok {
		return ""
	}
	return withCode.Code()
}
func AsWithCode(err error) *WithCodeError {
	withCode, ok := err.(*WithCodeError)
	if !ok {
		return nil
	}
	return withCode
}