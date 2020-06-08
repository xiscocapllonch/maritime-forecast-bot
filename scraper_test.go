package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFormatDate(t *testing.T) {
	expectedOut := "domingo, 7 junio 2020 a las 20:00 hora oficial"

	date := Date("2020-06-07T20:00:00")
	formattedDate := date.formatDate()

	if formattedDate != expectedOut {
		t.Errorf("Expected Output: %v\nBut got: %v", expectedOut, formattedDate)
	}
}

func TestGetXML(t *testing.T) {
	expectedWarning := "Posibilidad de aguaceros y tormentas muy fuertes."

	result, err := getMockXML(t)

	if err != nil {
		t.Errorf("error getting XML: %v", err)
	} else {
		if result.Warning.Text != expectedWarning {
			t.Errorf(
				"Warning Text should be: %v\nBut got: %v",
				expectedWarning,
				result.Warning.Text,
			)
		}
	}
}

func getMockXML(t *testing.T) (Result, error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("test_data/test.xml")
		if err != nil {
			t.Errorf("error opening test file: %v", err)
		}
		defer file.Close()

		bFile, err := ioutil.ReadAll(file)

		_, err = w.Write(bFile)
		if err != nil {
			t.Errorf("error writing test file: %v", err)
		}
	}))
	defer ts.Close()

	return getXML(ts.URL)
}
