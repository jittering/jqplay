package jmes

import (
	"context"
	"encoding/json"
	"io"

	"github.com/jmespath/go-jmespath"
	"github.com/pkg/errors"
)

func Eval(ctx context.Context, input string, search string, w io.Writer) error {
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return err
	}
	result, err := jmespath.Search(search, data)
	if err != nil {
		return errors.Wrap(err, "failed to search jmespath")
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return errors.Wrap(err, "failed to encode json")
	}
	return nil
}
