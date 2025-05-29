package main

import (
	"testing"
)

func TestExecute_InvalidDSN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ожидалась паника, но её не было")
		}
	}()
	err := execute("127.0.0.1", "8080", "")
	if err == nil {
		t.Errorf("ожидалась ошибка, но err == nil")
	}
}
