package spec

// JSON Schema validation with github.com/santhosh-tekuri/jsonschema (draft 2020-12).
//
// The schemas in spec/schema/ are the CONTRACT every language's output must honour.
// A "contract test" checks a real receipt (or saved order, or CLI error) against them.

import (
	"bytes"
	"encoding/json"
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// english prints the library's error messages in English.
var english = message.NewPrinter(language.English)

var (
	schemasMu sync.Mutex
	schemas   = map[string]*jsonschema.Schema{} // compiled schemas, by file name
)

// Schema compiles (once) and returns one of the shared schemas, e.g. "receipt.schema.json".
//
// The compiler reads the file from the spec folder, and resolves relative references like
// `"$ref": "receipt.schema.json"` in order.schema.json against that same folder.
func Schema(name string) (*jsonschema.Schema, error) {
	schemasMu.Lock()
	defer schemasMu.Unlock()

	if schema, ok := schemas[name]; ok {
		return schema, nil
	}
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	schema, err := jsonschema.NewCompiler().Compile(filepath.Join(dir, "schema", name))
	if err != nil {
		return nil, err
	}
	schemas[name] = schema
	return schema, nil
}

// SchemaErrors validates value against one of the shared schemas and returns one message
// per problem, sorted – an empty list means the value is valid. The second result is only set when
// the schema itself could not be loaded.
//
// value can be any Go value that encodes to JSON (a Receipt, a map, …): it is turned into
// JSON and back first, so the schema sees exactly what another program would see.
func SchemaErrors(name string, value any) ([]string, error) {
	schema, err := Schema(name)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	err = schema.Validate(instance)
	if err == nil {
		return nil, nil
	}
	var validationErr *jsonschema.ValidationError
	if !errors.As(err, &validationErr) {
		return nil, err
	}
	var problems []string
	collectProblems(validationErr, &problems)
	sort.Strings(problems)
	return problems, nil
}

// collectProblems walks the tree of validation errors and keeps the leaves – the actual
// problems ("/tax minimum: got -1, want 0") – not the nodes that only group them
// ("/tax validation failed"). Each problem reads "<where in the value> <what is wrong>".
func collectProblems(err *jsonschema.ValidationError, problems *[]string) {
	if len(err.Causes) == 0 {
		location := "/" + strings.Join(err.InstanceLocation, "/")
		*problems = append(*problems, location+" "+err.ErrorKind.LocalizedString(english))
		return
	}
	for _, cause := range err.Causes {
		collectProblems(cause, problems)
	}
}
