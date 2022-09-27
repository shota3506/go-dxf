package dxf

import (
	"strconv"
)

const (
	GroupCodeEntityType GroupCode = 0
)

type GroupCode int64

func (c GroupCode) String() string {
	return strconv.FormatInt(int64(c), 10)
}
