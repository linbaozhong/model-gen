package validation

import "testing"

func TestValidator_Alpha(t *testing.T) {
	if !Alpha("adf97953") {
		t.Error("error", "---1")
	}
	if !Email("w@w.w") {
		t.Error("error", "---2")
	}
	if !AlphaNumeric("adf97953") {
		t.Error("error", "---3")
	}
}

func TestName(t *testing.T) {
	var name = ""
	b := Required(name)
	t.Log(b, name)
}
