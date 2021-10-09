// Provides utilities to help create and update project config files.
package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

type File struct {
	Path string
	Data map[string]interface{}
}

type unmarshaller func([]byte, interface{}) error

var unmarshallers = map[string]unmarshaller{
	".json": json.Unmarshal,
	".yaml": yaml.Unmarshal,
}

func ReadBytes(content []byte, p string, ts ...string) (*File, error) {
	var (
		fns  []unmarshaller
		err  error
		data map[string]interface{}
	)

	if len(ts) == 0 {
		ts = []string{path.Ext(p)}
	}

	fns = make([]unmarshaller, 0, len(ts))

	for _, t := range ts {
		if t[0] != '.' {
			t = "." + t
		}

		if fn := unmarshallers[t]; fn != nil {
			fns = append(fns, fn)
		}
	}

	if len(fns) == 0 {
		return nil, fmt.Errorf("unknown config file type: %q", ts)
	}

	for _, fn := range fns {
		if err = fn(content, &data); err == nil {
			return &File{p, data}, nil
		}
	}

	return nil, err
}

func ReadFile(p string, t ...string) (*File, error) {
	if content, err := ioutil.ReadFile(p); err != nil {
		return nil, err
	} else {
		return ReadBytes(content, p, t...)
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
	data []byte,
	p string,
	opts ...func(*mergo.Config),
) (*Source, error) {
	if file, err := ReadBytes(data, p); err != nil {
		return nil, err
	} else {
		return &Source{file, opts}, nil
	}
}

func ReadSourceFile(p string, opts ...func(*mergo.Config)) (*Source, error) {
	if file, err := ReadFile(p); err != nil {
		return nil, err
	} else {
		return &Source{file, opts}, nil
	}
}

func PointerToMap(ptr string, v interface{}) (mp map[string]interface{}) {
	ptr = strings.Trim(ptr, "/")

	if ptr != "" {
		mp = make(map[string]interface{}, 1)
		p := strings.SplitN(ptr, "/", 2)

		if len(p) == 2 {
			mp[p[0]] = PointerToMap(p[1], v)
		} else {
			mp[p[0]] = v
		}
	}

	return mp
}
