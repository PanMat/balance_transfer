package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//InvokeChainCode invokes add, delete operations on the chaincode for transactions
func (hfc *FabricSetup) InvokeChainCode(channelName, chainCodeID, username string, orgName string, fcn string, args [][]byte, endpoints []string) []byte {
	invokeReq := channel.Request{
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

	res, err := client.Execute(invokeReq, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(endpoints...))
	if err != nil {
		logger.Errorf("Error executing on chaincode: %s\n", err)
		return []byte("Unable to invoke chaincode")
	}
	logger.Infof("Chaincode executed succesfully: %s\n", res)
	return res.Payload
}
