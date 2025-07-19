package router

type Router struct {
	addr string
	Todo *Todo
	User *User
}

func NewRouter(
	todo *Todo,
	user *User,
) *Router {
	return &Router{
		Todo: todo,
		User: user,
	}
}
