package utils

import (
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	mspIdentity "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("Setup")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// FabricSetup struct
type FabricSetup struct {
	AdminUser         string
	OrdererOrgName    string
	ConfigFileName    string
	Secret            []byte
	IdentityTypeUser  string
	Sdk               *fabsdk.FabricSDK
	RegistrarUsername string
	RegistrarPassword string
	ChannelID         string
	Identity          msp.Identity
	Targets           []string
}

// Init reads config file, setup client, CA
func (hfc *FabricSetup) Init() {
	//adding logger for outut
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(backendLeveled, backendFormatter)
	logger.Info("================ Creating New SDK Instance ================")

	//initializing SDK
	var config = config.FromFile(hfc.ConfigFileName)
	var err error
	hfc.Sdk, err = fabsdk.New(config)
	if err != nil {
		logger.Infof("Unable to create new instance of SDk: %s\n", err)
	}

	//clean up user data from previous runs
	configBackend, err := hfc.Sdk.Config()
	if err != nil {
		logger.Fatal(err)
	}

	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(configBackend)
	identityConfig, err := mspIdentity.ConfigFromBackend(configBackend)
	if err != nil {
		logger.Fatal(err)
	}

	keyStorePath := cryptoSuiteConfig.KeyStorePath()
	credentialStorePath := identityConfig.CredentialStorePath()
	hfc.cleanupPath(keyStorePath)
	hfc.cleanupPath(credentialStorePath)
}

func (hfc *FabricSetup) cleanupPath(storePath string) {
	err := os.RemoveAll(storePath)
	if err != nil {
		logger.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}
}
