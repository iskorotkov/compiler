package symbol

import (
	"crypto/md5"
)

type hasher struct {
	value int
}

func (h *hasher) Hash(value string) int {
	if h.value == 0 {
		hash := md5.Sum([]byte(value))
		h.value = int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])
	}

	return h.value
}
