package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/db"
	"todo/utils"
)

type Activities struct {
	ID   int    `json:"id,omitempty"`
	Task string `json:"task"`
}

var activity Activities

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	claims := utils.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("User Email from Token:", claims.EmailId)

	// Parse the JSON body

	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	query := "INSERT INTO TODOActivities (task) VALUES (?)"
	stmt, err := db.MySqlSession.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(activity.Task)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Failed to Add activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task added successfully"))

}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	claims := utils.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("User Email from Token:", claims.EmailId)

	query := "SELECT * FROM TODOActivities"
	rows, err := db.MySqlSession.Query(query)
	if err != nil {
		log.Println("Error querying database:", err)
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var activities []Activities
	for rows.Next() {
		//var activity Activity
		err := rows.Scan(&activity.ID, &activity.Task)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			return
		}
		activities = append(activities, activity)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func EditTaskHandler(w http.ResponseWriter, r *http.Request) {

	claims := utils.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("User Email from Token:", claims.EmailId)

	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "UPDATE TODOActivities SET task = ? WHERE id = ?"
	stmt, err := db.MySqlSession.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(activity.Task, activity.ID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task updated successfully"))
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	claims := utils.GetUserFromContext(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("User Email from Token:", claims.EmailId)

	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM TODOActivities WHERE id = ?"
	stmt, err := db.MySqlSession.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(activity.ID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task deleted successfully"))

}
