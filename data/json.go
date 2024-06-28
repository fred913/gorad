package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type JSONDatabase struct {
	flock  *sync.Mutex
	Path   string
	Indent bool
}

func NewJSONDatabase(path string, indent bool) *JSONDatabase {
	return &JSONDatabase{
		flock:  &sync.Mutex{},
		Path:   path,
		Indent: indent,
	}
}

func (j *JSONDatabase) Load(v interface{}, writeDefault bool) {
	j.flock.Lock()
	defer j.flock.Unlock()

	if _, err := os.Stat("config.json"); err != nil {
		// slog.Error("could not open config.json", "err", err)

		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			panic(fmt.Errorf("could not marshal default data: %w", err))
		} else {
			os.WriteFile("config.json", b, 0644)
		}
		panic("empty config.json created")
	}

	content, err := os.ReadFile(j.Path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, v)
	if err != nil {
		panic(err)
	}
}

func (j *JSONDatabase) Save(v interface{}) (err error) {
	j.flock.Lock()
	defer j.flock.Unlock()

	var content []byte
	if j.Indent {
		content, err = json.MarshalIndent(v, "", "    ")
	} else {
		content, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(j.Path, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
