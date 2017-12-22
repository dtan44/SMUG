package handler

import (
	"testing"
)

func setupHelper() {
	secretKey = "correct"
}

func TestValidateKey(t *testing.T) {
	setupHelper()
	if !validateKey("correct") {
		t.Errorf("Failed Correct Validate Key")
	}

	if validateKey("wrong") {
		t.Error("Failed Wrong Validate Key")
	}
}
