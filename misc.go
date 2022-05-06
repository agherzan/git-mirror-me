// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"crypto/sha256"
	"encoding/hex"
)

func mask(what string) string {
	var masked string

	if len(what) != 0 {
		h := sha256.Sum256([]byte(what))
		masked = hex.EncodeToString(h[:])
	}

	return masked
}
