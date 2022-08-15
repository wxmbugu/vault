package main

// import (
// 	"time"
// )

type (
	// Holds the secrets for a specific duration
	vault struct {
		id       int
		secret   string
		duration string
		uuid     string
	}
	secrets interface {
		Create(vault) (vault, error)
		Find(id int) (vault, error)
		Uuid(uuid string) (vault, error)
		Delete(id string) error
	}
)
