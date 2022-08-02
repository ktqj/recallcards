package storage

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_readJsonFile(t *testing.T) {
	// Test reading from an empty file
	tmpFile, err := ioutil.TempFile("", "data.json")
	path := tmpFile.Name()
	if err != nil {
		t.Fatal(err)
	}

	d, err := readJsonFile[string](path)
	if err != nil {
		t.Fatal(err)
	}

	if len(d) != 0 {
		t.Fatalf("Expected empty map, got %v", d)
	}

	// Test malformed json
	tmpFile.Write([]byte("{====}"))
	d, err = readJsonFile[string](path)
	if !strings.HasPrefix(err.Error(), "invalid character") {
		t.Fatal(err)
	}

	if d != nil {
		t.Fatalf("Expected nil map, got %v", d)
	}

	// Test reading a mapped card
	content := []byte(`{"wąs":{"ID":"","Phrase":"wąs","Translation":"усы","RecallAttempts":[],"Bucket":0,"Created_at":"2022-08-02T15:35:54.316447+02:00"}}`)
	tmpFile.WriteAt(content, 0)
	d, err = readJsonFile[string](path)
	if err != nil {
		t.Fatal(err)
	}

	if len(d) != 1 {
		t.Fatalf("Expected map with 1 elem, got %v", d)
	}

	// Remove temp file to test creating file on read
	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}

	d, err = readJsonFile[string](path)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 0 {
		t.Fatalf("Expecting empty map, got %v", d)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}

	// Reading from non-existing directory
	d, err = readJsonFile[string]("/nonexisting/directory/to/read/data.json")
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if d != nil {
		t.Fatalf("Expecting nil map, got %v", d)
	}
}
