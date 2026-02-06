package todo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreSaveLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "todo.json")
	store := NewStore(path)

	items := []ToDo{
		New("Test item", "desc"),
	}

	if err := store.Save(items); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected 1 item, got %d", len(loaded))
	}

	if loaded[0].Name != items[0].Name {
		t.Fatalf("name mismatch: %q vs %q", loaded[0].Name, items[0].Name)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
