package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"todo/db"
	"todo/utils"

	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
)

type SignUpInput struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	EmailId   string `json:"emailId" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type LoginInput struct {
	EmailId  string `json:"emailId" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var signupinput SignUpInput
	err := json.NewDecoder(r.Body).Decode(&signupinput)
	if err != nil {
		log.Println("Error while decoding:", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(signupinput)
	if err != nil {
		log.Println("Invalid SignUp Input")
		return
	}

	// Encrpyt pass
	hashedPassword, err := argon2id.CreateHash(signupinput.Password, argon2id.DefaultParams)
	if err != nil {
		log.Println("Failed to hash password")
		return
	}
	signupinput.Password = hashedPassword

	// Insert into DB
	query := "INSERT INTO TodoUsers (first_name, last_name, password,emailID) VALUES (?, ?, ?,?)"
	stmt, err := db.MySqlSession.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(signupinput.FirstName, signupinput.LastName, signupinput.Password, signupinput.EmailId)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginInput LoginInput
	err := json.NewDecoder(r.Body).Decode(&loginInput)
	if err != nil {
		log.Println("failed to decode req body", err)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginInput)
	if err != nil {
		log.Println("Invalid Input")
		http.Error(w, "Invalid login details", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	query := "SELECT password FROM TodoUsers WHERE emailID=?"
	err = db.MySqlSession.QueryRow(query, loginInput.EmailId).Scan(&hashedPassword)
	if err != nil {
		log.Println("error while query DB", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}
	match, err := argon2id.ComparePasswordAndHash(loginInput.Password, hashedPassword)
	if err != nil || !match {
		log.Println("Invalid Password")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(loginInput.EmailId)
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
