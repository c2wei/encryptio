# Encryptio
This is a program to encrypt or decrypt files based on Golang.

### Build
```shell
go build -o encryptio main.go
```

### How to use?
Execute the binary in command line with the following 3 arguments.

cmd - encrypt or decrypt \
dirPath - //your/dir/path/* \
16CharKey - 16 characters key

```shell
encryptio {cmd} {dirPath} {16CharKey}

# examples
encryptio encrypt ./* 1234567812345678
encryptio decrypt ./* 1234567812345678
```
