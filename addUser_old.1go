/*

Our Chaincode
NST/Salix
2017

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/crypto/primitives"
)

//----------------------------------------
// Create Owner
// Could establish valid UserTypes -
// USER - human for now
// Robots, IoT, M2M for later
//----------------------------------------

type User struct {
	UserID      int    `json:"userid"`      // Person unique ID
	RecType     string `json:"recordtype"`  // Type = USER
	Name        string `json:"name"`        // Full name of the person
	Affiliation string `json:"affiliation"` // University, company, etc
	Address     string `json:"address"`     //
	Phone       string `json:"phone"`       //
	Email       string `json:"email"`       //
}

// ChemChaincode example simple Chaincode implementation
type ChemChaincode struct {
}

var indexstr string = "recordUserID"

func main() {
	err := shim.Start(new(ChemChaincode))
	if err != nil {
		fmt.Printf("Error starting Chemistry chaincode: %s", err)
	}
}

// Init resets all the things
func (t *ChemChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var empty []string
	indexAsbytes, _ := json.Marshal(empty)
	err := stub.PutState(indexstr, indexAsbytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke issued entry point to invoke a chaincode function
func (t *ChemChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addRecord" {
		return t.addRecord(stub, args)
	} else if function == "modify" {
		return t.modify(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *ChemChaincode) addRecord(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	var userid string = args[1] //5 arguments for User + RecType and UserID are from inside
	var index []string
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of args, expected 5 for record entry")
	}

	str := `{"userid": ` + args[0] + `, "name": ` + args[1] + `, "affiliation": "` + args[2] + `", "address": ` + args[3] + `, "phone":` + args[4] + `, "email":"` + args[5] + `"}`
	err = stub.PutState(userid, []byte(str))
	if err != nil {
		return nil, err
	}

	indexAsbytes, err := stub.GetState(indexstr)
	json.Unmarshal(indexAsbytes, &index)
	index = append(index, args[0])
	newindexAsbytes, err := json.Marshal(index)
	err = stub.PutState(indexstr, newindexAsbytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

// Query is our entry point for queries

func (t *ChemChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getInfo" {
		return t.getInfo(stub, args)
	} else if function == "seeAll" {
		return t.seeAll(stub, args)
	}

	fmt.Println("didnt find any function" + function)

	return nil, errors.New("unknown query")
}

func (t *ChemChaincode) getInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	if len(args) != 1 {
		return nil, errors.New("wrong number of arguments to get info")
	}

	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("couldnt get the record, check id sent")
	}

	return valAsbytes, nil

}

func (t *ChemChaincode) seeAll(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var index []string

	if len(args) != 0 {
		return nil, errors.New("expected 0 arguments")
	}
	valAsbytes, err := stub.GetState(indexstr)
	if err != nil {
		return nil, errors.New("error!!")
	}

	json.Unmarshal(valAsbytes, &index)
	var allResults string
	for i := range index {
		oneResult, err := stub.GetState(index[i])
		if err != nil {
			return nil, errors.New("error!!")
		}
		allResults = allResults + string(oneResult[:])
	}
	return []byte(allResults), nil

}

func (t *ChemChaincode) modify(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	if len(args) != 3 {
		return nil, errors.New("number of arguments are wrong")
	}
	field := args[1]
	value := args[2]
	valAsbytes, err := stub.GetState(args[0])
	modifiedAC := User{}
	json.Unmarshal(valAsbytes, &modifiedAC)
	if field == "name" {
		modifiedAC.Name = value
	} else if field == "affiliation" {
		modifiedAC.Affiliation = value
	} else if field == "address" {
		modifiedAC.Address = value
	} else if field == "phone" {
		modifiedAC.Phone = value
	} else if field == "email" {
		modifiedAC.Email = value
	} else if field == "userid" {
		temp1, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.New("couldnt update")
		}
		modifiedAC.UserID = temp1
		err = stub.DelState(args[0])
	} else {
		return nil, errors.New("no right field to be changed")
	}

	str := `{"name": "` + modifiedAC.Name + `", "affiliation":` + modifiedAC.Affiliation + `, "address":"` + modifiedAC.Address + `", "phone":"` + modifiedAC.Phone + `", "email":"` + modifiedAC.Email + `"}`
	err = stub.PutState(strconv.Itoa(modifiedAC.UserID), []byte(str))

	if err != nil {
		return nil, errors.New("couldnt update")
	}
	return nil, nil
}

// Write - invoke function to write key / value pair

// Read - query function to read key / value pair
