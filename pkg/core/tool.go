//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
)

//=============================================================================

func WriteJsonInZip(w *zip.Writer, filename string, object any) error {
	data,err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	return WriteFileInZip(w, filename, data)
}

//=============================================================================

func WriteFileInZip(w *zip.Writer, filename string, data []byte) error {
	writer, err := w.Create(filename)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	_, err = io.Copy(writer, buffer)
	return err
}

//=============================================================================
