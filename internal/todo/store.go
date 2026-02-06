package todo

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type ToDo struct {
	ID          uuid.UUID `json:"ID"`
	Name        string    `json:"name"`
	IsDone      bool      `json:"isDone"`
	Description string    `json:"desc"`
}

type Store struct {
	Path string
}

func NewStore(path string) Store {
	return Store{Path: path}
}

func New(name string, desc string) ToDo {
	return ToDo{
		ID:          uuid.New(),
		Name:        name,
		IsDone:      false,
		Description: desc,
	}
}

func (s Store) Load() ([]ToDo, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ToDo{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []ToDo{}, nil
	}

	var todos []ToDo
	if err := json.Unmarshal(data, &todos); err != nil {
		return []ToDo{}, nil
	}

	return todos, nil
}

func (s Store) Save(todos []ToDo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0644)
}
