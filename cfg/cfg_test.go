package cfg

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/imdario/mergo"
	"gotest.tools/v3/assert"
)

func TestReadFile(t *testing.T) {
	var (
		content []byte
		data    interface{}
		err     error
		file    *File
		path    = "test/target.json"
	)

	if content, err = ioutil.ReadFile(path); err != nil {
		t.FailNow()
	}

	if err = json.Unmarshal(content, &data); err != nil {
		t.FailNow()
	}

	file, err = ReadFile(path)

	assert.NilError(t, err)
	assert.Equal(t, file.Path, path)
	assert.DeepEqual(t, file.Data, data)
}

func TestFileMergeWith(t *testing.T) {
	tests := map[string]struct {
		Src  string
		Res  string
		Opts []func(*mergo.Config)
	}{
		"1": {"test/1_src.json", "test/1_res.json", nil},
		"2": {"test/2_src.json", "test/2_res.json", nil},

		"3": {
			"test/3_src.json",
			"test/3_res.json",
			[]func(*mergo.Config){
				mergo.WithOverride,
			},
		},

		"4": {
			"test/4_src.json",
			"test/4_res.json",
			[]func(*mergo.Config){
				mergo.WithAppendSlice,
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			testMerge(t, tt.Src, tt.Res, tt.Opts...)
		})
	}
}

func testMerge(t *testing.T, src, res string, opts ...func(*mergo.Config)) {
	var (
		err     error
		file    *File
		source  *Source
		resjson []byte
		resdata map[string]interface{}
	)

	if file, err = ReadFile("test/target.json"); err != nil {
		t.Fatal("cannot open test/target.json")
	}

	if resjson, err = ioutil.ReadFile(res); err != nil {
		t.Fatalf("cannot read res: %s", res)
	}

	if err = json.Unmarshal(resjson, &resdata); err != nil {
		t.Fatal("cannot unmarshal res")
	}

	if source, err = ReadSourceFile(src, opts...); err != nil {
		t.Fatalf("cannot read src: %s", src)
	}

	err = file.MergeSource(source)

	assert.NilError(t, err)
	assert.DeepEqual(t, file.Data, resdata)
}

func TestPointerToMap(t *testing.T) {
	tests := map[string]struct {
		Pointer string
		Value   interface{}
		Map     map[string]interface{}
	}{
		"0": {
			"",
			nil,
			nil,
		},
		"1": {
			"a",
			true,
			map[string]interface{}{
				"a": true,
			},
		},
		"2": {
			"a/b",
			true,
			map[string]interface{}{
				"a": map[string]interface{}{
					"b": true,
				},
			},
		},
		"3": {
			"a/b/c",
			true,
			map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": true,
					},
				},
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			m := PointerToMap(tt.Pointer, tt.Value)
			assert.DeepEqual(t, m, tt.Map)
		})
	}
}
