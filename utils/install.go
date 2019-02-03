package utils

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//InstallChainCode install chain codes on peers with the called Organization
func (hfc *FabricSetup) InstallChainCode(peers string, chainCodeName string, chainCodePath string, chainCodeType string, chainCodeVersion string, orgname string, goPath string) []byte {
	logger.Info("Installing chaincode.........")
	clientContext := hfc.Sdk.Context(fabsdk.WithUser(hfc.AdminUser), fabsdk.WithOrg(orgname))
	//setupp chaincode package before installation
	ccPkg, err := packager.NewCCPackage(chainCodePath, goPath)
	if err != nil {
		logger.Error(err)
		return []byte("Failed to create chaincode package")
	}
	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		logger.Error(err)
		return []byte("Failed to create new resource management client")
	}
	installCCReq := resmgmt.InstallCCRequest{Name: chainCodeName, Path: chainCodePath, Version: chainCodeVersion, Package: ccPkg}
	res, err := resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		logger.Error(err)
		return []byte("Failed to install chaincode on the peers")
	}
	//check the qaulity and outcome of responses from count received
	if len(res) > 0 {
		return []byte("!!!!!Chain code installed successfully!!!!!")
	}
	return []byte("******* Chain code NOT installed *******")

}
