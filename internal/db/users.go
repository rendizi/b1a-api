package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

func InsertUser(email, password string) error {
	insertQuery := `INSERT INTO users (email, password, lastVisited, sharedWithMe) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(insertQuery, email, password, "{}", "{}")
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser(email, password string) error {
	var storedPassword string
	// Execute the SQL query to fetch the password associated with the email
	err := db.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("email not found")
		}
		return err
	}

	// Verify if the provided password matches the stored password
	if storedPassword != password {
		return errors.New("incorrect password")
	}

	return nil // User validated successfully
}

func GetHistory(email string) ([]string, error) {
	var storedHistory []byte
	err := db.QueryRow("SELECT lastVisited FROM users WHERE email = $1 LIMIT 1", email).Scan(&storedHistory)
	if err != nil {
		return nil, err
	}

	// Check if storedLinks is nil
	if storedHistory == nil {
		return nil, nil
	}

	// Convert []byte to []string
	links := strings.Split(string(storedHistory[1:len(storedHistory)-1]), ",")
	return links, nil
}

func GetShared(email string) ([]string, error) {
	var storedLinks []byte
	err := db.QueryRow("SELECT sharedWithMe FROM users WHERE email = $1 LIMIT 1", email).Scan(&storedLinks)
	if err != nil {
		return nil, err
	}

	// Check if storedLinks is nil
	if storedLinks == nil {
		return nil, nil
	}

	// Convert []byte to []string
	links := strings.Split(string(storedLinks[1:len(storedLinks)-1]), ",")
	return links, nil
}

func UpdateHistory(url, email string) error {
	// Check if the array is full
	fullQuery := `SELECT array_length(lastVisited, 1) >= 5 FROM users WHERE email = $1`
	var isFull bool
	err := db.QueryRow(fullQuery, email).Scan(&isFull)
	if err != nil {
		return err
	}

	log.Println("Is array full? ", isFull)
	if isFull {
		// If the array is full, remove the first element before appending
		updateQuery := `
            UPDATE users 
            SET lastVisited = array_append(lastVisited[2:5], $1) 
            WHERE email = $2
        `
		_, err := db.Exec(updateQuery, url, email)
		if err != nil {
			return err
		}
	} else {
		// If the array is not full, simply append the new element
		updateQuery := `UPDATE users SET lastVisited = array_append(lastVisited, $1) WHERE email = $2`
		_, err := db.Exec(updateQuery, url, email)
		if err != nil {
			return err
		}
	}
	return nil
}
