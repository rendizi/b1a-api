package handler

import (
	"b1a/internal/db"
	"b1a/internal/generator"
	"b1a/internal/url"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
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
	if _, ok := data["url"]; !ok {
		http.Error(w, "no url provided", http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("Authorization")
	var email string
	if tokenString != "" {
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
		email, ok = claims["email"].(string)
		if !ok {
			http.Error(w, "Error extracting email from token", http.StatusInternalServerError)
			return
		}
	} else {
		email = "users@gmail.com"
	}

	if !url.IsVal(data["url"]) {
		http.Error(w, "not valid url", http.StatusBadRequest)
		return
	}
	prefered, ok := data["prefered"]
	short := ""
	if ok {
		if db.IsFreeUrl(prefered) == nil {
			short = prefered
			err = db.InsertUrl(data["url"], short, email, data["sharedWith"], data["topic"], data["message"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dat := struct {
				Email    string `json:"email"`
				ShortUrl string `json:"shorturl"`
			}{
				Email:    email,
				ShortUrl: short,
			}

			jsonData, err := json.Marshal(dat)
			if err != nil {
				http.Error(w, "Error marshalling json", http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		} else {
			http.Error(w, "prefered url already in use", http.StatusBadRequest)
			return
		}
	} else {
		if is, stored, err := db.IsInDb(data["url"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if is {
			short = stored
		} else {
			for {
				short = generator.Do()
				if db.IsFreeUrl(short) == nil {
					break
				}
			}
			err = db.InsertUrl(data["url"], short, email, data["sharedWith"], data["topic"], data["message"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	dat := struct {
		Email    string `json:"email"`
		ShortUrl string `json:"shorturl"`
	}{
		Email:    email,
		ShortUrl: short,
	}

	jsonData, err := json.Marshal(dat)
	if err != nil {
		http.Error(w, "Error marshalling json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	trimmedPath := path[1:]

	url, err := db.GetUrl(trimmedPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		data := struct {
			Email string `json:"email"`
			Url   string `json:"url"`
		}{
			Email: "no",
			Url:   url,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error marshalling json", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
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
	err = db.UpdateHistory(trimmedPath, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Email string `json:"email"`
		Url   string `json:"url"`
	}{
		Email: email,
		Url:   url,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if len(url) == 0 {
		http.Error(w, "no url found", http.StatusBadRequest)
		return
	}
	dat, err := db.GetInfo(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data := struct {
		Url     string `json:"url"`
		Topic   string `json:"topic"`
		Message string `json:"message"`
		Clicked string `json:"clicked"`
	}{
		Url:     url,
		Topic:   dat[0],
		Message: dat[1],
		Clicked: dat[2],
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshalling json", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
