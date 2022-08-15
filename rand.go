package main

import (
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const char = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01233456789"

// type uuid string

func Generate(size int, char string) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = char[seededRand.Intn(len(char))]
	}
	return string(b)
}

func String(lenght int) string {
	return Generate(lenght, char)
}
