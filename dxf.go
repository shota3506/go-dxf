package dxf

type Drawing struct {
	Header   map[string]map[GroupCode]string
	Classes  []*Class
	Blocks   []*Block
	Entities []*Entity
}

type Class struct {
	Attributes map[GroupCode]string
}

type Block struct {
	Entities   []*Entity
	Attributes map[GroupCode]string
}

type Entity struct {
	Name       string
	Attributes map[GroupCode]string
}
