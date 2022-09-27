package dxf_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/shota3506/go-dxf"
)

func TestWriter_Write(t *testing.T) {
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

	buffer := bytes.NewBuffer(nil)
	w := dxf.NewWriter(buffer)
	if err := w.Write(drawing); err != nil {
		t.Fatal(err)
	}
	// TODO: assert buffer
}
