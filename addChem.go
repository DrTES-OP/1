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

// Item is for storing retreived molecules of records on ChemChaincode

type Item struct {
	ChemID         string // A chemical item identifier
	RecType        string // ChemINV
	ChemDesc       string // Description of state - solid, liquid, colorless, etc
	ChemFormula    string // Chemical Formula
	ChemCAS        string // CAS name
	ChemIUPAC      string // IUPAC name
	ChemSMI        string // Canonical smiles
	ChemInChi      string // InChi
	ChemInChiKey   string // InChi key
	ChemNMRC       string // 13C NMR data
	ChemNMRH       string // 1H PMR data
	ChemMS         string // Mass-spectra
	ChemUV         string // UV-Vis spectra
	ChemElementary string // Elemental Analysis (still must for newly-synthesized compounds)
	ChemXR         string // link to deposited CIF-file for X-Ray
}

// ListUser is for storing retreived User list with status
type ListItem struct {
	UserID         string // Person unique ID
	ChemID         string // A chemical item identifier
	RecType        string // ChemINV
	ChemDesc       string // Description of state - solid, liquid, colorless, etc
	ChemFormula    string // Chemical Formula
	ChemCAS        string // CAS name
	ChemIUPAC      string // IUPAC name
	ChemSMI        string // Canonical smiles
	ChemInChi      string // InChi
	ChemInChiKey   string // InChi key
	ChemNMRC       string // 13C NMR data
	ChemNMRH       string // 1H PMR data
	ChemMS         string // Mass-spectra
	ChemUV         string // UV-Vis spectra
	ChemElementary string // Elemental Analysis (still must for newly-synthesized compounds)
	ChemXR         string // link to deposited CIF-file for X-Ray
}

// CountUser is for storing retreived User count
type CountUser struct {
	Count int `json:"count"`
}

// Init initializes the smart contracts
func (t *ChemChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check if table already exists
	_, err := stub.GetTable("ItemTable")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create Item Table
	err = stub.CreateTable("ItemTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ChemID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "RecType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemDesc", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemFormula", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemCAS", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemIUPAC", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemSMI", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemInChi", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemInChiKeyil", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemNMRC", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemNMRH", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemMS", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemUV", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemElementary", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ChemXR", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating ItemTable.")
	}

	return nil, nil
}

func (t *ChemChaincode) getNumItems(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
	}

	var columns []shim.Column

	contractCounter := 0

	rows, err := stub.GetRows("ItemTable", columns)
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

func (t *ChemChaincode) getItem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ChemID to query")
	}

	ChemID := args[0]

	// Get the row pertaining to this UserId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: ChemID}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ItemTable", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the User " + ChemID + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the User " + ChemID + "\"}"
		return nil, errors.New(jsonResp)
	}

	res2E := Item{}

	res2E.ChemID = row.Columns[0].GetString_()
	res2E.RecType = row.Columns[1].GetString_()
	res2E.ChemDesc = row.Columns[2].GetString_()
	res2E.ChemFormula = row.Columns[3].GetString_()
	res2E.ChemCAS = row.Columns[4].GetString_()
	res2E.ChemIUPAC = row.Columns[5].GetString_()
	res2E.ChemSMI = row.Columns[6].GetString_()
	res2E.ChemInChi = row.Columns[7].GetString_()
	res2E.ChemInChiKey = row.Columns[8].GetString_()
	res2E.ChemNMRC = row.Columns[9].GetString_()
	res2E.ChemNMRH = row.Columns[10].GetString_()
	res2E.ChemMS = row.Columns[11].GetString_()
	res2E.ChemUV = row.Columns[12].GetString_()
	res2E.ChemElementary = row.Columns[13].GetString_()
	res2E.ChemXR = row.Columns[14].GetString_()

	mapB, _ := json.Marshal(res2E)
	fmt.Println(string(mapB))

	return mapB, nil

}

func (t *ChemChaincode) listAllItem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0.")
	}

	var columns []shim.Column

	rows, err := stub.GetRows("ItemTable", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}

	res2E := []*ListItem{}

	for row := range rows {
		newApp := new(ListItem)
		newApp.ChemID = row.Columns[0].GetString_()
		newApp.RecType = row.Columns[1].GetString_()
		newApp.ChemDesc = row.Columns[2].GetString_()
		newApp.ChemFormula = row.Columns[3].GetString_()
		newApp.ChemCAS = row.Columns[4].GetString_()
		newApp.ChemIUPAC = row.Columns[5].GetString_()
		newApp.ChemSMI = row.Columns[6].GetString_()
		newApp.ChemInChi = row.Columns[7].GetString_()
		newApp.ChemInChiKey = row.Columns[8].GetString_()
		newApp.ChemNMRC = row.Columns[9].GetString_()
		newApp.ChemNMRH = row.Columns[10].GetString_()
		newApp.ChemMS = row.Columns[11].GetString_()
		newApp.ChemUV = row.Columns[12].GetString_()
		newApp.ChemElementary = row.Columns[13].GetString_()
		newApp.ChemXR = row.Columns[14].GetString_()

		res2E = append(res2E, newApp)
	}

	res2F, _ := json.Marshal(res2E)
	fmt.Println(string(res2F))
	return res2F, nil

}

// Invoke invokes the chaincode
func (t *ChemChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "submitItem" {
		if len(args) != 15 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 15. Got: %d.", len(args))
		}

		ChemID := args[0]
		RecType := args[1]
		ChemDesc := args[2]
		ChemFormula := args[3]
		ChemCAS := args[4]
		ChemIUPAC := args[5]
		ChemSMI := args[6]
		ChemInChi := args[7]
		ChemInChiKey := args[8]
		ChemNMRC := args[9]
		ChemNMRH := args[10]
		ChemMS := args[11]
		ChemUV := args[12]
		ChemElementary := args[13]
		ChemXR := args[14]

		// Insert a row
		ok, err := stub.InsertRow("ItemTable", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: ChemID}},
				&shim.Column{Value: &shim.Column_String_{String_: RecType}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemDesc}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemFormula}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemCAS}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemIUPAC}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemSMI}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemInChi}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemInChiKey}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemNMRC}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemNMRH}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemMS}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemUV}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemElementary}},
				&shim.Column{Value: &shim.Column_String_{String_: ChemXR}},
			}})

		if err != nil {
			return nil, err
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		return nil, err
	}

	return nil, errors.New("Invalid invoke function name.")

}

func (t *ChemChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getItem" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting ChemID to query")
		}
		t := ChemChaincode{}
		return t.getItem(stub, args)
	} else if function == "listAllItem" {
		t := ChemChaincode{}
		return t.listAllItem(stub, args)
	} else if function == "getNumItems" {
		t := ChemChaincode{}
		return t.getNumItems(stub, args)
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
