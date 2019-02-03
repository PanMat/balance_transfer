package utils

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/test/integration"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("Helper")

// GetRegisteredUser get registered user. If user is not enrolled, enroll new user
func GetRegisteredUser(username, orgName, secret, identityTypeUser string, sdk *fabsdk.FabricSDK) (string, bool) {
	ctxProvider := sdk.Context(fabsdk.WithOrg(orgName))
	mspClient, err := msp.New(ctxProvider)
	if err != nil {
		log.Fatalf("Failed to create msp client: %s", err.Error())
	}
	signingIdentity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		log.Infof("Check if user %s is enrolled: %s", username, err.Error())
		testAttributes := []msp.Attribute{
			{
				Name:  integration.GenerateRandomID(),
				Value: fmt.Sprintf("%s:ecert", integration.GenerateRandomID()),
				ECert: true,
			},
			{
				Name:  integration.GenerateRandomID(),
				Value: fmt.Sprintf("%s:ecert", integration.GenerateRandomID()),
				ECert: true,
			},
		}

		// Register the new user
		identity, err := mspClient.GetIdentity(username)
		if true {
			log.Infof("User %s does not exist, registering new user", username)
			_, err = mspClient.Register(&msp.RegistrationRequest{
				Name:        username,
				Type:        identityTypeUser,
				Attributes:  testAttributes,
				Affiliation: orgName,
				Secret:      secret,
			})
		} else {
			log.Infof("Identity: %s", identity.Secret)
		}
		//enroll user
		err = mspClient.Enroll(username, msp.WithSecret(secret))
		if err != nil {
			log.Infof("enroll %s failed: %v", username, err)
			return "failed " + err.Error(), false
		}

		return username + " enrolled Successfully", true
	}
	log.Infof("%s: %s", signingIdentity.Identifier().ID, string(signingIdentity.EnrollmentCertificate()[:]))
	return username + " already enrolled", true
}

// GetArgs get [][]byte args from string array
func GetArgs(args []string) [][]byte {
	var result [][]byte
	for _, element := range args {
		result = append(result, []byte(element))
	}
	return result
}

//ConvertBytestoString converts [][]byte to []string
func ConvertBytestoString(data [][]byte) []string {
	s := make([]string, len(data))
	for row := range data {
		s[row] = string(data[row])
	}

	return s
}
