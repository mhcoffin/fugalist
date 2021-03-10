package fugalist

import (
	"fmt"
	"math/rand"
	"time"
)

// Note: characters are chosen to not require URL encoding.
const base64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-"
const sixBits = 0b111111
const sixtyBits = 0b111111111111111111111111111111111111111111111111111111111111

var ind []uint64

// We don't need to make Uniq() cryptographically secure, but not seeding at all
// results in a deterministic sequence, which will lead to collisions.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Map from char to index of chars in base64.
func init() {
	max := int32(0)
	for _, v := range base64 {
		if v > max {
			max = v
		}
	}
	ind = make([]uint64, max+1)
	for k, v := range base64 {
		ind[v] = uint64(k)
	}
}

/**
 * Returns 60 random bits encoded as a 10-character string.
 */
func Uniq() string {
	return base64encode(rand.Uint64() & sixtyBits)
}

func base64encode(a uint64) string {
	return fmt.Sprintf("%c%c%c%c%c%c%c%c%c%c",
		base64[(a>>54)&sixBits],
		base64[(a>>48)&sixBits],
		base64[(a>>42)&sixBits],
		base64[(a>>36)&sixBits],
		base64[(a>>30)&sixBits],
		base64[(a>>24)&sixBits],
		base64[(a>>18)&sixBits],
		base64[(a>>12)&sixBits],
		base64[(a>>6)&sixBits],
		base64[(a>>0)&sixBits],
	)
}

func base64decode(s string) uint64 {
	if len(s) != 10 {
		panic("bad string")
	}
	return ind[s[0]]<<54 | ind[s[1]]<<48 | ind[s[2]]<<42 | ind[s[3]]<<36 | ind[s[4]]<<30 | ind[s[5]]<<24 | ind[s[6]]<<18 | ind[s[7]]<<12 | ind[s[8]]<<6 | ind[s[9]]
}

/**
 * Returns the bit-wise XOR of the input strings.
 */
func Xor(s []string) string {
	accum := uint64(0)
	for _, v := range s {
		accum ^= base64decode(v)
	}
	return base64encode(accum)
}
