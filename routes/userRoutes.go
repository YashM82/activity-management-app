package routes

import (
	"github.com/gorilla/mux"
	user "todo/controller"

)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/signup", user.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", user.LoginHandler).Methods("POST")
	//r.HandleFunc("/logout", user.LogoutHandler).Methods("POST")
}
