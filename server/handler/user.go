package handler

import (
	"b1a/internal/db"
	"b1a/internal/email"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your_secret_key")

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data := make(map[string]string)
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	if _, ok := data["email"]; !ok {
		http.Error(w, "No email provided", http.StatusBadRequest)
		return
	}
	if _, ok := data["password"]; !ok {
		http.Error(w, "No password provided", http.StatusBadRequest)
		return
	}

	err = db.ValidateUser(data["email"], data["password"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data["email"],
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"message": "Signed-in successful", "token": tokenString}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data := make(map[string]string)

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	if _, ok := data["email"]; !ok {
		http.Error(w, "No email provided", http.StatusBadRequest)
		return
	}
	if _, ok := data["password"]; !ok {
		http.Error(w, "No password provided", http.StatusBadRequest)
		return
	}

	err = email.Verify(data["email"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.InsertUser(data["email"], data["password"])
	if err != nil {
		resp := "something went wrong in db"
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			resp = "email is already in use"
		}
		http.Error(w, err.Error()+resp, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Links(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Token is not valid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Error extracting token claims", http.StatusInternalServerError)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Error extracting email from token", http.StatusInternalServerError)
		return
	}

	history, err := db.GetHistory(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shared, err := db.GetShared(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Email   string   `json:"email"`
		History []string `json:"history"`
		Shared  []string `json:"shared"`
	}{
		Email:   email,
		History: history,
		Shared:  shared,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
