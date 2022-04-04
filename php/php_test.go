package php

import (
	"testing"
)

func TestParse(t *testing.T) {
	// Test PHP file.
	src := "<?php \\Drupal::service(\"test\"); ?>"

	// Parse.
	doc, err := Parse([]byte(src))
	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}

	if len(doc.StaticCalls) >= 2 || len(doc.StaticCalls) <= 0 {
		t.Errorf("Invalid number of static calls found")
		return
	}

	staticCall := doc.StaticCalls[0]
	if staticCall.Class.Name != "Drupal" {
		t.Errorf("Invalid class name")
		return
	}

	if staticCall.Method.Name != "service" {
		t.Errorf("Invalid method name")
		return
	}

	if len(staticCall.Args) >= 2 || len(staticCall.Args) <= 0 {
		t.Errorf("Invalid number of args found")
		return
	}
}
