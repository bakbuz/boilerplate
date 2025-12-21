package service

type AuthenticationService interface {
}

type authenticationService struct {
}

// NewAuthenticationService ...
func NewAuthenticationService() AuthenticationService {
	return &authenticationService{}
}
