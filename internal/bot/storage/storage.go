package storage

import "net/http"

type Storage interface {
	Save(file *http.Response) error
}
