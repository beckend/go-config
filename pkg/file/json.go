package file

import (
	bytes "bytes"
	jsonOriginal "encoding/json"
	fmt "fmt"
	time "time"

	jsonpatch "github.com/evanphx/json-patch"
)

var (
	JSONTimeRFC3399Layout = "2006-01-02T15:04:05.999999999Z0700"
	nilTime               = (time.Time{}).UnixNano()
)

// usage: `json:"time"`
type JSONTimeRFC3399 struct {
	time.Time
}

func (self *JSONTimeRFC3399) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(JSONTimeRFC3399Layout, s)
	}
	self.Time = t
	return nil
}

func (self *JSONTimeRFC3399) MarshalJSON() ([]byte, error) {
	if self.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", self.Time.Format(JSONTimeRFC3399Layout))), nil
}

func (ct *JSONTimeRFC3399) IsSet() bool {
	return ct.UnixNano() != nilTime
}

func JSONMerge(input [][]byte) ([]byte, error) {
	returned := []byte(`{}`)

	for _, toMerge := range input {
		combined, err := jsonpatch.MergeMergePatches(returned, toMerge)
		if err != nil {
			return nil, err
		}

		returned = combined
	}

	return returned, nil
}

// create a json reader https://github.com/gookit/config/blob/master/read.go which can handle env
// https://stackoverflow.com/questions/53152852/how-to-iterate-recursively-through-mapstringinterface

func JSONPrettify(input []byte) (*bytes.Buffer, error) {
	var prettyJSON bytes.Buffer
	err := jsonOriginal.Indent(&prettyJSON, input, "", "\t")
	if err != nil {
		return nil, err
	}

	return &prettyJSON, nil
}

func JSONPrettifyNoError(input []byte) *bytes.Buffer {
	var prettyJSON bytes.Buffer
	result, err := JSONPrettify(input)
	if err != nil {
		return &prettyJSON
	}

	return result
}

type JSONWriteInterfaceOptions struct {
	Input    interface{}
	PathFile string
	Pretty   bool
}

func JSONWriteInterface(options *JSONWriteInterfaceOptions) ([]byte, error) {
	contentBytes, err := json.Marshal(options.Input)
	if err != nil {
		return nil, err
	}

	if options.Pretty {
		contentBytesPretty, err := JSONPrettify(contentBytes)
		if err != nil {
			return nil, err
		}

		contentBytes = contentBytesPretty.Bytes()
	}

	WriteFile(&WriteFileOptions{
		ContentsBytes: contentBytes,
		PathFile:      options.PathFile,
	})

	return contentBytes, nil
}

func JSONGenericMapToBytes(input map[string]interface{}) ([]byte, error) {
	returned, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return returned, nil
}
