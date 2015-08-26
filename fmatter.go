// Package fmatter is a simple YAML Front Matter parser, using the
// gopkg.in/yaml.v2 package.
package fmatter

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
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
// file contents are returned. For details on the frontmatter
// parameter, please see the gopkg.in/yaml.v2 package.
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

	yamlStart := len(data) - r.Len()
	yamlEnd := yamlStart

	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return data, nil
			}
			return nil, err
		}

		if strings.TrimSpace(line) == "---" {
			yamlEnd = len(data) - r.Len()
			break
		}
	}

	err = yaml.Unmarshal(data[yamlStart:yamlEnd], frontmatter)
	if err != nil {
		return nil, err
	}
	content = data[yamlEnd:]
	err = nil
	return
}
