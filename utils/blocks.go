package utils

import (
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//QueryBlockHeight queries the blocks on ledger
func (hfc *FabricSetup) QueryBlockHeight(channelName string, username string, orgName string, endpoint string) []byte {
	channelContext := hfc.Sdk.ChannelContext(channelName, fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)

	if err != nil {
		logger.Errorf("Unable to create new ledger client %s", err)
		return []byte("Unable to create new ledger client")
	}
	bci, err := client.QueryInfo(ledger.WithTargetEndpoints(endpoint))
	ht := bci.BCI.Height
	logger.Debugf("BCI output is : %s", bci)

	return []byte("Current Height of ledger is " + strconv.FormatUint(ht, 10))
}

//QueryBlockByID queries blocks based on input ID
func (hfc *FabricSetup) QueryBlockByID(channelName string, username string, orgName string, endpoint string, blockID uint64) []byte {
	logger.Debugf("----------Recd request to retrieve data for block# :%s ----------", strconv.Itoa(int(blockID)))

	channelContext := hfc.Sdk.ChannelContext(channelName, fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)

	if err != nil {
		logger.Errorf("Unable to create new ledger client %s", err)
		return []byte("Unable to create new ledger client")
	}

	block, err := client.QueryBlock(blockID, ledger.WithTargetEndpoints(endpoint))
	if err != nil {
		logger.Fatalf("Query by Block number returned error: %s", err)
	}
	if block.Data == nil {
		logger.Fatal("Query by Block number has nil data")
		return []byte("Block has nil data")
	}

	return processBlock(block)
}
