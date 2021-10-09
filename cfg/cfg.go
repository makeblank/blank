// Provides utilities to help create and update project config files.
package cfg

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/imdario/mergo"
)

type File struct {
	Path string
	Data map[string]interface{}
}

func ReadBytes(content []byte, path string) (*File, error) {
	var data map[string]interface{}

	// TODO: read YAML files
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	} else {
		return &File{path, data}, nil
	}
}

func ReadFile(path string) (*File, error) {
	if content, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		return ReadBytes(content, path)
	}
}

func (f *File) MergeSource(srcs ...*Source) error {
	for _, s := range srcs {
		if err := mergo.Merge(&f.Data, s.File.Data, s.Options...); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) MergeMap(
	m map[string]interface{},
	opts ...func(*mergo.Config),
) error {
	if err := mergo.Merge(&f.Data, m, opts...); err != nil {
		return err
	}
	return nil
}

type Source struct {
	File    *File
	Options []func(*mergo.Config)
}

func ReadSourceBytes(
	content []byte,
	path string,
	opts ...func(*mergo.Config),
) (*Source, error) {
	if file, err := ReadBytes(content, path); err != nil {
		return nil, err
	} else {
		return &Source{file, opts}, nil
	}
}

func ReadSourceFile(path string, opts ...func(*mergo.Config)) (*Source, error) {
	if file, err := ReadFile(path); err != nil {
		return nil, err
	} else {
		return &Source{file, opts}, nil
	}
}

func PointerToMap(pointer string, v interface{}) (mp map[string]interface{}) {
	pointer = strings.Trim(pointer, "/")

	if pointer != "" {
		mp = make(map[string]interface{}, 1)
		p := strings.SplitN(pointer, "/", 2)

		if len(p) == 2 {
			mp[p[0]] = PointerToMap(p[1], v)
		} else {
			mp[p[0]] = v
		}
	}

	return mp
}
