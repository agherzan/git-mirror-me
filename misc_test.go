// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"testing"
)

// TestMask tests the mask function.
func TestMask(t *testing.T) {
	t.Parallel()

	if m := mask(""); m != "" {
		t.Fatalf("unexpected output for \"\" input, got %s", m)
	}

	if m := mask("foo"); m != "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae" {
		t.Fatalf("unexpected output for \"foo\" input, got %s", m)
	}
}
