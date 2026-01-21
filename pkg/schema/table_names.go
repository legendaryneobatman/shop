package schema

var TableNames = struct {
	User         string
	List         string
	Todo         string
	RefreshToken string
}{
	User:         "users",
	List:         "lists",
	Todo:         "todos",
	RefreshToken: "refresh_tokens",
}
