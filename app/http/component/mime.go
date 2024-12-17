package component

import (
	"bufio"
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"log"
	"os"
	"strings"
)

type MimeTypes interface {
	ContentType(extension string) string
}

type DefaultMimeTypes struct {
	mapping map[string]string
}

func (defaultMime *DefaultMimeTypes) ContentType(extension string) string {
	if extension == "" {
		return common.ContentTypeOctetStream
	}
	contentType, ok := defaultMime.mapping[extension]
	if ok {
		return contentType
	}

	return common.ContentTypeOctetStream
}

func MimeTypeFromMapping(mapping map[string]string) MimeTypes {
	log.Printf("MimeTypes from mapping, mapping=%v", mapping)
	return &DefaultMimeTypes{mapping: mapping}
}
func MimeTypeFromFile(filename string) (MimeTypes, error) {
	mapping := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(';')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		line = strings.TrimSuffix(line, ";")
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			log.Printf("line %s syntax error", line)
			continue
		}
		for i := 1; i < len(parts); i++ {
			mapping[strings.ToLower(parts[i])] = strings.ToLower(parts[0])
		}
	}

	return MimeTypeFromMapping(mapping), nil
}
