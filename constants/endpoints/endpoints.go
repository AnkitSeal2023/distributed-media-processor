package endpoints

import "fmt"

type authEndpoints struct {
	SignUp string
	SignIn string
}

type userEndpoints struct {
	UploadVideo string
}

const auth = "/auth"
const user = "/user"
const apiV1 = "/api/v1"

var Ping = fmt.Sprintf("%s/ping", apiV1)

var AuthEndpoints = authEndpoints{
	SignUp: fmt.Sprintf("%s%s/signup", apiV1, auth),
	SignIn: fmt.Sprintf("%s%s/signin", apiV1, auth),
}

var UserEndpoints = userEndpoints{
	UploadVideo: fmt.Sprintf("%s%s/upload-video", apiV1, user),
}
