//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package importexport

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/db"
	"github.com/algotiqa/inventory-server/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

const (
	MetadataFileName  = "metadata.json"
	InventoryFileName = "inventory.json"
	PortfolioFileName = "portfolio.json"
)

//=============================================================================

func CollectTradingSystemsData(tx *gorm.DB, username string, ids []uint) (*TradingSystemsData, error) {
	tss,err := db.GetTradingSystemsById(tx, username, ids)
	if err != nil {
		return nil, err
	}
	if len(*tss) != len(ids) {
		return nil, req.NewBadRequestError("Some trading systems were not found or not accessible")
	}

	filter := map[string]any{}
	filter["username"] = username

	dps,err := db.GetDataProductsFull(tx, filter,0, 50000)
	if err != nil {
		return nil, err
	}

	bps,err := db.GetBrokerProductsFull(tx, filter,0, 50000)
	if err != nil {
		return nil, err
	}

	sess,err := db.GetTradingSessions(tx, username)
	if err != nil {
		return nil, err
	}

	aps,err := db.GetAgentProfiles(tx, filter,0, 50000)
	if err != nil {
		return nil, err
	}

	tsd := NewTradingSystemsData(*tss, *dps, *bps, *sess, *aps)
	if tsd == nil {
		return nil, req.NewServerError("Database constraints mismatch! Username='%v'", username)
	}

	tsd.Ids = ids
	return tsd,nil
}

//=============================================================================

func WriteMetadata(zipWriter *zip.Writer) error {
	md := &Metadata{
		Version: &Version{
			Major: 1,
			Minor: 0,
		},
		ExportDate: time.Now(),
	}

	return writeFile(zipWriter, MetadataFileName, md)
}

//=============================================================================

func WriteData(zipWriter *zip.Writer, data *TradingSystemsData) error {
	ed := &ExportedData{}

	//--- Write data products

	for _, dp := range data.DataProducts {
		dpc := NewDataProduct(dp)
		ed.DataProducts = append(ed.DataProducts, dpc)
	}

	//--- Write broker products

	for _, bp := range data.BrokerProducts {
		bpc := NewBrokerProduct(bp)
		ed.BrokerProducts = append(ed.BrokerProducts, bpc)
	}

	//--- Write trading sessions

	for _, s := range data.TradingSessions {
		sc := NewTradingSession(s)
		ed.TradingSessions = append(ed.TradingSessions, sc)
	}

	//--- Write agent profiles

	for _, ap := range data.AgentProfiles {
		apc := NewAgentProfile(ap)
		ed.AgentProfiles = append(ed.AgentProfiles, apc)
	}

	//--- Write trading systems

	for _, ts := range data.TradingSystems {
		tsc := NewTradingSystem(&ts)
		ed.TradingSystems = append(ed.TradingSystems, tsc)
	}

	return writeFile(zipWriter, InventoryFileName, ed)
}

//=============================================================================

func WriteTradingSystemsData(c *auth.Context, zipWriter *zip.Writer, ids []uint) error {
	ed,err := platform.ExportTradingSystemsFromPortfolio(c, ids)
	if err != nil {
		return err
	}

	dirs,err := platform.ExportTradingSystemsFromStorage(c, ids)
	if err != nil {
		return err
	}

	dirsMap,err := extractZipInMemory(dirs)
	if err != nil {
		return err
	}

	for file, content := range dirsMap {
		err = writeData(zipWriter, file, content)
		if err != nil {
			return err
		}
	}

	return writeFile(zipWriter, PortfolioFileName, ed)
}

//=============================================================================
//===
//=== Private function
//===
//=============================================================================

func extractZipInMemory(zipData []byte) (map[string][]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	extractedFiles := make(map[string][]byte)

	for _, file := range reader.File {
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

func writeFile(zipWriter *zip.Writer, filename string, object any) error {
	data,err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	return writeData(zipWriter, filename, data)
}

//=============================================================================

func writeData(zipWriter *zip.Writer, filename string, data []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)
	_, err = io.Copy(writer, buffer)
	return err
}

//=============================================================================
