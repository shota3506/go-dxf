package dxf_test

import (
	"os"
	"testing"

	"github.com/shota3506/go-dxf"
)

func TestReader_Read(t *testing.T) {
	file, err := os.Open("./testdata/cube.dxf")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	r := dxf.NewReader(file)
	drawing, err := r.Read()
	if err != nil {
		t.Fatal(err)
	}
	if drawing == nil {
		t.Fatal("expected drawing value not to be nil")
	}
}
