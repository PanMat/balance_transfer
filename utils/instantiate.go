package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// InstantiateChainCode instantiates the installed chaincode on the network
func (hfc *FabricSetup) InstantiateChainCode(user string, orgName string, channelName string, chainCodeName string, chainCodePath string, chainCodeVersion string, args [][]byte) []byte {
	logger.Info("Instantiating chain code.........")
	clientContext := hfc.Sdk.Context(fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgName))

	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error("=======================N O T E   E R R O R=======================")
		logger.Error(err)
		logger.Error("=======================N O T E   E R R O R=======================")
		return []byte("Failed to create new resource management client")
	}
	target := []string{"Org1MSP", "Org2MSP"}
	ccPolicy := cauthdsl.SignedByAnyMember(target)

	//create instantiate request
	instantiateReq := resmgmt.InstantiateCCRequest{
		Name:    chainCodeName,
		Path:    chainCodePath,
		Version: chainCodeVersion,
		Args:    args,
		Policy:  ccPolicy,
		//CollConfig:
	}

	res, err := resMgmtClient.InstantiateCC(channelName, instantiateReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error("=======================N O T E   E R R O R=======================")
		logger.Error(err)
		logger.Error("=======================N O T E   E R R O R=======================")
		return []byte("Unable to instantiate chaincode")
	}
	if res.TransactionID == "" {
		return []byte("Failed to instantiate chaincode")
	}

	return []byte("Chain code instantiated successfully........")
}
