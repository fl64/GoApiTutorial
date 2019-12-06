// main_test.go

package main

import (
  "fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "password", "gm_licenses")
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM licenses")
	a.DB.Exec("ALTER TABLE licenses AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
 CREATE TABLE IF NOT EXISTS licenses
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    edition VARCHAR(50) NOT NULL, 
    devices INT NOT NULL,
    issued_to VARCHAR(50) NOT NULL,
    issued_on VARCHAR(50) NOT NULL
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/licenses", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentLicense(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/license/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "License not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'License not found'. Got '%s'", m["error"])
	}
}

func TestCreateLicense(t *testing.T) {
	clearTable()

	payload := []byte(`{"edition":"Enterprise","devices":100,"issued_on":"2019-12-06","issued_to":"GMTEST"}`)

	req, _ := http.NewRequest("POST", "/api/v1/license", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["edition"] != "Enterprise" {
		t.Errorf("Expected license edition to be 'Enterprise'. Got '%v'", m["name"])
	}

	if m["issued_to"] != "GMTEST" {
		t.Errorf("Expected license devices to be 'GMTEST'. Got '%v'", m["devices"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func addLicenses(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
    statement := fmt.Sprintf("INSERT INTO licenses(edition,devices,issued_on,issued_to) VALUES('%s', %d, '%s','%s')", "Enterprise", ((i+1) * 10), "2019-01-01", "License " + strconv.Itoa(i+1) )
		a.DB.Exec(statement)
	}
}

func TestGetLicense(t *testing.T) {
	clearTable()
	addLicenses(1)

	req, _ := http.NewRequest("GET", "/api/v1/license/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateLicense(t *testing.T) {
	clearTable()
	addLicenses(1)

	req, _ := http.NewRequest("GET", "/api/v1/license/1", nil)
	response := executeRequest(req)
	var originalLicense map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalLicense)

	payload := []byte(`{"edition":"Enterprise","devices":88,"issued_on":"2019-12-06","issued_to":"GMTEST1"}`)

	req, _ = http.NewRequest("PUT", "/api/v1/license/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalLicense["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalLicense["id"], m["id"])
	}

	if m["devices"] == originalLicense["devices"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalLicense["devices"], m["devices"], m["devices"])
	}

	if m["issued_to"] == originalLicense["issued_to"] {
		t.Errorf("Expected the age to change from '%v' to '%v'. Got '%v'", originalLicense["issued_to"], m["issued_to"], m["issued_to"])
	}
}

func TestDeleteLicense(t *testing.T) {
	clearTable()
	addLicenses(1)

	req, _ := http.NewRequest("GET", "/api/v1/license/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/api/v1/license/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/api/v1/license/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
