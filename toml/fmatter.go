// Package fmatter is a TOML Front Matter parser, using the
// github.com/BurntSushi/toml package.
//
// DEPRECATED. This was a terrible idea.
package fmatter

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

// ReadFile read an entire file into memory, and calls Read which
// parses the front matter data and returns the remaining file
// contents.
func ReadFile(filename string, frontmatter interface{}) (content []byte, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Read(data, frontmatter)
}

// Read detects and parses the front matter data, and returns the
// remaining contents. If no front matter is found, the entire
// file contents are returned. frontmatter is passed to toml.Decode
func Read(data []byte, frontmatter interface{}) (content []byte, err error) {
	r := bytes.NewBuffer(data)

	// eat away starting whitespace
	var ch rune = ' '
	for unicode.IsSpace(ch) {
		ch, _, err = r.ReadRune()
		if err != nil {
			// file is just whitespace
			return []byte{}, nil
		}
	}
	r.UnreadRune()

	// check if first line is ---
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if strings.TrimSpace(line) != "---" {
		// no front matter, just content
		return data, nil
	}

	matterStart := len(data) - r.Len()
	matterEnd := matterStart

	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return nil, err
		}

		if strings.TrimSpace(line) == "---" {
			matterEnd = len(data) - r.Len() - len(line)
			content = data[matterEnd+len(line):]
			break
		}
	}

	_, err = toml.Decode(string(data[matterStart:matterEnd]), frontmatter)
	if err != nil {
		return nil, err
	}
	err = nil
	return
}
