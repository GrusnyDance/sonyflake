package snowflake

import (
	"crypto/sha256"
	"encoding/binary"
)

func hash(input string, mask uint16) uint16 {
	hashStr := sha256.Sum256([]byte(input))
	h := binary.BigEndian.Uint16(hashStr[:])

	return h & mask
}
