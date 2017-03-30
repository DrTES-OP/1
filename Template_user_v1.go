package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/op/go-logging"
)

var myLogger = logging.MustGetLogger("asset_mgm")

// ChemChaincode example simple Chaincode implementation
type ChemChaincode struct {
}

// ===================================================================================
// Main
// ===================================================================================

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(ChemChaincode))
	if err != nil {
		fmt.Printf("Error starting ChemChaincode: %s", err)
	}
}

// ===================================================================================
// Custom stuff
// ===================================================================================

//custom data models
type user struct {
	ObjectType  string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	UserId      string `json:"userId"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Affiliation string `json:"affiliation"` // University, company, etc
	Address     string `json:"address"`     //
	Phone       string `json:"phone"`       //
	Email       string `json:"email"`       //
}

// Init initializes chaincode
// ===========================
func (t *ChemChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	myLogger.Info("Successfully deployed chain code")
	return shim.Success(nil)
}

//////////////////////////////////////////////////////////////
// Invoke Functions based on Function name
// The function name gets resolved to one of the following calls
// during an invoke
//
//////////////////////////////////////////////////////////////
func InvokeFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	InvokeFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response{}
	return InvokeFunc[fname]
}

// Putting things into ledger
func (t *ChemChaincode) CreateChemData(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	myLogger.Debug("adding data...")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
// ===========================

func (t *ChemChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	//function, args := stub.GetFunctionAndParameters()
	//myLogger.Debugf("Query [%s]", function)
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *ChemChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initUser" { //create new user
		return t.initUser(stub, args)
	} else if function == "readUser" { //read user
		return t.readUser(stub, args)
	} else if function == "queryUsers" { //find Users based on an ad hoc rich query
		return t.queryUsers(stub, args)
	} else if function == "getHistoryForUser" { //get history of a user
		return t.getHistoryForUser(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// init users- create new user info, store into chaincode state
// ============================================================
func (t *ChemChaincode) initUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1      2      	3						4					 5				6					7					8
	//UserID, status, Title, First Name, last Name, Affiliation, Address, Phone, Email

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init of user")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}
	if len(args[8]) <= 0 {
		return shim.Error("9th argument must be a non-empty string")
	}
	UserId := args[0]
	status := args[1]
	title := args[2]
	firstName := args[3]
	lastName := args[4]
	affiliation := args[5]
	address := args[6]
	phone := args[7]
	email := args[8]

	// ==== Check if user already exists ====
	userAsBytes, err := stub.GetState(UserId)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + UserId)
		return shim.Error("This user already exists: " + UserId)
	}

	// ==== Create marble object and marshal to JSON ====
	objectType := "user"
	user := &user{objectType, UserId, status, title, firstName, lastName, affiliation, address, phone, email}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	// === Save marble to state ===
	err = stub.PutState(UserId, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// ==== user saved and indexed. Return success ====
	fmt.Println("- end init user")
	return shim.Success(nil)
}

// ===============================================
// readUser - read a marble from chaincode state
// ===============================================
func (t *ChemChaincode) readUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var lastName, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting last name of the user to query")
	}

	lastName = args[0]
	valAsbytes, err := stub.GetState(lastName) //get the user from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + lastName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + lastName + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===== Example: Ad hoc rich query ========================================================
// queryMarbles uses a query string to perform a query for marbles.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *ChemChaincode) queryUsers(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResultKey, queryResultRecord, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultRecord))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

//=========================================================================================
//History of User
//=========================================================================================

func (t *ChemChaincode) getHistoryForUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	UserId := args[0]

	fmt.Printf("- start getHistoryForUser: %s\n", UserId)

	resultsIterator, err := stub.GetHistoryForKey(UserId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		txID, historicValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(txID)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// historicValue is a JSON marble, so we write as-is
		buffer.WriteString(string(historicValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
