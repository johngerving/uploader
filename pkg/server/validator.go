package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Validator is an object that can be validated.
type Validator interface {
	// Valid checks the object and returns
	// any problems. If len(problems) == 0
	// then the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}