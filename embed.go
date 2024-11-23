package embed

import (
	_ "embed"
)

//go:embed schema.sql
var DBSchema string