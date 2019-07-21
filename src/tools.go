package src

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

func writeTmpFile(buf []byte) []byte {
	if err := ioutil.WriteFile(tmpFile, buf, 0664); err != nil {
		panic(err)
		os.Exit(-1)
	}
	cmd := exec.Command("gofmt", "-w", tmpFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("run gofmt err: %s\n", err.Error())
	}
	f, _ := os.Open(tmpFile)
	defer f.Close()

	bts, _ := ioutil.ReadAll(f)
	return bts
}

//get absolute the output file path
func absFile() {
	if outputFile == nil || *outputFile == "" {
		return
	}

	var fileName = *outputFile
	if strings.HasPrefix(fileName, "~/") {
		fileName = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(fileName, "~"))
	}

	outputFile = &fileName
}

func parseTable(slice []byte) (table tables) {
	if len(slice) == 0 {
		return
	}
	var continueComment bool

	bts := bytes.Split(slice, []byte{'\n'})
	for _, line := range bts {
		line = trimLine(line)

		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, slashStar) {
			continueComment = true
			continue
		}

		if bytes.HasPrefix(line, starSlash) {
			continueComment = false
			continue
		}

		if bytes.HasPrefix(line, linLin) || continueComment || bytes.HasPrefix(line, set) || bytes.HasPrefix(line, drop) {
			continue
		}

		if bytes.HasPrefix(line, []byte("CREATE TABLE")) {
			name := trim(bytes.Split(line, space)[2])
			table.rawName = string(name)
			table.name = toString(title(name))
			table.field = nil
			table.index = nil
			continue
		}

		switch line[0] {
		case ')':
			break
		case '`':
			table.parseField(line)
		case 'P':
			table.parseKey(line)
		case 'I':
			table.parseIndex(line)
		case 'U':
			table.parseUniqueIndex(line)
		}
	}
	return
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

	name = bytes.Trim(name, "`")

	if i := bytes.Index(name, point); i > 0 {
		name = name[i:]
	}

	return name
}

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
			names[k][1] = 'D'
		case "Uuid":
			names[k] = []byte("UUID")
		case "Http":
			names[k] = []byte("HTTP")
		case "Url":
			names[k] = []byte("URL")
		}

		result = append(result, names[k]...)
	}
	return result
}

func marshalTable(table *tables) []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("package %s\n\n", *pkg))
	if len(table.imports) > 0 {
		buf.WriteString("import (\n")
		for _, v := range table.imports {
			buf.WriteString(fmt.Sprintf("\t\"%s\"", v))
		}
		buf.WriteString("\n)\n\n")
	}
	buf.WriteString(fmt.Sprintf("type %s struct {\n", table.name))

	if *sortField {
		sort.Strings(table.orderFields)
	}

	for _, k := range table.orderFields {

		//v: [field,comment,json_tag,gorm_tag...]
		v := table.field[k]
		buf.WriteString(fmt.Sprintf("\t%s\t%s", k, v[0]))

		if (jsonTag != nil && *jsonTag) || (gormTag != nil && *gormTag) {
			buf.Write([]byte{' ', '`'})
			tag(&buf, v, table.index, k)
			buf.WriteByte('`')
		}
		if string(v[1]) != "" {
			buf.WriteString(fmt.Sprintf(" //%s", v[1]))
		}
		buf.WriteByte('\n')
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')
	buf.WriteString(tableNameFunc(table.name, table.rawName))
	return buf.Bytes()
}

func tag(buf *bytes.Buffer, tags []string, index content, key string) {
	if jsonTag != nil && *jsonTag {
		buf.WriteString(fmt.Sprintf("json:\"%s\"", tags[2]))
	}
	if gormTag != nil && *gormTag {
		var data []string
		if jsonTag != nil && *jsonTag {
			buf.WriteByte(' ')
		}
		buf.WriteString("gorm:\"")
		if len(tags) > 3 {
			data = append(data, tags[3:]...)
		} else if index != nil && len(index[key]) > 0 {
			data = append(data, index[key]...)
		}
		if len(data) > 0 {
			buf.WriteString(strings.Join(data, ";"))
		}
		buf.WriteByte('"')
	}
}

func tableNameFunc(name, rawName string) string {
	return fmt.Sprintf(tableNameFuncStr, name, rawName)
}

func toString(str []byte) string {
	return *(*string)(unsafe.Pointer(&str))
}

func toByte(str string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := &reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
