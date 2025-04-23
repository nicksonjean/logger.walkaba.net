package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		r.Uint32(), r.Uint32()&0xffff, (r.Uint32()&0x0fff)|0x4000,
		(r.Uint32()&0x3fff)|0x8000, r.Int63n(0xffffffffffff))
	return uuid
}
