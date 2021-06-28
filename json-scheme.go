package helpers

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func IsValidObject(body interface{}, objScheme map[string]interface{}) error {
	obj, scheme := gojsonschema.NewGoLoader(body), gojsonschema.NewGoLoader(objScheme)

	if result, err := gojsonschema.Validate(scheme, obj); err != nil {
		return err
	} else if result.Valid() {
		return nil
	} else {
		e := ""
		for _, desc := range result.Errors() {
			e += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(e)
	}
}
