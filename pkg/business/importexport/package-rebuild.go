//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package importexport

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/platform"
)

//=============================================================================
//===
//=== Package rebuilding
//===
//=============================================================================

func RebuildPackage(data []byte) (*InMemoryPackage, error) {
	pack := &InMemoryPackage{}

	contentMap,err := unzipToMemory(data)
	if err != nil {
		return nil, err
	}

	err = rebuildStruct(contentMap, MetadataFileName, &pack.Metadata)
	if err != nil {
		return nil, err
	}

	err = rebuildStruct(contentMap, InventoryFileName, &pack.Data)
	if err != nil {
		return nil, err
	}

	portEd := &platform.ExportedData{}
	err = rebuildStruct(contentMap, PortfolioFileName, &portEd)
	if err != nil {
		return nil, err
	}

	tsMap := buildTradingSystemsMap(pack.Data)
	addPortfolioData(tsMap, portEd)
	addStorageData(tsMap, contentMap)

	return pack, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func unzipToMemory(zipData []byte) (map[string][]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	extractedFiles := make(map[string][]byte)

	for _, file := range reader.File {
		// If it's a directory entry, we skip it (maps only store file data)
		if file.FileInfo().IsDir() {
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(fileReader)
		_=fileReader.Close()
		if err != nil {
			return nil, err
		}

		extractedFiles[file.Name] = content
	}

	return extractedFiles, nil
}

//=============================================================================

func rebuildStruct(contentMap map[string][]byte, filename string, out any) error {
	data,okm := contentMap[filename]
	if !okm {
		return req.NewBadRequestError("missing file: %v", filename)
	}

	err := json.Unmarshal(data, &out)
	if err != nil {
		return req.NewBadRequestError("bad json format in file: %v", filename)
	}

	return nil
}

//=============================================================================

func buildTradingSystemsMap(data *ExportedData) map[uint]*TradingSystem{
	res := make(map[uint]*TradingSystem)

	for _, ts := range data.TradingSystems {
		res[ts.Id] = ts
	}

	return res
}

//=============================================================================

func addPortfolioData(tsMap map[uint]*TradingSystem, data *platform.ExportedData) {
	for _, pts := range data.TradingSystems {
		ts,ok := tsMap[pts.Id]
		if ok {
			ts.portfolioData = pts.JsonData
		} else {
			slog.Warn("addPortfolioData: Skipping portfolio data for unknown trading system: %v", pts.Id)
		}
	}
}

//=============================================================================

func addStorageData(tsMap map[uint]*TradingSystem, contentMap map[string][]byte) {
	for name, data := range contentMap {
		if strings.HasSuffix(name, ".zip") {
			name = strings.TrimSuffix(name, ".zip")
			id, err := strconv.Atoi(name)
			if err != nil {
				slog.Warn("addStorageData: Skipping non-numeric trading system: %v", name)
			} else {
				ts,ok := tsMap[uint(id)]
				if !ok {
					slog.Warn("addStorageData: Skipping storage data for unknown trading system: %v", id)
				} else {
					ts.storageData = data
				}
			}
		}
	}
}

//=============================================================================
