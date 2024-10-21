package errs

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
)

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Other          Kind = iota // Unclassified error. DefaultValue
	Invalid                    // Invalid operation for this type of item.
	IO                         // External I/O error such as network failure.
	Exist                      // Item already exists.
	NotExist                   // Item does not exist.
	Private                    // Information withheld.
	Internal                   // Internal error or inconsistency.
	BrokenLink                 // Link target does not exist.
	Database                   // Error from database.
	Validation                 // Input validation error.
	Unanticipated              // Unanticipated error.
	InvalidRequest             // Invalid Request
	// Unauthenticated is used when a request lacks valid authentication credentials.
	//
	// For Unauthenticated errors, the response body will be empty.
	// The error is logged and http.StatusUnauthorized (401) is sent.
	Unauthenticated // Unauthenticated Request
	// Unauthorized is used when a user is authenticated, but is not authorized
	// to access the resource.
	//
	// For Unauthorized errors, the response body should be empty.
	// The error is logged and http.StatusForbidden (403) is sent.
	Unauthorized
	UnsupportedMediaType // Unsupported Media Type
)

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// UserName is a string representing a user
type UserName string

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

// Parameter represents the parameter related to the error.
type Parameter string

// Code is a human-readable, short representation of the error
type Code string

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked.
	Op Op
	// User is the name of the user attempting the operation.
	User UserName
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// Param represents the parameter related to the error.
	Param Parameter
	// Code is a human-readable, short representation of the error
	Code Code
	// The underlying error that triggered this one, if any.
	Err error
}

func (e *Error) isZero() bool {
	return e.User == "" && e.Kind == 0 && e.Param == "" && e.Code == "" && e.Err == nil
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//
//	UserName
//		The username of the user attempting the operation.
//	string
//		Treated as an error message and assigned to the
//		Error field after a call to errors.New.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
func E(op Op, err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	e := &Error{
		Op:  op,
		Err: err,
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case UserName:
			e.User = arg
		case Code:
			e.Code = arg
		case Parameter:
			e.Param = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	if prev.Code == e.Code {
		prev.Code = ""
	}
	if e.Code == "" {
		e.Code = prev.Code
		prev.Code = ""
	}

	if prev.Param == e.Param {
		prev.Param = ""
	}
	if e.Param == "" {
		e.Param = prev.Param
		prev.Param = ""
	}

	return e
}

// OpStack returns the op stack information for an error
func OpStack(err error) []string {
	e := err
	stack := make([]string, 0)

	for errors.Unwrap(e) != nil {
		var errsError *Error

		if errors.As(e, &errsError) && errsError.Op != "" {
			stack = append(stack, e.Error())
		}

		e = errors.Unwrap(e)
	}

	slices.Reverse(stack)
	return stack
}

// TopError recursively unwraps all errors and retrieves the topmost error
func TopError(err error) error {
	currentErr := err
	for errors.Unwrap(currentErr) != nil {
		currentErr = errors.Unwrap(currentErr)
	}

	return currentErr
}

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case BrokenLink:
		return "link target does not exist"
	case Private:
		return "information withheld"
	case Internal:
		return "internal error"
	case Database:
		return "database error"
	case Validation:
		return "input validation error"
	case Unanticipated:
		return "unanticipated error"
	case InvalidRequest:
		return "invalid request error"
	case Unauthenticated:
		return "unauthenticated request"
	case Unauthorized:
		return "unauthorized request"
	case UnsupportedMediaType:
		return "unsupported media type"
	default:
		return "unknown error kind"
	}
}

func KindIs(err error, kind Kind) bool {
	var e *Error

	if errors.As(err, &e) {
		if e.Kind != Other {
			return e.Kind == kind
		}
	}

	return false
}
