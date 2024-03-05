package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func InsertUrl(url, shortUrl, author, sharedWith, topic, message string) error {
	insertQuery := `INSERT INTO urls (url, shortUrl, author, sharedWith, topic, message, clicked, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(insertQuery, url, shortUrl, author, sharedWith, topic, message, 0, 0)
	if err != nil {
		return err
	}

	arr := strings.Fields(sharedWith)
	for _, email := range arr {
		// Check if the array is full
		fullQuery := `SELECT array_length(sharedWithMe, 1) >= 5 FROM users WHERE email = $1`
		var isFull bool
		err = db.QueryRow(fullQuery, email).Scan(&isFull)
		if err != nil {
			return err
		}

		if isFull {
			// If the array is full, remove the first element before appending
			updateQuery := `
                UPDATE users 
                SET sharedWithMe = array_append(sharedWithMe[2:5], $1) 
                WHERE email = $2
            `
			_, err = db.Exec(updateQuery, shortUrl, email)
			if err != nil {
				return err
			}
		} else {
			// If the array is not full, simply append the new element
			updateQuery := `UPDATE users SET sharedWithMe = array_append(sharedWithMe, $1) WHERE email = $2`
			_, err = db.Exec(updateQuery, shortUrl, email)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func IsFreeUrl(shortenUrl string) error {
	var storedId int
	err := db.QueryRow("SELECT id FROM urls WHERE shortUrl = $1", shortenUrl).Scan(&storedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return errors.New("already in use")
}

func GetUrl(shortUrl string) (string, error) {
	var storedUrl string
	err := db.QueryRow("SELECT url FROM urls WHERE shortUrl = $1", shortUrl).Scan(&storedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("no such shorten url")
		}
		return "", err
	}

	// Update the clicked count in the database
	_, err = db.Exec("UPDATE urls SET clicked = clicked + 1 WHERE shortUrl = $1", shortUrl)
	if err != nil {
		return "", err
	}

	return storedUrl, nil
}

func GetInfo(shortUrl string) ([]string, error) {
	var topic, message string
	var clicked int
	err := db.QueryRow("SELECT topic,message,clicked FROM urls WHERE shortUrl = $1", shortUrl).Scan(&topic, &message, &clicked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no such shorten url")
		}
		return nil, err
	}
	return []string{topic, message, fmt.Sprintf("%d", clicked)}, nil
}
