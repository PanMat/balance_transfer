package main

//Version: 7.0
// This is a test chaincode that can create user name with balance and also have the ability to move balance
import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("cafe2")

// TestChaincode : Set up structure for the chaincode
type TestChaincode struct {
}

// Init : Calling Init method
func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Infof("<<<<<<<<<<<<<<<<<<Starting Pankaj's TestChainCode Init Method!!!>>>>>>>>>>>>>")
	//first call will always have only one dummy user with balance
	_, args := stub.GetFunctionAndParameters()

	user := args[0]
	bal, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Error with value received for balance")
	}
	logger.Infof("######### Received %s, with balance of %d ###########", user, bal)

	//Write the state to the ledger
	err = stub.PutState(user, []byte(strconv.Itoa(bal)))
	//check for errors and updates
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info("!!!!!!! Succesfully added First user to the state, Congratulations!!!!!!!!")
	return shim.Success(nil)
}

// Invoke : get function query and invoke appropriately
func (t *TestChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Infof("########### Invoking Cafe2 ###########")
	function, args := stub.GetFunctionAndParameters()

	if function == "add" {
		// Add a new user to its state
		return t.add(stub, args)
	}
	if function == "move" {
		// Move balance from one user to another
		return t.move(stub, args)
	}

	if function == "query" {
		// Query balance for a user
		return t.query(stub, args)
	}

	if function == "delete" {
		// Query balance for a user
		return t.delete(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'add', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// add : Method for adding a new user
func (t *TestChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person & starting balance to add")
	}
	logger.Infof("Received request to add %s with starting balance %s \n", args[0], args[1])
	//Write the state to the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}
	msg := []byte("Successfully Added " + args[0] + " with balance of $" + args[1])
	return shim.Success(msg)
}

// move : move money from one user to another
func (t *TestChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting - function, first name, second name and amount to transfer")
	}
	user1 := args[0]
	user2 := args[1]
	bal, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transfer amount")
	}

	bal1, err := stub.GetState(user1)
	bal2, err := stub.GetState(user2)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Abal1, _ := strconv.Atoi(string(bal1))
	Abal2, _ := strconv.Atoi(string(bal2))

	if Abal1 < bal {
		return shim.Error("User does not have enough balance in account")
	}
	logger.Infof("Starting balance before TX: %s = %d &  %s = %d", user1, Abal1, user2, Abal2)
	Abal1 = Abal1 - bal
	Abal2 = Abal2 + bal

	//write the updated information to state
	err = stub.PutState(user1, []byte(strconv.Itoa(Abal1)))
	err = stub.PutState(user2, []byte(strconv.Itoa(Abal2)))
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Infof(".................... transfered %d from %s to %s", bal, user1, user2)
	logger.Infof("Updated balance after TX : %s = %d &  %s = %d", user1, Abal1, user2, Abal2)
	msg := []byte("Moved $" + args[2] + " from " + user1 + " to " + user2)
	return shim.Success(msg)
}

func (t *TestChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	bal, err := stub.GetState(args[0])
	Abal, _ := strconv.Atoi(string(bal))
	if err != nil {
		//		jsonResp := "{\"Error\" : \"Failed to get state for " +  + "\"}"
		return shim.Error("Failed to get state for the user")
	}

	//	jsonResp := "{\"Name\":\"" + args[0] + "\",\"Balance\":\"" + string(bal) + "\"}"
	logger.Infof("Balance for %s : %d", args[0], Abal)
	msg := []byte("Balance for " + args[0] + " is : " + string(bal))
	return shim.Success(msg)
}

func (t *TestChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Enter only one user at a time for deletion")
	}

	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("Failed to delete user")
	}
	msg := []byte("DELETED " + args[0] + " from the ledger...")
	return shim.Success(msg)
}

func main() {
	err := shim.Start(new(TestChaincode))
	if err != nil {
		logger.Errorf("Error starting Test chaincode: %s", err)
	}
}
