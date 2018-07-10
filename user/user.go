package user

// New returns a pointer to a new User instance
func New() *User {
	return new(User)
}

// User represents a user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	APIToken string `json:"api_token"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}

// SetLogin sets name and password fields of a User structure
func (u *User) SetLogin(name, pass string) {
	u.Username = name
	u.Password = pass
}
