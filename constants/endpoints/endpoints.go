package endpoints

import "fmt"

type userEndpoints struct {
	Register string
	SignIn   string
}

var user = "/user"

var UserEndpoints = userEndpoints{
	Register: fmt.Sprintf("%s/register", user),
	SignIn:   fmt.Sprintf("%s/signin", user),
}
