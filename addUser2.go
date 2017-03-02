package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

// ChemChaincode example simple Chaincode implementation
type ChemChaincode struct {
}

// User is for storing retreived owners of records on ChemChaincode

type User struct {
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

// ListUser is for storing retreived User list with status
type ListUser struct {
	UserId string `json:"userId"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

// CountUser is for storing retreived User count
type CountUser struct {
	Count int `json:"count"`
}

// Init initializes the smart contracts
func (t *ChemChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check if table already exists
	_, err := stub.GetTable("UserTable")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create User Table
	err = stub.CreateTable("UserTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "userId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "status", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "title", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "firstName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "lastName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "affiliation", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "address", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "phone", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "email", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating UserTable.")
	}

	return nil, nil
}

func (t *ChemChaincode) getNumUsers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
	}

	var columns []shim.Column

	contractCounter := 0

	rows, err := stub.GetRows("UserTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	for row := range rows {
		if len(row.Columns) != 0 {
			contractCounter++
		}
	}

	res2E := CountUser{}
	res2E.Count = contractCounter
	mapB, _ := json.Marshal(res2E)
	fmt.Println(string(mapB))

	return mapB, nil
}

func (t *ChemChaincode) UpdateStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	UserId := args[0]
	newStatus := args[1]

	// Get the row pertaining to this UserId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("UserTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving User with UserId %s. Error %s", UserId, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}

	//currStatus := row.Columns[1].GetString_()

	//End- Check that the currentStatus to newStatus transition is accurate
	// Delete the row pertaining to this UserId
	err = stub.DeleteRow(
		"UserTable",
		columns,
	)
	if err != nil {
		return nil, errors.New("Failed deleting row.")
	}

	//UserId := row.Columns[0].GetString_()
	status := newStatus
	title := row.Columns[2].GetString_()
	firstName := row.Columns[3].GetString_()
	lastName := row.Columns[4].GetString_()
	affiliation := row.Columns[5].GetString_()
	address := row.Columns[6].GetString_()
	phone := row.Columns[7].GetString_()
	email := row.Columns[8].GetString_()

	//Insert the row pertaining to this UserId with new status
	_, err = stub.InsertRow(
		"UserTable",
		shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: UserId}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},
				&shim.Column{Value: &shim.Column_String_{String_: title}},
				&shim.Column{Value: &shim.Column_String_{String_: firstName}},
				&shim.Column{Value: &shim.Column_String_{String_: lastName}},
				&shim.Column{Value: &shim.Column_String_{String_: affiliation}},
				&shim.Column{Value: &shim.Column_String_{String_: address}},
				&shim.Column{Value: &shim.Column_String_{String_: phone}},
				&shim.Column{Value: &shim.Column_String_{String_: email}},
			}})
	if err != nil {
		return nil, errors.New("Failed inserting row.")
	}

	return nil, nil

}

func (t *ChemChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Userid to query")
	}

	UserId := args[0]

	// Get the row pertaining to this UserId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("UserTable", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the User " + UserId + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the User " + UserId + "\"}"
		return nil, errors.New(jsonResp)
	}

	res2E := User{}

	res2E.UserId = row.Columns[0].GetString_()
	res2E.Status = row.Columns[1].GetString_()
	res2E.Title = row.Columns[2].GetString_()
	res2E.FirstName = row.Columns[3].GetString_()
	res2E.LastName = row.Columns[4].GetString_()
	res2E.Affiliation = row.Columns[5].GetString_()
	res2E.Address = row.Columns[6].GetString_()
	res2E.Phone = row.Columns[7].GetString_()
	res2E.Email = row.Columns[8].GetString_()

	mapB, _ := json.Marshal(res2E)
	fmt.Println(string(mapB))

	return mapB, nil

}

func (t *ChemChaincode) listAllUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
	}

	var columns []shim.Column

	rows, err := stub.GetRows("UserTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	res2E := []*ListUser{}

	for row := range rows {
		newApp := new(ListUser)
		newApp.UserId = row.Columns[0].GetString_()
		newApp.Status = row.Columns[1].GetString_()
		newApp.Title = row.Columns[2].GetString_()
		res2E = append(res2E, newApp)
	}

	res2F, _ := json.Marshal(res2E)
	fmt.Println(string(res2F))
	return res2F, nil

}

// Invoke invokes the chaincode
func (t *ChemChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "submitUser" {
		if len(args) != 9 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 9. Got: %d.", len(args))
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

		// Insert a row
		ok, err := stub.InsertRow("UserTable", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: UserId}},
				&shim.Column{Value: &shim.Column_String_{String_: status}},
				&shim.Column{Value: &shim.Column_String_{String_: title}},
				&shim.Column{Value: &shim.Column_String_{String_: firstName}},
				&shim.Column{Value: &shim.Column_String_{String_: lastName}},
				&shim.Column{Value: &shim.Column_String_{String_: affiliation}},
				&shim.Column{Value: &shim.Column_String_{String_: address}},
				&shim.Column{Value: &shim.Column_String_{String_: phone}},
				&shim.Column{Value: &shim.Column_String_{String_: email}},
			}})

		if err != nil {
			return nil, err
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		return nil, err
	} else if function == "updateUserStatus" {
		t := ChemChaincode{}
		return t.UpdateStatus(stub, args)
	}

	return nil, errors.New("Invalid invoke function name.")

}

func (t *ChemChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getUser" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting Userid to query")
		}
		t := ChemChaincode{}
		return t.getUser(stub, args)
	} else if function == "listAllUser" {
		t := ChemChaincode{}
		return t.listAllUser(stub, args)
	} else if function == "getNumUsers" {
		t := ChemChaincode{}
		return t.getNumUsers(stub, args)
	}

	return nil, nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(ChemChaincode))
	if err != nil {
		fmt.Printf("Error starting ChemChaincode: %s", err)
	}
}
