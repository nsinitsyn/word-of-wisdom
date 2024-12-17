package challenge

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math"
)

const (
	inputSize  = 16
	targetSize = 16
	nonceSize  = 8
)

type Challenge struct {
	Input  []byte
	Target []byte
}

func NewRandomChallenge(complexity uint8) *Challenge {
	var c = &Challenge{}
	c.populateByRandom(complexity)
	return c
}

func NewChallenge(bytes []byte) *Challenge {
	var c = &Challenge{}
	c.Input = bytes[:inputSize]
	c.Target = bytes[inputSize:]
	return c
}

func (c *Challenge) populateByRandom(leadingZeros uint8) {
	c.Input = make([]byte, inputSize)
	c.Target = make([]byte, targetSize)

	_, _ = rand.Read(c.Input)
	_, _ = rand.Read(c.Target)

	copy(c.Target[:leadingZeros], make([]byte, leadingZeros))
}

func (c *Challenge) Solve() []byte {
	nonce := make([]byte, nonceSize)

	for i := uint64(0); i < math.MaxUint64; i++ {
		binary.BigEndian.PutUint64(nonce, i)
		if c.VerifyNonce(nonce) {
			// log.Printf("found solution on iteration %d\n", i)
			return nonce
		}
	}

	return nil
}

func (c *Challenge) VerifyNonce(nonce []byte) bool {
	h := sha256.New()
	h.Write(c.Input)
	h.Write(nonce)
	hash := h.Sum(nil)
	return bytes.Compare(hash, c.Target) < 0
}
