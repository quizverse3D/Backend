package authgateway

type Auth struct {
	ID           string
	Email        string
	Password     string
	PasswordSalt string
}
