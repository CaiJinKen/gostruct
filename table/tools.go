package table

import (
	"bytes"
)

func title(name []byte) []byte {
	names := bytes.Split(name, []byte{'_'})
	result := make([]byte, 0)
	for k := range names {
		char := names[k][0]
		if 'a' <= char && char <= 'z' {
			names[k][0] = char - ('a' - 'A')
		}
		switch string(names[k]) {
		case "Id":
			names[k] = []byte("ID")
		case "Uuid":
			names[k] = []byte("UUID")
		case "Http":
			names[k] = []byte("HTTP")
		case "Https":
			names[k] = []byte("HTTPS")
		case "Url":
			names[k] = []byte("URL")
		case "Html":
			names[k] = []byte("HTML")
		}

		result = append(result, names[k]...)
	}
	return result
}

func trimLine(line []byte) []byte {
	line = bytes.TrimSpace(line)
	lt := len(line)
	for lt > 0 && (line[lt-1] == '\n' || line[lt-1] == '\r') {
		lt--
		if lt == 0 {
			line = line[:0]
		} else {
			line = line[:lt-1]
		}
	}
	return line
}

func trim(name []byte) []byte {
	if len(name) < 2 {
		return name
	}

	if i := bytes.Index(name, []byte{'.'}); i > 0 {
		name = name[i:]
	}

	name = bytes.Trim(name, "`")

	return name
}
