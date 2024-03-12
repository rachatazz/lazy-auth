## Lazy Auth
Authentication services are developed on the Gin framework.

## Installation
```bash
$ go get .
```

## Running the app
```bash
# development
$ go run main.go

# production build
$ go build -o ./dist/main

# production run
$ ./dist/main
```

## Reference documents
- HTTP framework - [Gin](https://gin-gonic.com/docs/)
- ORM - [GORM](https://gorm.io/docs/)
- Model binding and validation - [Validator](https://github.com/go-playground/validator)
- Configuration - [Viper](https://github.com/spf13/viper)