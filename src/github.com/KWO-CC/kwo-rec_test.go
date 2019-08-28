package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Test_createRecord(t *testing.T) {
	scc := new(KilloWattRecords)
	stub := shim.NewMockStub("KilloWattRecords", scc)
	stub.MockPeerChaincode("KilloWattRecords", stub)
	checkInit(t, stub, nil)

	record1ID := "record1"
	record1Body := RecordModel{
		RID:        record1ID,
		ProducerID: "producer1",
		ConsumerID: "consumer1",
		Kw:         "177",
		Amount:     "3440",
		Timestamp:  53344,
	}

	checkInvoke(t, stub, [][]byte{
		[]byte("createRecord"),
		[]byte(record1ID),
		[]byte(stringify(record1Body)),
	})
	checkQuery(t, stub, "queryRecord", record1ID, record1Body)

	record2ID := "record2"
	record2Body := RecordModel{
		RID:        record1ID,
		ProducerID: "producer1",
		ConsumerID: "consumer1",
		Kw:         "332",
		Amount:     "4000",
		Timestamp:  53344,
	}

	checkInvoke(t, stub, [][]byte{
		[]byte("createRecord"),
		[]byte(record2ID),
		[]byte(stringify(record2Body)),
	})
	checkQuery(t, stub, "queryRecord", record2ID, record2Body)

	record3ID := "record3"
	record3Body := RecordModel{
		RID:        record1ID,
		ProducerID: "producer1",
		ConsumerID: "consumer1",
		Kw:         "17",
		Amount:     "3230",
		Timestamp:  53344,
	}

	checkInvoke(t, stub, [][]byte{
		[]byte("createRecord"),
		[]byte(record3ID),
		[]byte(stringify(record3Body)),
	})
	checkQuery(t, stub, "queryRecord", record3ID, record3Body)
}
