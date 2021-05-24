package file

import (
	io "io"
	fs "io/fs"

	gotoml "github.com/pelletier/go-toml"
)

func TOMLFileReaderCallbackToJSON(getFileCallback func() (fs.File, error)) ([]byte, error) {
	file, err := getFileCallback()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	res, err := gotoml.LoadReader(file)
	if err != nil {
		return nil, err
	}

	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLFileReaderToJSON(input fs.File, doClose bool) ([]byte, error) {
	res, err := gotoml.LoadReader(input)
	if err != nil {
		return nil, err
	}

	if doClose {
		defer input.Close()
	}
	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLReaderToJSON(input io.Reader) ([]byte, error) {
	res, err := gotoml.LoadReader(input)
	if err != nil {
		return nil, err
	}

	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLStringToJSON(input string) ([]byte, error) {
	res, err := gotoml.Load(input)
	if err != nil {
		return nil, err
	}

	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLBytesToJSON(input []byte) ([]byte, error) {
	res, err := gotoml.LoadBytes(input)
	if err != nil {
		return nil, err
	}

	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLFileToJSON(input string) ([]byte, error) {
	res, err := gotoml.LoadFile(input)
	if err != nil {
		return nil, err
	}

	return JSONGenericMapToBytes(res.ToMap())
}

func TOMLFilesToMergedJSON(pathFiles []string) ([]byte, error) {
	var toBeMerged [][]byte

	for _, pathFile := range pathFiles {
		bytesToMerge, err := TOMLFileToJSON(pathFile)
		if err != nil {
			return nil, err
		}

		result, err := JSONMerge([][]byte{bytesToMerge})
		if err != nil {
			return nil, err
		}

		toBeMerged = append(toBeMerged, result)
	}

	return JSONMerge(toBeMerged)
}

func TOMLBytesToMergedJSON(bytesSlices [][]byte) ([]byte, error) {
	var toBeMerged [][]byte

	for _, slice := range bytesSlices {
		result, err := JSONMerge([][]byte{slice})
		if err != nil {
			return nil, err
		}

		toBeMerged = append(toBeMerged, result)
	}

	return JSONMerge(toBeMerged)
}
