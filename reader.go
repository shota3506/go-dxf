package dxf

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidFormat    = errors.New("dxf: invalid format")
	ErrInvalidGroupCode = errors.New("dxf: invalid group code")
)

// A Reader reads a drawing object from a DXF file.
type Reader struct {
	r *bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// Read reads a drawing object from r.
func (r *Reader) Read() (*Drawing, error) {
	drawing := &Drawing{}

	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode != GroupCodeEntityType {
			return nil, ErrInvalidFormat
		}

		if value == "EOF" {
			break
		} else if value != "SECTION" {
			return nil, ErrInvalidFormat
		}

		// read next codes
		groupCode, value, err = r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode != GroupCode(2) {
			return nil, ErrInvalidFormat
		}
		switch value {
		case "HEADER":
			header, err := r.readHeader()
			if err != nil {
				return nil, err
			}
			drawing.Header = header
		case "CLASSES":
			classes, err := r.readClasses()
			if err != nil {
				return nil, err
			}
			drawing.Classes = classes
		case "BLOCKS":
			blocks, err := r.readBlocks()
			if err != nil {
				return nil, err
			}
			drawing.Blocks = blocks
		case "ENTITIES":
			entities, err := r.readEntities()
			if err != nil {
				return nil, err
			}
			drawing.Entities = entities
		default:
			// TODO: support TABLES, OBJECTS, THUMBNAILIMAGE

			// ignore section
			for {
				groupCode, value, err := r.readCodes()
				if err != nil {
					if err == io.EOF {
						return nil, ErrInvalidFormat
					}
					return nil, err
				}

				if groupCode == GroupCodeEntityType && value == "ENDSEC" {
					break
				}
			}
		}
	}

	return drawing, nil
}

func (r *Reader) readCodes() (GroupCode, string, error) {
	code, err := r.r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if code == "" {
				return 0, "", io.EOF
			}
			return 0, "", ErrInvalidFormat
		}
		return 0, "", err
	}
	value, err := r.r.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			return 0, "", err
		}
		if value == "" {
			return 0, "", ErrInvalidFormat
		}
	}

	codeInt, err := strconv.ParseInt(strings.TrimSpace(code), 10, 64)
	if err != nil {
		return 0, "", ErrInvalidGroupCode
	}

	return GroupCode(codeInt), strings.TrimSpace(value), nil
}

func (r *Reader) readHeader() (map[string]map[GroupCode]string, error) {
	header := map[string]map[GroupCode]string{}

	key := ""
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode == GroupCodeEntityType && value == "ENDSEC" {
			break
		}

		if groupCode == GroupCode(9) {
			key = value
			header[key] = map[GroupCode]string{}
			continue
		}

		if key == "" {
			return nil, ErrInvalidFormat
		}

		header[key][groupCode] = value
	}
	return header, nil
}

// TODO
func (r *Reader) readClasses() ([]*Class, error) {
	classes := []*Class{}

	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode == GroupCodeEntityType {
			if value == "ENDSEC" {
				break
			} else if value == "CLASS" {
				class, err := r.readClass()
				if err != nil {
					return nil, err
				}
				classes = append(classes, class)
				continue
			}
		}
		return nil, ErrInvalidFormat
	}
	return classes, nil
}

func (r *Reader) readClass() (*Class, error) {
	class := &Class{
		Attributes: map[GroupCode]string{},
	}

	// GroupCode 1, 2, 3, 90, 91, 280, 281 are required
	for i := 0; i < 7; i++ {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		class.Attributes[groupCode] = value
	}
	return class, nil
}

func (r *Reader) readBlocks() ([]*Block, error) {
	blocks := []*Block{}

outer:
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode == GroupCodeEntityType {
			switch value {
			case "ENDSEC":
				break outer
			case "BLOCK":
				block, err := r.readBlock()
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, block)
				continue
			}
		}
		return nil, ErrInvalidFormat
	}
	return blocks, nil
}

func (r *Reader) readBlock() (*Block, error) {
	block := &Block{
		Entities:   []*Entity{},
		Attributes: map[GroupCode]string{},
	}

	var entity *Entity
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode == GroupCodeEntityType {
			if entity != nil {
				block.Entities = append(block.Entities, entity)
			}
			if value == "ENDBLK" {
				break
			}

			entity = &Entity{
				Name:       value,
				Attributes: map[GroupCode]string{},
			}
			continue
		}

		if entity != nil { // entity attibutes
			entity.Attributes[groupCode] = value
		} else { // block attributes
			block.Attributes[groupCode] = value
		}
	}
	return block, nil
}

func (r *Reader) readEntities() ([]*Entity, error) {
	entities := []*Entity{}

	var entity *Entity
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}
			return nil, err
		}

		if groupCode == GroupCodeEntityType {
			if entity != nil {
				entities = append(entities, entity)
			}
			if value == "ENDSEC" {
				break
			}

			entity = &Entity{
				Name:       value,
				Attributes: map[GroupCode]string{},
			}
			continue
		}

		if entity == nil {
			return nil, ErrInvalidFormat
		}

		entity.Attributes[groupCode] = value
	}
	return entities, nil
}

// TODO
func (r *Reader) readObjects() error {
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return ErrInvalidFormat
			}
			return err
		}

		if groupCode == GroupCodeEntityType && value == "ENDSEC" {
			break
		}
	}
	return nil
}

// TODO
func (r *Reader) readThumbnailimage() error {
	for {
		groupCode, value, err := r.readCodes()
		if err != nil {
			if err == io.EOF {
				return ErrInvalidFormat
			}
			return err
		}

		if groupCode == GroupCodeEntityType && value == "ENDSEC" {
			break
		}
	}
	return nil
}
