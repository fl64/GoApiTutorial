// model.go

package main

import (
	"database/sql"
	"fmt"
)

type license struct {
	ID        int    `json:"id"`
	Edition   string `json:"edition"`
	Devices   int    `json:"devices"`
	IssuedTo  string `json:"issued_to"`
	IssusedOn string `json:"issued_on"`
}

func (lic *license) getLicense(db *sql.DB) error {
	s := "SELECT edition, devices, issued_to, issued_on FROM licenses WHERE id=%d"
	statement := fmt.Sprintf(s, lic.ID)
	return db.QueryRow(statement).Scan(&lic.Edition, &lic.Devices, &lic.IssuedTo, &lic.IssusedOn)
}

func (lic *license) updateLicense(db *sql.DB) error {
	s := "UPDATE licenses SET edition='%s', devices=%d, issued_to='%s', issued_on='%s' WHERE id=%d"
	statement := fmt.Sprintf(s, lic.Edition, lic.Devices, lic.IssuedTo, lic.IssusedOn, lic.ID)
	_, err := db.Exec(statement)
	return err
}

func (lic *license) deleteLicense(db *sql.DB) error {
	s := "DELETE FROM licenses WHERE id=%d"
	statement := fmt.Sprintf(s, lic.ID)
	_, err := db.Exec(statement)
	return err
}

func (lic *license) createLicense(db *sql.DB) error {
	s := "INSERT INTO licenses(edition, devices, issued_to, issued_on) VALUES('%s', %d, '%s', '%s')"
	statement := fmt.Sprintf(s, lic.Edition, lic.Devices, lic.IssuedTo, lic.IssusedOn)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&lic.ID)

	if err != nil {
		return err
	}

	return nil
}

func getLicenses(db *sql.DB, start, count int) ([]license, error) {
	s := "SELECT id, edition, devices, issued_to, issued_on FROM licenses LIMIT %d OFFSET %d"
	statement := fmt.Sprintf(s, count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	licenses := []license{}

	for rows.Next() {
		var lic license
		if err := rows.Scan(&lic.ID, &lic.Edition, &lic.Devices, &lic.IssuedTo, &lic.IssusedOn); err != nil {
			return nil, err
		}
		licenses = append(licenses, lic)
	}

	return licenses, nil
}
