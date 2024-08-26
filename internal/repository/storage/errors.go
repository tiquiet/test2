package storage

import "fmt"

var ErrorAlreadyExists = fmt.Errorf("key already exists")
var ErrorNotFound = fmt.Errorf("key not found")
