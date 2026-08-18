package workspace

import (
	"testing"
	"errors"
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

func TestInvalidDirIsRejected(t *testing.T) {
	dirs := []string{
		"",
		"project/not_exists",
	}

	for _, tt := range dirs {
		t.Run("is rejected dir", func(t *testing.T) {
			_, err := OpenWorkspace(tt)
			if err == nil {
				t.Fatalf("want Err, got nil")
			}

			if !errors.Is(err, ErrOpenWorkspace) {
				t.Fatalf("want: %v, got: %v", ErrOpenWorkspace, err)
			}
		})
	}

}
