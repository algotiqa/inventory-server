//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"slices"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/inventory-server/pkg/core"
	"github.com/algotiqa/inventory-server/pkg/core/process/agentscanner"
	"github.com/algotiqa/inventory-server/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

const ConfigAgentName = "config/agent.yaml"
const ConfigKeyName   = "config/agent.key"
const ConfigCertName  = "config/agent.crt"
const LogName         = "log/"

//=============================================================================

const agentCertValidityDays = 9999
var certDomain = pkix.Name{
	Country           : []string{"EU"},
	Province          : []string{"Italy"},
	Locality          : []string{"Rome"},
	Organization      : []string{"Algotiqa"},
	OrganizationalUnit: []string{"AlgotiqaAgent"},
	CommonName        : "*",
}

//=============================================================================

type AgentCertificateBundle struct {
	Key  []byte
	Cert []byte
}

//=============================================================================

func GetAgentProfiles(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.AgentProfile, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	list, err := db.GetAgentProfiles(tx, filter, offset, limit)

	if err != nil {
		return nil, err
	}

	return list, nil
}

//=============================================================================

func GetAgentProfileById(tx *gorm.DB, c *auth.Context, id uint) (*AgentProfileExt, error) {
	c.Log.Info("GetAgentProfileById: Getting an agent profile", "id", id)

	ap, err := getAgentProfile(tx, c, id, "GetAgentProfileById")
	if err != nil {
		return nil, err
	}

	//--- Get trading systems

	filter := make(map[string]any)
	filter["agent_profile_id"] = id
	tss,err := db.GetTradingSystemsFull(tx, filter, 0, 5000)
	if err != nil {
		c.Log.Error("GetAgentProfileById: Could not retrieve trading systems", "error", err.Error())
		return nil, err
	}

	//--- Just remove some sensitive information

	ap.SslKey  = nil
	ap.SslCert = nil

	//--- Put all together

	ape := AgentProfileExt{
		AgentProfile  : *ap,
		TradingSystems: tss,
	}

	return &ape, nil
}

//=============================================================================

func AddAgentProfile(tx *gorm.DB, c *auth.Context, aps *AgentProfileSpec) (*db.AgentProfile, error) {
	c.Log.Info("AddAgentProfile: Adding a new agent profile", "name", aps.Name)

	//TODO: Validate input

	var ap db.AgentProfile
	ap.Username      = c.Session.Username
	ap.Name          = aps.Name
	ap.Host          = aps.Host
	ap.Port          = aps.Port
	ap.ScanInterval  = aps.ScanInterval
	ap.ScanFolder    = aps.ScanFolder
	ap.FileExtension = aps.FileExtension
	ap.HostType      = aps.HostType

	bundle,err := generateAgentCertificates(certDomain, aps.Host)
	if err != nil {
		return nil, err
	}

	ap.SslKey  = bundle.Key
	ap.SslCert = bundle.Cert

	err = db.AddAgentProfile(tx, &ap)

	if err != nil {
		c.Log.Error("AddAgentProfile: Could not add a new agent profile", "error", err.Error())
		return nil, req.NewServerErrorByError(err)
	}

	c.Log.Info("AddAgentProfile: Agent profile added", "id", ap.Id, "name", ap.Name)
	return &ap, nil
}

//=============================================================================

func UpdateAgentProfile(tx *gorm.DB, c *auth.Context, id uint, aps *AgentProfileSpec) (*db.AgentProfile, error) {
	c.Log.Info("UpdateAgentProfile: Updating an agent profile", "id", id, "name", aps.Name)

	ap, err := getAgentProfile(tx, c, id, "UpdateAgentProfile")
	if err != nil {
		return nil, err
	}

	//TODO: Validate input

	ap.Name          = aps.Name
	ap.Host          = aps.Host
	ap.Port          = aps.Port
	ap.ScanInterval  = aps.ScanInterval
	ap.ScanFolder    = aps.ScanFolder
	ap.FileExtension = aps.FileExtension
	ap.HostType      = aps.HostType

	err = db.UpdateAgentProfile(tx, ap)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	c.Log.Info("UpdateAgentProfile: Agent profile updated", "id", ap.Id, "name", ap.Name)
	return ap, nil
}

//=============================================================================

func DeleteAgentProfile(tx *gorm.DB, c *auth.Context, id uint) (string, error) {
	c.Log.Info("DeleteAgentProfile: Deleting agent profile", "id", id)

	bp, err := getAgentProfile(tx, c, id, "DeleteAgentProfile")
	if err != nil {
		return "", err
	}

	//--- Check if there are references (not the efficient way, but...)

	filter := map[string]any{}
	filter["agent_profile_id"] = id

	tss,err := db.GetTradingSystems(tx, filter, 0, 10)
	if err != nil {
		return "", err
	}
	if len(*tss) > 0 {
		return DeleteStatusTradingSystems, err
	}

	//--- Proper delete

	err = db.DeleteAgentProfile(tx, id)
	if err != nil {
		c.Log.Error("DeleteAgentProfile: Cannot delete agent profile", "id", id, "error", err.Error())
		return "", req.NewServerErrorByError(err)
	}

	c.Log.Info("DeleteAgentProfile: Agent profile deleted", "id", id, "name", bp.Name)
	return DeleteStatusOk, nil
}

//=============================================================================

