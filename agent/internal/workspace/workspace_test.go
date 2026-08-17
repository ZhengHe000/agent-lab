package workspace

import (
	"testing"
)

func TestLifeCycle(t *testing.T) {
	dir := t.TempDir()

	wok, err := OpenWorkspace(dir)
	if err != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	cls := wok.root.Close()
	if cls != nil {
		t.Fatalf("want nil, got: %v", err)
	}

	// marginal testing: wok.root.Close Does not support repeated calls
	cls = wok.root.Close()
	if cls != nil {
		t.Fatalf("want nil, got: %v", err)
	}
}
