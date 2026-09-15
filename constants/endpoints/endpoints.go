package endpoints

import "fmt"

type userEndpoints struct {
	Register string
	SignIn   string
}

const user = "/user"
const apiV1 = "/api/v1"

var Ping = fmt.Sprintf("%s/ping", apiV1)

var UserEndpoints = userEndpoints{
	Register: fmt.Sprintf("%s%s/register", apiV1, user),
	SignIn:   fmt.Sprintf("%s%s/signin", apiV1, user),
}
