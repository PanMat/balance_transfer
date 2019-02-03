package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

)

//QueryChainCode queries the chaincode for current balances
func (hfc *FabricSetup) QueryChainCode(channelName string, chainCodeID string, username string, orgName string, fcn string, args [][]byte, endpoint string) []byte {
	queryReq := channel.Request{
		ChaincodeID: chainCodeID,
		Fcn:         fcn,
		Args:        args,
	}

	clienthannelContext := hfc.Sdk.ChannelContext(channelName, fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgName))
	client, err := channel.New(clienthannelContext)
	if err != nil {
		logger.Errorf("Error: %s", err)
		return []byte("Unable to create new channel")
	}
	res, err := client.Query(queryReq, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoint))
	if err != nil {
		logger.Fatalf("Unable to query ledger : %s", err)
	}
	return res.Payload
}
