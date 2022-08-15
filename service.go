package main

import (
	// "bytes"
	// "time"
	//"crypto/internal/randutil"

	"database/sql"

	//	"fmt"
	"log"
	"sync"
	"time"
	//	"golang.org/x/crypto/bcrypt"
	//"log"
	// "math/rand"
)

type Service struct {
	secretservice secrets
}

func NewService(conn *sql.DB) Service {
	vault := NewVault(conn)
	return Service{
		secretservice: vault,
	}
}

func (service *Service) CreateSecret(secret vault) (vault, error) {
	secret, err := service.secretservice.Create(secret)
	if err != nil {
		return vault{}, err
	}

	duration, _ := time.ParseDuration(secret.duration)
	// execute in future where the duration of a secret expires and gets deleted
	newtimer := time.NewTimer(duration)
	go func() {
		<-newtimer.C
		err := service.secretservice.Delete(secret.uuid)
		if err != nil {
			log.Println(err)
		}
	}()

	return secret, nil
}

func (service *Service) FindSecret(id string) (vault, error) {
	var wg sync.WaitGroup
	secret, err := service.secretservice.Uuid(id)
	if err != nil {
		return vault{}, err
	}
	wg.Add(1)
	go func() {
		err := service.secretservice.Delete(id)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	wg.Wait()
	return secret, err
}

func (service *Service) DeleteSecret(id string) error {
	err := service.secretservice.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// func (service *Service) Encrypt(key, text []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	b := base64.StdEncoding.EncodeToString(text)
// 	ciphertex := make([]byte, aes.BlockSize+len(b))
// 	iv := ciphertex[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return nil, err
// 	}
// 	cfb := cipher.NewCFBEncrypter(block, iv)
// 	cfb.XORKeyStream(ciphertex[aes.BlockSize:], []byte(b))
// 	return ciphertex, nil
// }

// func (service *Service) Decrypt(key, text []byte) ([]byte, error) {

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(text) < aes.BlockSize {
// 		return nil, errors.New("ciphertext too short")
// 	}
// 	iv := text[:aes.BlockSize]
// 	text = text[aes.BlockSize:]
// 	cfb := cipher.NewCFBDecrypter(block, iv)
// 	cfb.XORKeyStream(text, text)
// 	data, err := base64.StdEncoding.DecodeString(string(text))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
