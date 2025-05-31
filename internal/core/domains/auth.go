package domains

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
