package main

import (
	bytesUtils "bytes"
	"encoding/json"
	"io/ioutil"

	xml2json "github.com/basgys/goxml2json"
	yaml2json "github.com/ghodss/yaml"
)

//ConvertToJSON converts file of a given type to JSON
func ConvertToJSON(fileType string, filePath string) json.RawMessage {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	switch fileType {
	case "YAML":
		return unmarshalJSONFromYAML(bytes)
	case "XML":
		return unmarshalXMLFromYAML(bytes)
	default:
		return nil
	}
}

func unmarshalJSONFromYAML(bytes []byte) json.RawMessage {
	bytes, err := yaml2json.YAMLToJSON(bytes)
	if err != nil {
		panic(err)
	}

	var target json.RawMessage
	err = json.Unmarshal(bytes, &target)
	if err != nil {
		panic(err)
	}
	return target
}

func unmarshalXMLFromYAML(bytes []byte) json.RawMessage {
	r := bytesUtils.NewReader(bytes)
	buf, err := xml2json.Convert(r)
	if err != nil {
		panic(err)
	}

	var target json.RawMessage
	err = json.Unmarshal(buf.Bytes(), &target)
	if err != nil {
		panic(err)
	}
	return target
}
