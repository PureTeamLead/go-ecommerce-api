package errs

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrWrongUsername = errors.New("wrong username")
	ErrCreatingUser  = errors.New("failed to create user")
	ErrNoUserFound   = errors.New("no user found")
	ErrDeletingUser  = errors.New("failed to delete user")
	ErrUpdatingUser  = errors.New("failed to update user's info")
)

var (
	ErrNoProductFound  = errors.New("no product found")
	ErrDeletingProduct = errors.New("failed to delete product")
	ErrUpdatingProduct = errors.New("failed to update product's info")
	ErrCreatingProduct = errors.New("failed to create product")
)

var (
	ErrDB         = errors.New("db error")
	ErrHashing    = errors.New("error while hashing password")
	ErrBadRequest = errors.New("got a wrong request")
)
