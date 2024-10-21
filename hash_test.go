package snowflake

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCharset    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	testPodsNumber = 100
)

func Test_hash(t *testing.T) {
	t.Parallel()

	prefix := "origin-vod-"
	uniquePodNames := make(map[string]struct{})
	for len(uniquePodNames) < testPodsNumber {
		str := prefix + generateRandomString(10) + "-" + generateRandomString(5)
		uniquePodNames[str] = struct{}{}
	}

	hashMap := make(map[uint16]struct{})
	for name := range uniquePodNames {
		h := hash(name, MachineIDBitMask)
		assert.LessOrEqual(t, h, uint16(MachineIDBitMask))
		hashMap[h] = struct{}{}
	}

	fmt.Printf("collisions generated: %d\n", len(uniquePodNames)-len(hashMap))
}

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(testCharset))))
		if err != nil {
			panic(err)
		}
		result[i] = testCharset[num.Int64()]
	}
	return string(result)
}
