/*
* Letter of Credit Trade Finance
* @author Sunbum Lee (jake@netobjex.com)
* Copyright NetObjex, Inc. 2019 All Rights Reserved.
**/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Response Model
type Response struct {
	Key    string `json:"key"`
	Record string `json:"record"`
}

// History Model
type History struct {
	TxID      string `json:"tx_id"`
	Value     string `json:"value"`
	TimeStamp string `json:"timestamp"`
	IsDelete  bool   `json:"is_delete"`
}

func getSampleData(fileName string, objectName string) []byte {
	var objs map[string]interface{}
	var tempByte []byte
	file, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(file, &objs)
	tempByte, _ = json.Marshal(objs[objectName])
	if tempByte == nil {
		return nil
	}

	return tempByte
}

func createByID(stub shim.ChaincodeStubInterface, args []string, template interface{}) error {
	logger.Debugf("ARGS %v: ", args)

	if len(args) != 2 {
		return fmt.Errorf("[ERROR] The number of arguments should be 2")
	}

	recordKey := args[0]
	recordBody := args[1]

	recordByte, err := stub.GetState(recordKey)
	if err != nil || recordByte != nil {
		return fmt.Errorf("[ERROR] The record with the ID %s already exist", recordKey)
	}

	json.Unmarshal([]byte(recordBody), &template)

	newRecordByte, err := json.Marshal(template)

	if err != nil {
		return fmt.Errorf("[ERROR] On Marshaling JSON %s", template)
	}

	stub.PutState(recordKey, newRecordByte)
	logger.Infof("A new loc added with Key: %s, Data: %s", recordKey, string(newRecordByte))
	return nil
}

func queryByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, fmt.Errorf("[ERROR] The number of arguments should be 1")
	}

	aKey := args[0]
	returnData, err := stub.GetState(aKey)

	if err != nil || returnData == nil {
		return nil, fmt.Errorf("[ERROR] Failed to get state for %s", aKey)
	}

	logger.Infof("Search Key: %s, Data: %s", aKey, returnData)
	return returnData, nil
}

func updateRecord(stub shim.ChaincodeStubInterface, args []string, dataByte []byte, completeModel interface{}, updateModel interface{}) error {

	if len(args) != 2 {
		return fmt.Errorf("[ERROR] The number of arguments should be 2")
	}

	searchKey := args[0]
	updateRecord := args[1]

	json.Unmarshal(dataByte, &completeModel)
	json.Unmarshal([]byte(updateRecord), &updateModel)
	updateRecordByte, _ := json.Marshal(updateModel)
	json.Unmarshal(updateRecordByte, &completeModel)
	newRecordUpdateByte, _ := json.Marshal(completeModel)
	stub.PutState(searchKey, newRecordUpdateByte)

	logger.Infof("Update User Response:%s\n", string(newRecordUpdateByte))
	return nil
}

func checkRestrictedQuery(t *testing.T, stub *shim.MockStub, query string, key string, value string, access string) {
	response := stub.MockInvoke("1", [][]byte{[]byte(query), []byte(key), []byte(access)})
	if response.Status == 500 {
		t.Skipf("RESPONSE: %s", string(response.Status))
	}
	if response.Status != shim.OK {
		fmt.Println("Query: ", key, " failed ", string(response.Message))
		t.FailNow()
	}
	if response.Payload == nil {
		fmt.Println("Data with Query: ", key, " not found with ")
		t.FailNow()
	}
	// if string(response.Payload) != value {
	// 	fmt.Println("Data with Query, ", key, " are not same with ", value)
	// 	t.FailNow()
	// }
	if reflect.DeepEqual(string(response.Payload), value) {
		fmt.Println("Data with Query, ", key, " are not same with ", value)
		t.FailNow()
	}
}

func stringify(instance interface{}) string {
	instanceByte, _ := json.Marshal(instance)
	instanceStr := string(instanceByte)
	return instanceStr
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInit("1", args)
	if response.Status != shim.OK {
		fmt.Println("ERROR: checkInit: Chaincode Init fail")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("1", args)
	if response.Status != shim.OK {
		fmt.Println("ERROR: checkInvoke due to " + response.Message)
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, query string, key string, value interface{}) {
	response := stub.MockInvoke("1", [][]byte{[]byte(query), []byte(key)})
	var temp interface{}
	json.Unmarshal(response.Payload, &temp)
	if response.Status == 500 {
		fmt.Println("RESPONSE::: ", response.Payload)
		t.Skipf("RESPONSE: %s", string(response.Status))
	}
	if response.Status != shim.OK {
		fmt.Println("Query: ", key, " failed ", string(response.Message))
		t.FailNow()
	}
	if response.Payload == nil {
		fmt.Println("Data with Query: ", key, " not found with ")
		t.FailNow()
	}
	// if string(response.Payload) != value {
	// 	fmt.Println("Data with Query, ", key, " are not same with ", value)
	// 	t.FailNow()
	// }

	if reflect.DeepEqual(temp, value) {
		fmt.Println("Data with Query, ", key, " are not same with ", value)
		t.FailNow()
	}
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*[]Response, error) {

	var responseList []Response

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		response := Response{
			Key: queryResponse.Key,
			//Record: string(queryResponse.Value),
		}
		responseList = append(responseList, response)
	}
	return &responseList, nil
}

// func getAccessControl(userID string, accessControls []AccessControl) int {
// 	for _, ac := range accessControls {
// 		if userID == ac.ID {
// 			return ac.Type
// 		}
// 	}
// 	return 0
// }

func authenticate(stub shim.ChaincodeStubInterface, searchKey string, userEmail string) (bool, []byte) {
	var temp map[string]interface{}
	dataByte, err := stub.GetState(searchKey)
	if dataByte == nil || err != nil {
		return false, nil
	}
	json.Unmarshal(dataByte, &temp)
	if temp["email"] != userEmail {
		return false, nil
	}
	return true, dataByte
}

func jsonPretty(jsonObject interface{}) []byte {
	b, _ := json.MarshalIndent(jsonObject, "", "\t")
	return b
}
