package jsonstreamsplit

import (
	"bufio"
	"encoding/json"
	"io"
	"unicode/utf8"
)

const (
	open   = '{'
	close  = '}'
	quote  = '"'
	escape = '\\'
)

type Handler func(object []byte)

// Split takes an io.Reader and returns a slice of byte slices
// representing the json objects stored in the stream
func Split(reader io.Reader) (results [][]byte, err error) {
	err = SplitWithHandler(reader, func(object []byte) {
		results = append(results, object)
	})
	return
}

// Split takes an io.Reader and returns a slice of strings
// representing the json objects stored in the stream
func SplitString(reader io.Reader) (results []string, err error) {
	err = SplitWithHandler(reader, func(object []byte) {
		results = append(results, string(object))
	})
	return
}

// Split takes an io.Reader and returns a slice of raw json messages
// representing the json objects stored in the stream
// Note, invalid json strings are dropped without error
func SplitJsonRaw(reader io.Reader) (results []json.RawMessage, err error) {
	err = SplitWithHandler(reader, func(object []byte) {
		var o json.RawMessage
		_ = json.Unmarshal(object, &o)
		results = append(results, o)
	})
	return
}

// Split takes an io.Reader and a Handler function which receives slice a
// byte slice containing a single object from the stream.
func SplitWithHandler(reader io.Reader, handler Handler) error {
	r := bufio.NewReader(reader)
	var (
		currentMatch []rune
		isEscaped    bool
		isQuoted     bool
		depth        int64
	)
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch char {
		case open:
			if !isQuoted && !isEscaped {
				depth++
			}

		case close:
			if !isQuoted && !isEscaped {
				depth--
			}

		case quote:
			if !isEscaped {
				isQuoted = !isQuoted
			}
		}

		if char == escape {
			if !isEscaped {
				isEscaped = true
			} else {
				isEscaped = false
			}
		} else {
			isEscaped = false
		}

		currentMatch = append(currentMatch, char)

		if depth == 0 {
			handler(runeSliceToByteSlice(currentMatch))
			currentMatch = []rune{}
		}
	}
}

// a slightly more performant version of []byte(string(rs))
func runeSliceToByteSlice(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}
