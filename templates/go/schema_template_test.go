package gotpl

import (
	"os"
	"strings"
	"testing"
)

func TestCompositeNullValueUsesJSONMarshal(t *testing.T) {
	data, err := os.ReadFile("schema.dbtpl.go.tpl")
	if err != nil {
		t.Fatalf("failed to read schema template: %v", err)
	}
	if !strings.Contains(string(data), "return json.Marshal({{ short $nullName }}.{{ $c.GoName }})") {
		t.Fatalf("null composite Value should marshal to JSON for driver compatibility")
	}
}
