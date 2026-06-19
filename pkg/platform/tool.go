//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package platform

import (
	"strconv"
	"strings"
)

//=============================================================================

func addParameters(ids []uint) string {
	var sb strings.Builder

	for i, id := range ids {
		if i != 0 {
			sb.WriteString("&")
		}
		sb.WriteString("id=")
		sb.WriteString(strconv.FormatUint(uint64(id), 10))
	}

	return sb.String()
}

//=============================================================================
