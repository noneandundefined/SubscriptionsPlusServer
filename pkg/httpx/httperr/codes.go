package httperr

import "errors"

var (
	Err_DuplicateEmail          = errors.New("a user with this email address already exists")
	Err_UserNotFound            = errors.New("user was not found")
	Err_ContextDeadlineExceeded = errors.New("database response timeout, please try again later")
	Err_ContextCanceled         = errors.New("operation was canceled")
	Err_UniqueViolation         = errors.New("entered data already exists")
	Err_DbTimeout               = errors.New("database connection timed out")
	Err_DbNetworkTemporary      = errors.New("temporary network issue. Please try again")
	Err_DbNetwork               = errors.New("network error while connecting to database")
	Err_NotDeleted              = errors.New("error: no record was deleted")
	Err_NotUpdated              = errors.New("error: no record was updated")
)

var (
	Err_RedisKeyNotFound      = errors.New("the requested data was not found")
	Err_RedisTimeout          = errors.New("the request timed out")
	Err_RedisDeadlineExceeded = errors.New("the operation took too long")
	Err_RedisCanceled         = errors.New("the request was canceled")
	Err_RedisNetwork          = errors.New("a network error occurred")
	Err_RedisOperationFailed  = errors.New("operation failed")
)
