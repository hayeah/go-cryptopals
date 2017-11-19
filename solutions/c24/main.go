package main

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"math"

	"github.com/pkg/errors"

	"github.com/hayeah/go-cryptopals/random"

	"github.com/hayeah/go-cryptopals"
)

func CrackOracle(o *Oracle) (key uint64, err error) {
	ptextbuf := make([]byte, 16)

	var ptext []byte
	var ctext []byte

	// Attempt to align chosen ptext with random prefix so the last 8 bytes are all zeroes
	for {
		var n uint64
		n, err = cryptopals.CryptoRandInt()
		if err != nil {
			return
		}

		ptextSize := n%8 + 8

		ptext = ptextbuf[:ptextSize+1]

		ctext, err = o.Encrypt(ptext)
		if err != nil {
			return
		}

		if len(ctext)%8 == 0 {
			break
		}
	}

	// spew.Dump("ptext", ptext)
	// spew.Dump("ctext", ctext)

	// last 8 bytes would've been xored with 0. Treat it as a number generated by Mersenne Twister
	nth := len(ctext)/8 - 1
	i64b := ctext[len(ctext)-8 : len(ctext)]
	nthUInt64 := binary.BigEndian.Uint64(i64b)

	// Now find a seed where the nth number of the sequence is `nthUInt64`
	for i := 0; i < math.MaxUint16; i++ {
		r := random.NewMersenneTwister()
		r.Seed(uint64(i))

		// generate `nth - 1` number
		for j := 0; j < nth-1; j++ {
			r.Next()
		}

		// nth's number
		if r.Next() == nthUInt64 {
			// win
			return uint64(i), nil
		}
	}

	return 0, errors.New("Failed to crack")
}

// Oracle an encryption oracle
type Oracle struct {
	// 16 bit seed
	key uint64
}

// Encrypt returns the cipher text
func (o *Oracle) Encrypt(ptext []byte) ([]byte, error) {
	prefix, err := cryptopals.CryptoRandPrefix(64)
	if err != nil {
		return nil, err
	}

	var outputW bytes.Buffer

	w := cipher.StreamWriter{
		S: random.NewStream(uint64(o.key)),
		W: &outputW,
	}

	w.Write(prefix)
	w.Write(ptext)

	if w.Err != nil {
		return nil, w.Err
	}

	return outputW.Bytes(), nil
}
