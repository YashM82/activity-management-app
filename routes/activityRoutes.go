package routes

import (
	"net/http"
	"todo/controller"
	"todo/utils"

	"github.com/gorilla/mux"
)

func RegisterActivityRoutes(r *mux.Router) {
	r.Handle("/addtask", utils.AuthMiddleware(http.HandlerFunc(controller.AddTaskHandler))).Methods("POST")
	r.Handle("/getalltasks", utils.AuthMiddleware(http.HandlerFunc(controller.GetAllTasks))).Methods("GET")
	// r.Handle("/gettaskid", utils.AuthMiddleware(http.HandlerFunc(controller.GetTaskId))).Methods("GET")
	r.Handle("/edittask", utils.AuthMiddleware(http.HandlerFunc(controller.EditTaskHandler))).Methods("POST")
	r.Handle("/deletetask", utils.AuthMiddleware(http.HandlerFunc(controller.DeleteTaskHandler))).Methods("POST")
}
