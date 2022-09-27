package dxf

import (
	"bufio"
	"fmt"
	"io"
)

// A Writer writes a drawing object using DXF format.
type Writer struct {
	w *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: bufio.NewWriter(w),
	}
}

func (w *Writer) Write(drawing *Drawing) error {
	if drawing.Header != nil {
		if err := w.writeHeader(drawing.Header); err != nil {
			return err
		}
	}

	if drawing.Classes != nil {
		if err := w.writeClasses(drawing.Classes); err != nil {
			return err
		}
	}
	// TODO: write tables

	if drawing.Blocks != nil {
		if err := w.writeBlocks(drawing.Blocks); err != nil {
			return err
		}
	}

	if drawing.Entities != nil {
		if err := w.writeEntities(drawing.Entities); err != nil {
			return err
		}
	}

	// TODO: write objects
	// TODO: write thumbnailimage

	if _, err := w.w.WriteString(" 0\nEOF"); err != nil {
		return err
	}

	if err := w.w.Flush(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) startSection() error {
	if _, err := w.w.WriteString(" 0\nSECTION\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) endSection() error {
	if _, err := w.w.WriteString(" 0\nENDSEC\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeEntity(entity *Entity) error {
	if _, err := w.w.WriteString(fmt.Sprintf(" 0\n%s\n", entity.Name)); err != nil {
		return err
	}
	for groupCode, value := range entity.Attributes {
		if _, err := w.w.WriteString(fmt.Sprintf(" %s\n%s\n", groupCode, value)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeHeader(header map[string]map[GroupCode]string) error {
	if err := w.startSection(); err != nil {
		return err
	}
	if _, err := w.w.WriteString(" 2\nHEADER\n"); err != nil {
		return err
	}

	for key, values := range header {
		// write header key
		if _, err := w.w.WriteString(fmt.Sprintf(" 0\n%s\n", key)); err != nil {
			return err
		}

		for groupCode, value := range values {
			if _, err := w.w.WriteString(fmt.Sprintf(" %s\n%s\n", groupCode, value)); err != nil {
				return err
			}
		}
	}

	if err := w.endSection(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeClasses(classes []*Class) error {
	if err := w.startSection(); err != nil {
		return err
	}
	if _, err := w.w.WriteString(" 2\nCLASSES\n"); err != nil {
		return err
	}

	for _, class := range classes {
		if err := w.writeClass(class); err != nil {
			return err
		}
	}

	if err := w.endSection(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeClass(class *Class) error {
	if _, err := w.w.WriteString(" 0\nCLASS\n"); err != nil {
		return err
	}

	for groupCode, value := range class.Attributes {
		if _, err := w.w.WriteString(fmt.Sprintf(" %s\n%s\n", groupCode, value)); err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) writeBlocks(blocks []*Block) error {
	if err := w.startSection(); err != nil {
		return err
	}
	if _, err := w.w.WriteString(" 2\nBLOCKS\n"); err != nil {
		return err
	}

	for _, block := range blocks {
		if err := w.writeBlock(block); err != nil {
			return err
		}
	}

	if err := w.endSection(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeBlock(block *Block) error {
	if _, err := w.w.WriteString(" 0\nBLOCK\n"); err != nil {
		return err
	}

	// write block attribute
	for groupCode, value := range block.Attributes {
		if _, err := w.w.WriteString(fmt.Sprintf(" %s\n%s\n", groupCode, value)); err != nil {
			return err
		}
	}

	for _, entity := range block.Entities {
		if err := w.writeEntity(entity); err != nil {
			return err
		}
	}

	if _, err := w.w.WriteString(" 0\nENDBLK\n"); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeEntities(entities []*Entity) error {
	if err := w.startSection(); err != nil {
		return err
	}
	if _, err := w.w.WriteString(" 2\nENTITIES\n"); err != nil {
		return err
	}

	for _, entity := range entities {
		if err := w.writeEntity(entity); err != nil {
			return err
		}
	}

	if err := w.endSection(); err != nil {
		return err
	}
	return nil
}
