package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var schemalookup map[string]*gojsonschema.Schema

func igluval(iglu string, document interface{}) bool {
	iglu = strings.Replace(iglu, "iglu:", os.Getenv("IGLU_PATH"), 1)
	loader := gojsonschema.NewGoLoader(document)
	schema, err := schemaget(iglu)
	if err == nil {
		result, error := schema.Validate(loader)
		if result.Valid() && error == nil {
			//fmt.Println("Valid Schema")
			return true
		}
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			fmt.Printf("- %s\n", err)
		}
		return false
	}
	fmt.Println(err)
	return false
}

func schemaget(iglu string) (*gojsonschema.Schema, error) {
	var err error
	if val, ok := schemalookup[iglu]; ok {
		return val, nil
	}

	schemaLoader := gojsonschema.NewReferenceLoader(iglu)
	schemalookup[iglu], err = gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		delete(schemalookup, iglu)
		return schemalookup[iglu], err
	}
	return schemalookup[iglu], nil
}
