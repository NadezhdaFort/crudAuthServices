package authclient

type UserIdRole struct {
	UserId   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type HTTPResponse struct {
	Success  bool       `json:"success"`
	Error    string     `json:"error,omitempty"`
	UserInfo UserIdRole `json:"data"`
}
