package file

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"example.com/recallcards/pkg/cards"
)

func Test_fileReader(t *testing.T) {
	// Test reading from an empty file
	tmpFile, err := ioutil.TempFile("", "cards.json")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpFile.Name()
	defer os.Remove(path)

	fileReader := readJsonFile[cards.Cards]
	d, err := fileReader(path)
	if err != nil {
		t.Fatal(err)
	}

	emptryStorage := make(cards.Cards, 0)
	if reflect.DeepEqual(d, emptryStorage) {
		t.Fatalf("Expected %v, got %v", emptryStorage, d)
	}

	// Test malformed json
	_, err = tmpFile.Write([]byte("{====}"))
	if err != nil {
		t.Fatal(err)
	}

	d, err = fileReader(path)
	if !strings.HasPrefix(err.Error(), "invalid character") {
		t.Fatal(err)
	}

	if reflect.DeepEqual(d, emptryStorage) {
		t.Fatalf("Expected %v, got %v", emptryStorage, d)
	}

	// Test reading a mapped card
	content := []byte(`[{"ID":1,"Phrase":"wąs","Translation":"усы","Bucket":0,"Created_at":"2022-09-02T15:35:54.316447+02:00"}]`)
	_, err = tmpFile.WriteAt(content, 0)
	if err != nil {
		t.Fatal(err)
	}

	d, err = fileReader(path)
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

	d, err = fileReader(path)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(d, emptryStorage) {
		t.Fatalf("Expecting empty map, got %v", d)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}

	// Reading from non-existing directory
	d, err = fileReader("/nonexisting/directory/to/read/data.json")
	if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if reflect.DeepEqual(d, emptryStorage) {
		t.Fatalf("Expecting nil map, got %v", d)
	}
}
