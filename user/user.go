package user

// New creates new User.
func New() *User {
	return new(User)
}

// User is a struct which represents a user.
type User struct {
	Username string
	Password string
	APIToken string
}

// Login sets fields in User struct.
func (u *User) Login(name, pass, token string) {
	u.Username = name
	u.Password = pass
	u.APIToken = token
}
