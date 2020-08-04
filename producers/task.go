package producers

import (
	"bytes"
	"encoding/json"
	"io"
)

type Task struct {
	Body io.Reader
}

func NewTask(body interface{}) (*Task, error) {
	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(body); err != nil {
		return nil, err
	}
	return &Task{Body: buff}, nil
}
