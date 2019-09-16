// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatedUser struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
