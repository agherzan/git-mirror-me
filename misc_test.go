// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"testing"
)

// TestMask tests the mask function.
func TestMask(t *testing.T) {
	if m := mask(""); m != "" {
		t.Fatalf("unexpected output for \"\" input, got %s", m)
	}
	if m := mask("foo"); m != "acbd18db4cc2f85cedef654fccc4a4d8" {
		t.Fatalf("unexpected output for \"foo\" input, got %s", m)
	}
}
