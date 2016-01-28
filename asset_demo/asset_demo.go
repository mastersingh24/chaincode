/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Vehicle struct {
	Make            string `json:"make"`
	Model           string `json:"model"`
	Reg             string `json:"reg"`
	VIN             int    `json:"VIN"`
	Owner           string `json:"owner"`
	Scrapped        bool   `json:"scrapped"`
	Status          int    `json:"status"`
	Colour          string `json:"colour"`
	V5cID           string `json:"v5cID"`
	LeaseContractID string `json:"leaseContractID"`
}

func (t *SimpleChaincode) create(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting vehicle record")
	}

	var v Vehicle
	var err error

	err = json.Unmarshal([]byte(args[0]), &v)
	if err != nil {
		return nil, errors.New("Invalid vehicle record")
	}

	if v.V5cID == "" {
		return nil, errors.New("Invalid vehicle record - 'v5cID' missing")
	}

	record, err := stub.GetState(v.V5cID)
	if record != nil {
		fmt.Println("ERROR: Vehicle already exists")
		return nil, errors.New("Vehicle already exists")
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, errors.New("Error creating vehicle record")
	}
	err = stub.PutState(v.V5cID, bytes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created new vehicle %+v\n", v)
	return nil, nil

}

func (t *SimpleChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var v Vehicle
	var err error
	var bytes []byte

	bytes, err = stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error retrieving vehicle with v5cID = " + args[0])
	}

	err = json.Unmarshal(bytes, &v)
	if err != nil {
		return nil, errors.New("Corrupt vehicle record")
	}

	v.Owner = args[1]

	bytes, err = json.Marshal(v)

	if err != nil {
		return nil, errors.New("Error creating vehicle record")
	}

	err = stub.PutState(v.V5cID, bytes)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//need one arg
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 'v5cID'")
	}

	var err error

	bytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error retrieving vehicle with v5cID = " + args[0])
	}
	fmt.Printf("Found vehicle bytes:%d\n", len(bytes))
	fmt.Printf("Found vehicle:\n%s\n", string(bytes))

	return bytes, nil
}

func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "create" {
		//Create an asset with some value
		return t.create(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "transfer" {
		//Create an asset with some value
		return t.transfer(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
