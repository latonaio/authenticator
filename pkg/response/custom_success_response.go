package custmres

type JWTResponseFormat struct {
	Jwt string `json:"jwt"`
}

type UserResponseFormat struct {
	LoginID string `json:"login_id"`
}
