package loader

import (
	"testing"

	xo "github.com/xo/dbtpl/types"
)

func TestPostgresCompositeGoType(t *testing.T) {
	tests := []struct {
		name   string
		typ    xo.Type
		goType string
		zero   string
		mapper func(xo.Type, string, string, string) (string, string, error)
	}{
		{
			name:   "composite type",
			typ:    xo.Type{Type: "public.address"},
			goType: "Address",
			zero:   "Address{}",
			mapper: StdlibPostgresGoType,
		},
		{
			name:   "nullable composite type",
			typ:    xo.Type{Type: "public.address", Nullable: true},
			goType: "NullAddress",
			zero:   "NullAddress{}",
			mapper: StdlibPostgresGoType,
		},
		{
			name:   "array of composites stdlib",
			typ:    xo.Type{Type: "public.address", IsArray: true},
			goType: "[]Address",
			zero:   "nil",
			mapper: StdlibPostgresGoType,
		},
		{
			name:   "array of composites pq",
			typ:    xo.Type{Type: "public.address", IsArray: true},
			goType: "pq.GenericArray",
			zero:   "nil",
			mapper: PQPostgresGoType,
		},
	}
	for i, test := range tests {
		goType, zero, err := test.mapper(test.typ, "public", "int", "uint")
		if err != nil {
			t.Fatalf("test %d (%s) unexpected error: %v", i, test.name, err)
		}
		if goType != test.goType {
			t.Errorf("test %d (%s) expected go type %q, got %q", i, test.name, test.goType, goType)
		}
		if zero != test.zero {
			t.Errorf("test %d (%s) expected zero %q, got %q", i, test.name, test.zero, zero)
		}
	}
}
