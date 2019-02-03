package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//JoinChannel joins peer(s) to a channel
func (hfc *FabricSetup) JoinChannel(channelName string, peer string, orgname string) []byte {
	logger.Infof("Joining %s from %s to join the channel %s ", peer, orgname, channelName)
	clientContext := hfc.Sdk.Context(fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgname))
	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error(err)
		return []byte("Failed to create new resource management client")
	}
	if err = resMgmtClient.JoinChannel(channelName, resmgmt.WithRetry(retry.DefaultResMgmtOpts)); err != nil {
		logger.Error(err)
		return []byte("Peer unable to join the channel")
	}

	//send confirmation of peer joining back to main program
	s := peer + "Joined the channel successfully............"
	return []byte(s)
}
