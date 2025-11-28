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

	got := namesFn("", Query{Params: []QueryParam{{Name: "names", Type: "[]string"}}})
	if got != "pq.Array(names)" {
		t.Fatalf("expected pq.Array wrapping for string slice params, got %q", got)
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

	got := dbPrefixFn("Exec", false, Table{GoName: "Customer", Fields: []Field{{GoName: "Tags", Type: "[]string"}}})
	if !strings.Contains(got, "pq.Array(c.Tags)") {
		t.Fatalf("expected pq.Array wrapping for table array fields, got %q", got)
	}
}

func TestCompositeArrayParamsNotWrapped(t *testing.T) {
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

	got := namesFn("", Query{Params: []QueryParam{{Name: "addresses", Type: "AddressArray"}}})
	if got != "addresses" {
		t.Fatalf("expected composite array params to be unwrapped, got %q", got)
	}
}

func TestCompositeArrayFieldsNotWrapped(t *testing.T) {
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

	got := dbPrefixFn("Exec", false, Table{GoName: "Customer", Fields: []Field{{GoName: "Addresses", Type: "AddressArray"}}})
	if strings.Contains(got, "pq.Array") {
		t.Fatalf("expected composite array fields to remain unwrapped, got %q", got)
	}
}
