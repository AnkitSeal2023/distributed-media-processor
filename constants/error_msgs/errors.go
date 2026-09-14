package error_msgs

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("User with this email id already exists")
var ErrInternalServer = errors.New("Internal Server Error")
