package models

import (
	"fmt"
	"testing"
)

func Test_getPlaintextAndMarkIndices(t *testing.T) {
	rebuilt, err := rebuild(getPlaintextAndMarkIndices("<maRk>Foo</mArk> bar #<marK>bat#</maRk> #<marK>!az#</maRk>"))
	if err != nil {
		t.Fatalf("Could not rebuild string: %s", err)
	}
	if rebuilt != "<mark>Foo</mark> bar <mark>#bat#</mark> #<mark>!az#</mark>" {
		t.Fatalf("Unexpected reconstituted string: %s", rebuilt)
	}
	rebuilt, err = rebuild(getPlaintextAndMarkIndices("foobar <maRk>Foo</mArk> bar #<marK>bat#</maRk> #<marK>!az#</maRk>"))
	if err != nil {
		t.Fatalf("Could not rebuild string: %s", err)
	}
	if rebuilt != "foobar <mark>Foo</mark> bar <mark>#bat#</mark> #<mark>!az#</mark>" {
		t.Fatalf("Unexpected reconstituted string: %s", rebuilt)
	}
	rebuilt, err = rebuild(getPlaintextAndMarkIndices("foobar no marks here"))
	if err != nil {
		t.Fatalf("Could not rebuild string: %s", err)
	}
	if rebuilt != "foobar no marks here" {
		t.Fatalf("Unexpected reconstituted string: %s", rebuilt)
	}
}

func rebuild(plain string, marksArr []int) (string, error) {
	marks := make(map[int]bool)
	for _, mark := range marksArr {
		marks[mark] = true
	}
	rebuilt := ""
	var isOpen bool
	for i, ch := range plain {
		if _, ok := marks[i]; ok {
			if !isOpen {
				rebuilt += "<mark>"
			} else {
				rebuilt += "</mark>"
			}
			isOpen = !isOpen
		}
		rebuilt += string(ch)
	}
	i := len(plain)
	if _, ok := marks[i]; ok {
		if !isOpen {
			return "", fmt.Errorf("Imbalanced mark tags")
		} else {
			rebuilt += "</mark>"
		}
	}
	return rebuilt, nil
}
