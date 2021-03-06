package pkcs7

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaddedReader(t *testing.T) {
	is := assert.New(t)

	ts := []struct {
		bsize int
		in    string
		out   string
		err   error
	}{
		{8, "abcd123", "abcd123\x01", nil},

		{8, "abcd", "abcd\x04\x04\x04\x04", nil},

		{8, "abcd123\x01", "abcd123\x01\x08\x08\x08\x08\x08\x08\x08\x08", nil},
		{8, "abcd123\x08", "abcd123\x08\x08\x08\x08\x08\x08\x08\x08\x08", nil},
		{8, "abcd123\x09", "abcd123\x09\x08\x08\x08\x08\x08\x08\x08\x08", nil},

		{16, "abcd1234abcd", "abcd1234abcd\x04\x04\x04\x04", nil},
	}

	for _, t := range ts {
		r := NewPaddedReader(bytes.NewReader([]byte(t.in)), t.bsize)

		out, err := ioutil.ReadAll(r)
		is.NoError(err)

		is.Equal([]byte(t.out), out)
	}
}

func TestPad(t *testing.T) {
	is := assert.New(t)

	ts := []struct {
		bsize int
		in    string
		out   string
		err   error
	}{
		{8, "", "\x08\x08\x08\x08\x08\x08\x08\x08", nil},
		{8, "abcd123", "abcd123\x01", nil},

		{8, "abcd", "abcd\x04\x04\x04\x04", nil},
		{16, "abcd1234abcd", "abcd1234abcd\x04\x04\x04\x04", nil},

		{8, "abcd123\x01", "abcd123\x01\x08\x08\x08\x08\x08\x08\x08\x08", nil},
		{8, "abcd123\x08", "abcd123\x08\x08\x08\x08\x08\x08\x08\x08\x08", nil},
		{8, "abcd123\x09", "abcd123\x09", nil},
	}

	for _, t := range ts {
		var src []byte
		if t.in != "" {
			src = []byte(t.in)
		}
		out := Pad(src, t.bsize)
		is.Equal([]byte(t.out), out)
	}
}
