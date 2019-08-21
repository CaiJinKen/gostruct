# gostruct
This is an open source tool convert mysql table structure to golang`s struct and written by pure golang.
If you want get a beautify struct, make sure the gofmt has been installed and in the $PATH.

Install:
```bash
go get github.com/CaiJinKen/gostruct
```

Usage:

* get struct from sql file
```bash
gostruct -i users.sql -o users.go
```

* get struct from db connection
```bash
gostruct -d `user:password@tcp(host:port)/db_name` -t users -o ./models/users.go
```

This tool also can:
* generate `json` tag(default) and `gorm` tag
- print the struct(default)
* sort struct fields
     


Help:
```bash
gostruct --help
```
or
```bash
gostruct -h
```

