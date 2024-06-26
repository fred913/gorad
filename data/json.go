package data

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONDatabase struct {
	DataLock sync.RWMutex
	flock    sync.Mutex
	Path     string
	Indent   bool
}

func (j *JSONDatabase) LoadInto(v interface{}) {
	j.DataLock.Lock()
	defer j.DataLock.Unlock()
	j.flock.Lock()
	defer j.flock.Unlock()

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
	j.DataLock.RLock()
	defer j.DataLock.RUnlock()
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
