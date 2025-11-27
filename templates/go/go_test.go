//go:build dbtpl

package gotpl

import (
	"context"
	"strings"
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

func TestQueryArrayParamsWrappedWithPQArray(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, xo.DriverKey, "postgres")
	ctx = context.WithValue(ctx, xo.OutKey, ".")

	funcs, err := NewFuncs(ctx)
	if err != nil {
		t.Fatalf("NewFuncs returned error: %v", err)
	}

	namesFn, ok := funcs["names"].(func(string, ...any) string)
	if !ok {
		t.Fatalf("names func missing")
	}

	got := namesFn("", Query{Params: []QueryParam{{Name: "addresses", Type: "[]AddressType"}}})
	if got != "pq.Array(addresses)" {
		t.Fatalf("expected pq.Array wrapping for composite slice params, got %q", got)
	}
}

func TestTableFieldsWrappedWithPQArray(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, xo.DriverKey, "postgres")
	ctx = context.WithValue(ctx, xo.OutKey, ".")
	ctx = context.WithValue(ctx, ShortsKey, map[string]string{})

	funcs, err := NewFuncs(ctx)
	if err != nil {
		t.Fatalf("NewFuncs returned error: %v", err)
	}

	dbPrefixFn, ok := funcs["db_prefix"].(func(string, bool, ...any) string)
	if !ok {
		t.Fatalf("db_prefix func missing")
	}

	got := dbPrefixFn("Exec", false, Table{GoName: "Customer", Fields: []Field{{GoName: "Addresses", Type: "[]AddressType"}}})
	if !strings.Contains(got, "pq.Array(c.Addresses)") {
		t.Fatalf("expected pq.Array wrapping for table array fields, got %q", got)
	}
}
