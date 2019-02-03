package utils

import (
	"os"
	"strings"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config/lookup"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//CreateChannel creates the channel
func (hfc *FabricSetup) CreateChannel(user string, orgname string, channel string, filepath string) []byte {
	channelCreated := []byte("!!!!!Channel Created Successfully!!!")
	channelNotCreated := []byte("*********Channel NOT Created ************")
	//set admin context for the org to create channel
	logger.Info("Trying to get context for org admin to create channel")
	mspClient, err := mspclient.New(hfc.Sdk.Context(), mspclient.WithOrg(orgname))
	adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		logger.Fatalf("Error getting admin context: %s\n", err)
		return channelNotCreated
	}
	hfc.Identity = adminIdentity

	var cfgBackends []core.ConfigBackend
	configBackend, err := hfc.Sdk.Config()
	if err != nil {
		logger.Errorf("Error grabbing configuration: %s\n", err)
		return channelNotCreated
	}
	cfgBackends = append(cfgBackends, configBackend)

	endpointConfig, err := fab.ConfigFromBackend(configBackend)
	netConfig, _ := endpointConfig.NetworkConfig()
	//logger.Infof("Network configuration is %+v\n", netConfig)

	err = lookup.New(configBackend).UnmarshalKey("organizations", &netConfig.Organizations)
	if err != nil {
		logger.Errorf("failed to get organizations from config: %s\n", err)
		return channelNotCreated
	}

	for _, org := range orgname {
		orgConfig, ok := netConfig.Organizations[strings.ToLower(orgname)]
		if !ok {
			logger.Info(org)
			continue
		}
		hfc.Targets = append(hfc.Targets, orgConfig.Peers...)
	}
	logger.Infof("Target Peers are %+v\n", hfc.Targets)

	r, err := os.Open(filepath)
	if err != nil {
		logger.Fatalf("Unable to open channel file from the path")
		return channelNotCreated
	}
	defer func() {
		if err = r.Close(); err != nil {
			logger.Debugf("File closing error: %v", err)
		}
	}()

	//create channel creation request
	req := resmgmt.SaveChannelRequest{ChannelID: channel, ChannelConfig: r, SigningIdentities: []msp.SigningIdentity{adminIdentity}}

	//prepare context
	clientContext := hfc.Sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(hfc.OrdererOrgName))
	// Channel management client is responsible for managing channels (create/update)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Errorf("Failed to create new channel management client: %s\n", err)
		return channelNotCreated
	}
	// Create channel (or update if it already exists)
	if _, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		logger.Errorf("Unable to create channel: %s\n", err)
		return channelNotCreated
	}

	return channelCreated
}