func GetExternalRefs(tx *gorm.DB, c *auth.Context, id uint) ([]string, error) {
	c.Log.Info("GetExternalRefs: Getting external refs from profile", "id", id)

	ap, err := getAgentProfile(tx, c, id, "GetExternalRefs")
	if err != nil {
		return nil, err
	}

	list, err := callAgentToGetExternalRefs(c, ap)
	if err != nil {
		return nil, err
	}

	dbRefs,err := getExternalRefsForProfile(tx, c, id)
	if err != nil {
		return nil, err
	}

	res := []string{}

	for _, xref := range list {
		if _,ok := dbRefs[xref]; !ok {
			res = append(res, xref)
		}
	}

	slices.Sort(res)

	c.Log.Info("GetExternalRefs: Got new list of external refs", "id", id, "size", len(res))
	return res, nil
}

//=============================================================================

func GetAgentPackage(tx *gorm.DB, c *auth.Context, id uint) ([]byte, error) {
	c.Log.Info("GetAgentPackage: Building package from profile", "id", id)

	ap, err := getAgentProfile(tx, c, id, "GetAgentPackage")
	if err != nil {
		return nil, err
	}

	return buildAgentPackage(ap)
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func callAgentToGetExternalRefs(c *auth.Context, ap *db.AgentProfile) ([]string, error) {
	client := agentscanner.CreateClient(ap.SslCert, ap.SslKey)
	if client == nil {
		return nil, req.NewServerError("Cannot create client for agent: %v", ap.Id)
	}

	var list []string
	err := req.DoGet(client, ap.RemoteUrl() + agentscanner.UrlTradingSystems, &list, "")
	if err != nil {
		c.Log.Error("callAgentToGetExternalRefs: Agent raised an error", "id", ap.Id, "error", err.Error())
		return nil, req.NewServiceUnavailableError("Agent raised an error : " + err.Error())
	}

	return list, nil
}

//=============================================================================

func getExternalRefsForProfile(tx *gorm.DB, c *auth.Context, id uint) (map[string]bool, error) {
	filter := map[string]any{}
	filter["username"]         = c.Session.Username
	filter["agent_profile_id"] = id
	list, err := db.GetTradingSystems(tx, filter, 0, 5000)

	if err != nil {
		return nil, err
	}

	result := map[string]bool{}
	for _,ts := range *list {
		result[ts.ExternalRef] = true
	}

	return result, nil
}

//=============================================================================

func getAgentProfile(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.AgentProfile, error) {
	ap, err := db.GetAgentProfileById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve agent profile", "error", err.Error())
		return nil, err
	}

	if ap == nil {
		c.Log.Error(function+": Agent profile was not found", "id", id)
		return nil, req.NewNotFoundError("Agent profile was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if ap.Username != c.Session.Username {
			c.Log.Error(function+": Agent profile not owned by user", "id", id)
			return nil, req.NewForbiddenError("Agent profile is not owned by user: %v", id)
		}
	}

	return ap, nil
}

//=============================================================================

func generateAgentCertificates(subject pkix.Name, agentIP string) (*AgentCertificateBundle, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, req.NewServerError("Cannot generate agent key: %v", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, req.NewServerError("Cannot generate serial number: %v", err)
	}

	var ipAddresses []net.IP

	ip := net.ParseIP(agentIP)
	if ip == nil {
		return nil, req.NewBadRequestError("Invalid agent host: %s", agentIP)
	}
	ipAddresses = []net.IP{ip}

	now := time.Now()

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               subject,
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, agentCertValidityDays),
		SignatureAlgorithm:    x509.SHA256WithRSA,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IPAddresses:           ipAddresses,
		BasicConstraintsValid: true,
	}

	certDer, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, req.NewServerError("Cannot create agent certificate: %v", err)
	}

	agentCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	agentKey  := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return &AgentCertificateBundle{
		Key : agentKey,
		Cert: agentCert,
	}, nil
}

//=============================================================================

func buildAgentPackage(ap *db.AgentProfile) ([]byte, error) {
	buf := new(bytes.Buffer)
	w   := zip.NewWriter(buf)

	err := writeConfig(w, ap)
	if err == nil {
		err = core.WriteFileInZip(w, ConfigKeyName, ap.SslKey)
		if err == nil {
			err = core.WriteFileInZip(w, ConfigCertName, ap.SslCert)
			if err == nil {
				err = writeAgentBinary(w, ap)
				if err == nil {
					err = writeLogDir(w)
					if err == nil {
						//--- We cannot defer this, because buf will remain empty
						_=w.Close()
						return buf.Bytes(), nil
					}
				}
			}
		}
	}

	_=w.Close()
	return nil, err
}

//=============================================================================

func writeConfig(w *zip.Writer, ap *db.AgentProfile) error {
	config := fmt.Sprintf(
		"application:\n"+
		"  bindAddress: :%v\n"+
		"  production: true\n"+
		"  debug: false\n" +
		"scan:\n" +
		"  dir: %v\n" +
		"  extension: %v\n",
		ap.Port, ap.ScanFolder, ap.FileExtension)

	return core.WriteFileInZip(w, ConfigAgentName, []byte(config))
}

//=============================================================================

func writeAgentBinary(w *zip.Writer, ap *db.AgentProfile) error {
	localFolder := "agent"
	agentFile   := "agent"

	if ap.HostType == db.HostTypeWindows {
		agentFile += ".exe"
	}

	data,err := os.ReadFile(localFolder +"/" +agentFile)
	if err != nil {
		return err
	}

	return core.WriteFileInZip(w, agentFile, data)
}

//=============================================================================

func writeLogDir(w *zip.Writer) error {
	fh := &zip.FileHeader{
		Name  : LogName,
		Method: zip.Store,
	}
	fh.SetMode(os.ModeDir | 0755)
	_, err := w.CreateHeader(fh)
	return err
}

//=============================================================================
