package src

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
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
		panic(err)
		os.Exit(-2)
	}
	f, _ := os.Open(tmpFile)
	defer f.Close()

	bts, _ := ioutil.ReadAll(f)
	return bts
}

func parseFile(file *os.File) (table tables) {
	var continueComment bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := trimLine(scanner.Bytes())

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

		if bytes.HasPrefix(line, linlin) || continueComment || bytes.HasPrefix(line, set) || bytes.HasPrefix(line, drop) {
			continue
		}

		if bytes.HasPrefix(line, []byte("CREATE TABLE")) {
			table.name = string(title(trim(bytes.Split(line, space)[2])))
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
	line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
	line = bytes.TrimSuffix(line, []byte{'\n'})
	line = bytes.TrimSuffix(line, []byte{'\r'})
	return line
}

func trim(name []byte) []byte {
	if len(name) < 2 {
		return name
	}
	if bytes.Contains(name, point) {
		name = bytes.Split(name, point)[1]
	}
	if name[0] == '`' {
		name = name[1 : len(name)-1]
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

func parseTable(table tables) []byte {
	var buf bytes.Buffer
	buf.WriteString("package main\n\n")
	if len(table.imports) > 0 {
		buf.WriteString("import (\n")
		for _, v := range table.imports {
			buf.WriteString(fmt.Sprintf("\t\"%s\"", v))
		}
		buf.WriteString("\n)\n\n")
	}
	buf.WriteString(fmt.Sprintf("type %s struct {\n", table.name))

	var keys []string
	for k := range table.field {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := table.field[k]
		if len(v) > 2 {
			buf.WriteString(fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\t//%s\n", k, v[0], v[1], v[2]))
		} else {
			buf.WriteString(fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", k, v[0], v[1]))
		}
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')
	return buf.Bytes()
}
