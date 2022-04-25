// SPDX-FileCopyrightText: Andrei Gherzan <andrei@gherzan.com>
//
// SPDX-License-Identifier: MIT

package mirror

import (
	"crypto/md5"
	"encoding/hex"
)

func mask(what string) string {
	var masked string
	if len(what) != 0 {
		h := md5.Sum([]byte(what))
		masked = hex.EncodeToString(h[:])
	}
	return masked
}
