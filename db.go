package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Vault struct {
	db *sql.DB
}

// Create(vault) (vault, error)
// 		Find(id int) (vault, error)
// 		Delete(id int) error

func NewVault(conn *sql.DB) Vault {
	return Vault{
		db: conn,
	}
}

func (v Vault) Create(vt vault) (vault, error) {

	sqlStatement := `
  INSERT INTO vault (secret,duration,uuid) 
  VALUES ($1,$2,$3)
  RETURNING *;
  `
	err := v.db.QueryRow(sqlStatement, vt.secret, vt.duration, vt.uuid).Scan(
		&vt.id,
		&vt.secret,
		&vt.duration,
		&vt.uuid,
	)
	return vt, err
}

func (v Vault) Find(id int) (vault, error) {

	sqlStatement := `
 SELECT * FROM vault
  WHERE vault.id = $1
  `
	var vt vault
	err := v.db.QueryRow(sqlStatement, id).Scan(
		&vt.id,
		&vt.secret,
		&vt.duration,
		&vt.uuid,
	)
	return vt, err
}

func (v Vault) Uuid(uuid string) (vault, error) {

	sqlStatement := `
 SELECT * FROM vault
  WHERE vault.uuid = $1
  `
	var vt vault
	err := v.db.QueryRow(sqlStatement, uuid).Scan(
		&vt.id,
		&vt.secret,
		&vt.duration,
		&vt.uuid,
	)
	return vt, err
}
func (v Vault) Delete(id string) error {
	sqlStatement := `DELETE FROM vault
  WHERE vault.uuid = $1
  `
	_, err := v.db.Exec(sqlStatement, id)
	return err
}
