//go:build dbtpl

package gotpl

import (
	"context"
	"testing"

	xo "github.com/xo/dbtpl/types"
)

func TestArrayModeDefaultsToPQ(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, xo.DriverKey, "postgres")
	ctx = context.WithValue(ctx, xo.OutKey, ".")

	funcs, err := NewFuncs(ctx)
	if err != nil {
		t.Fatalf("NewFuncs returned error: %v", err)
	}

	importsFn, ok := funcs["imports"].(func() []PackageImport)
	if !ok {
		t.Fatalf("imports func missing")
	}

	got := importsFn()
	foundPQ := false
	for _, imp := range got {
		if imp.Pkg == "\"github.com/lib/pq\"" {
			foundPQ = true
			break
		}
	}
	if !foundPQ {
		t.Fatalf("expected pq import to be injected for default array mode")
	}
}
