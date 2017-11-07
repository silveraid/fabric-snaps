/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mgmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/securekey/fabric-snaps/configmanager/api"
)

const (
	//number of records to be inserted in hyperledger for valid configuration test
	numOfRecords = 6
	//config contentCannot save state Configuration must be provided
	testConfigValue = "configstringgoeshere"
	//valid messages contain one MSP, one or more peers each having one or more apps
	//for testing valid messages were configured with one msp: two peers each having three apps
	validMsg          = `{"MspID":"msp.one","Peers":[{"PeerID":"peer.zero.example.com","App":[{"AppName":"testAppName","Config":"config for test app name on peer zero"},{"AppName":"appNameOne","Config":"config for appNameOne"},{"AppName":"appNameTwo","Config":"mnopq"}]},{"PeerID":"peer.one.example.com","App":[{"AppName":"appNameOneOnPeerOne","Config":"config for appNameOneOnPeerOne goes here"},{"AppName":"appNameOneOne","Config":"config for appNameOneOne goes here"},{"AppName":"appNameTwo","Config":"BLOne"}]}]}`
	validMsgOne       = `{"MspID":"msp.one","Peers":[{"PeerID":"peer.one.one.example.com","App":[{"AppName":"appNameR","Config":"configstringgoeshere"},{"AppName":"appNameTwo","Config":"config for appNametwo"},{"AppName":"appNameTwo","Config":"mnopq"}]},{"PeerID":"peer.two.two.example.com","App":[{"AppName":"appNameTwoOnPeerOne","Config":"config for appNameTwoOnPeerOne goes here"},{"AppName":"appNameOneTwo","Config":"config for appNameOneTwo goes here"},{"AppName":"appNameTwo","Config":"BLTwo"}]}]}`
	validMsgForMspTwo = `{"MspID":"msp.two","Peers":[{"PeerID":"peer.one.one.example.com","App":[{"AppName":"appNameP","Config":"msptwoconfigforfirstpeer"},{"AppName":"appNameThree","Config":"config for appNameThree"},{"AppName":"appNameTwo","Config":"mnopq"}]},{"PeerID":"peer.two.two.example.com","App":[{"AppName":"appNameThreeOnPeerOne","Config":"config for appNameThreeOnPeerOne goes here"},{"AppName":"appNameOneThree","Config":"config for appNameOneOnThree goes here"},{"AppName":"appNameTwo","Config":"BLThree"}]}]}`
	//misconfigured messages
	nonsenseMsg     = `{"MspID":"asd"}`
	noAppMsg        = `{"MspID":"asd","Peers":[{"PeerID":"peer.zero.example.com"}]}`
	noConfigMsg     = `{"MspID":"asd","Peers":[{"PeerID":"peer.zero.example.com","App":[{"AppName":"testApp"}]}]}`
	noAppIDMsg      = `{"MspID":"asd","Peers":[{"PeerID":"peer.zero.example.com","App":[{"Config":"Qkw="}]}]}`
	emptyAppNameMsg = `{"MspID":"asd","Peers":[{"PeerID":"peer.zero.example.com","App":[{"AppName":""}]}]}`
	noPeerIDMsg     = `{"MspID":"","Peers":[{"App":[{"AppName":"aaa"}]}]}`
	emptyPeerIDMsg  = `{"MspID":"asd","Peers":[{"PeerID":"","App":[{"AppName":"name","Config":"Qkw="}]}]}`
	emptyPeersMsg   = `{"MspID":"asd","Peers":[]}`
)

func TestValidConfiguration(t *testing.T) {
	b := []byte(validMsg)
	keyConfigMap, err := parseConfigMessage(b)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	//verify that key exists in map
	key, err := createConfigKey("msp.one", "peer.zero.example.com", "appNameTwo")
	if err != nil {
		t.Fatalf("Cannot create key %v", err)
	}
	_, present := keyConfigMap[key]
	if present == false {
		t.Fatalf("Key : %s should be in the map", key)
	}

	//verify that key does not exists in map
	key, _ = createConfigKey("non.existing.msp", "peer.zero.example.com", "appName")
	_, present = keyConfigMap[key]
	if present == true {
		t.Fatalf("Key : %s should NOT be in map", key)
	}

	//verify number of rCannot create config key usingecords to be saved to hyperledger
	if len(keyConfigMap) != numOfRecords {
		t.Fatalf("Expected : %d key/value records. Got %d", numOfRecords, len(keyConfigMap))
	}
}

func TestInvalidConfigurations(t *testing.T) {
	//loop through list of misconfigured configuration message
	invalidMessages := []string{nonsenseMsg, noAppMsg,
		noConfigMsg, noAppIDMsg,
		emptyAppNameMsg, noPeerIDMsg,
		emptyPeerIDMsg, emptyPeersMsg}

	for _, message := range invalidMessages {
		b := []byte(message)
		_, err := parseConfigMessage(b)
		if err == nil {
			t.Fatalf("ExCannot create config key usingpected error for message %s", message)
		}
	}

}

func TestConfigKeyToString(t *testing.T) {
	key := api.ConfigKey{MspID: "abc", PeerID: "peer.zero.sk.example", AppName: "testApp"}
	keyStr := configKeyToString(key)
	expectedKeyString := "abc!peer.zero.sk.example!testApp"
	if keyStr != expectedKeyString {
		t.Fatalf("Expected key string %s. Got %s", expectedKeyString, keyStr)
	}
}

func TestInstantiateConfigManager(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
}

func TestSaveValidConfig(t *testing.T) {

	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save state %s", err)
	}
	stub.MockTransactionEnd("saveConfiguration")
	//get saved configs
	mspID := "msp.one"
	criteria, err := api.NewSearchCriteriaByMspID(mspID)
	if err != nil {
		t.Fatalf("Error creating message status search criteria: %v", err)
	}
	stub.MockTransactionStart("queryConfiguration")
	//use criteria by mspID=msp.one
	configMessages, err := configManager.QueryForConfigs(criteria)
	if err != nil {
		t.Fatalf("Cannot query for configs %v", err)
	}
	fmt.Printf("%d\n", len(*configMessages))
	if len(*configMessages) != numOfRecords {
		t.Fatalf("Did not retrieve all configs %v", err)
	}

}

func TestPutStateFailed(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	//put state should fail
	if err := configManager.Save(b); err == nil {
		t.Fatalf("PutState failed %s", err)
	}
}

func TestSearchCriteria(t *testing.T) {
	mspID := "msp.one"
	criteria, err := api.NewSearchCriteriaByMspID(mspID)
	if err != nil {
		t.Fatalf("Error creating message status search criteria: %v", err)
	}
	mspid := criteria.GetMspID()
	if mspid != mspID {
		t.Fatalf("Expected %s retrieved %s for mspid.: ", mspID, mspid)
	}
	criteria, err = api.NewSearchCriteriaByMspID("")
	if err == nil {
		t.Fatalf("Expected error. An empty criterie is not valid")
	}

}

func TestGetFieldsForIndex(t *testing.T) {
	key := api.ConfigKey{}
	if _, err := getFieldsForIndex("abc", key); err == nil {
		t.Fatalf("Expected error:'unknown index'")
	}
	fields, err := getFieldsForIndex(indexMspID, key)
	if err == nil {
		t.Fatalf("Error 'invalid key' expected")
	}
	fmt.Printf("%d\n", len(fields))
}

func TestAddIndexes(t *testing.T) {
	key := api.ConfigKey{}
	key.MspID = "msp"
	key.PeerID = "peer"
	key.AppName = "appname"
	configManagerImpl := configManagerImpl{}

	if err := configManagerImpl.addIndex("", key); err == nil {
		t.Fatalf("Expected error:'Index is empty' ")
	}

	if err := configManagerImpl.addIndex("dddd", key); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using ...' ")
	}
	key.AppName = ""
	if err := configManagerImpl.addIndexes(key); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using empty AppName")
	}

}

func TestGetIndexKey(t *testing.T) {
	configManagerImpl := configManagerImpl{}

	if _, err := configManagerImpl.getIndexKey("", "", nil); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using ...' ")
	}
	if _, err := configManagerImpl.getIndexKey("aaa", "", nil); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using ...' ")
	}

	if _, err := configManagerImpl.getIndexKey("aaa", "sdfsdfs", nil); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using ...' ")
	}
	if _, err := configManagerImpl.getIndexKey("", "", []string{"a"}); err == nil {
		t.Fatalf("Expected error:'Cannot create config key using ...' ")
	}

}
func TestQueryForValidConfigs(t *testing.T) {
	// mspID
	mspID := "msp.one"
	criteria, err := api.NewSearchCriteriaByMspID(mspID)
	if err != nil {
		t.Fatalf("Error creating message status search criteria: %v", err)
	}
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	//store  config messages
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save state %s", err)
	}
	//store another  config messages for msp.one
	b = []byte(validMsgOne)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save state %s", err)
	}

	//store config messages for msp.two
	b = []byte(validMsgForMspTwo)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save state %s", err)
	}

	stub.MockTransactionEnd("saveConfiguration")
	stub.MockTransactionStart("queryConfiguration")
	//use criteria by mspID=msp.one
	configMessages, err := configManager.QueryForConfigs(criteria)
	if err != nil {
		t.Fatalf("Cannot query for configs %v", err)
	}
	//verify that all configs for mspID were returned
	expectedLength := numOfRecords*2 - 1
	if len(*configMessages) != expectedLength {
		t.Fatalf("Expected %d configs for %s", numOfRecords*2, mspID)
	}
	//verify that config map contains proper keys and values
	key, _ := createConfigKey("msp.one", "peer.one.example.com", "appNameOneOne")
	keystr := configKeyToString(key)
	expectedConfigValue := "config for appNameOneOne goes here"
	retrievedConfigValue := (*configMessages)[keystr]
	if expectedConfigValue != retrievedConfigValue {
		t.Fatalf("Expected %s value for key %s but got %s ", expectedConfigValue, keystr, retrievedConfigValue)
	}

	mspID = "msp.two"
	criteria, err = api.NewSearchCriteriaByMspID(mspID)
	if err != nil {
		t.Fatalf("Error creating message status search criteria: %v", err)
	}
	//use criteria by mspID=msp.two
	configMessages, err = configManager.QueryForConfigs(criteria)
	if err != nil {
		t.Fatalf("Cannot query for configs %v", err)
	}
	if len(*configMessages) != numOfRecords {
		t.Fatalf("Expected %d configs for %s", numOfRecords, mspID)
	}

	stub.MockTransactionEnd("queryConfiguration")

}

func TestNoConfigsReturnedForBogusMSP(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("queryConfiguration")

	mspID := "bogus.msp.id"
	criteria, err := api.NewSearchCriteriaByMspID(mspID)
	if err != nil {
		t.Fatalf("Error creating message status search criteria: %v", err)
	}
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	//use criteria by mspID=bogus.msp.id
	configMessages, err := configManager.QueryForConfigs(criteria)
	if err != nil {
		t.Fatalf("Cannot query for configs %v", err)
	}
	if len(*configMessages) > 0 {
		t.Fatalf("Expected no messaged for for %s", mspID)
	}
	stub.MockTransactionEnd("queryConfiguration")

}

func TestParseConfigMessage(t *testing.T) {
	apConfig := api.AppConfig{}
	apConfig.AppName = "abc"
	b, err := json.Marshal(apConfig)
	if err != nil {
		t.Fatalf("Cannot get ApiConfig bytes %v", err)
	}
	if _, err := parseConfigMessage(b); err == nil {
		t.Fatalf("Expected error: 'Cannot unmarshal config message...'")
	}
	if _, err := parseConfigMessage(nil); err == nil {
		t.Fatalf("Expected error 'Cannot unmarshal config message...'%v", err)
	}
	var config []byte
	if _, err := parseConfigMessage(config); err == nil {
		t.Fatalf("Expected error 'Cannot unmarshal config message'")
	}

}

func TestSaveEmptyConfig(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte("")
	if err := configManager.Save(b); err == nil {
		t.Fatalf("Expected error 'Cannot save state Configuration must be provided")
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestSaveInvalidConfig(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(noAppIDMsg)
	if err := configManager.Save(b); err == nil {
		t.Fatalf("Expected error 'Configuration message does not have proper App'")
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestGetWithValidKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save configuration message %s", err)
	}
	key, _ := createConfigKey("msp.one", "peer.zero.example.com", "appNameTwo")
	config, err := configManager.Get(key)
	if err != nil {
		t.Fatalf("Cannot get config for key %s %s", key, err)
	}
	if len(config) == 0 {
		t.Fatalf("Cannot get config content for key %s", key)
	}
	if bytes.Equal([]byte(testConfigValue), config[:]) {
		t.Fatalf("Stored and retrieved content are not the same.Expected %s received %s", testConfigValue, string(config[:]))
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestGetWithInvalidKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save configuration message %s", err)
	}
	key := api.ConfigKey{MspID: "abc", PeerID: ""}
	if _, err := configManager.Get(key); err == nil {
		t.Fatalf("Expected error 'Cannot create key using mspID:abc , peerID , appName'")
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestGetWithNonExistingKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save configuration message %s", err)
	}
	key, _ := createConfigKey("msp.one.does.not.exist", "peer.zero.example.com", "appName")
	config, err := configManager.Get(key)
	if err != nil {
		t.Fatalf("Cannot get config for key %s", err)
	}
	if len(config) > 0 {
		t.Fatalf("Should not get any config for non-existing key %s", key)
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestDeleteWithValidKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save configuration message %s", err)
	}
	key, _ := createConfigKey("msp.one", "peer.zero.example.com", "appName")
	if err := configManager.Delete(key); err != nil {
		t.Fatalf("Cannot delete config for  key %s %s", key, err)
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestDeleteWithNonExistingKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	b := []byte(validMsg)
	if err := configManager.Save(b); err != nil {
		t.Fatalf("Cannot save configuration message %s", err)
	}
	key, _ := createConfigKey("msp.one.some.bogus.key", "peer.zero.example.com", "appName")
	if err := configManager.Delete(key); err != nil {
		t.Fatalf("Cannot delete config for  key %s %s", key, err)
	}
	stub.MockTransactionEnd("saveConfiguration")
}
func TestDeleteWithInvalidKey(t *testing.T) {
	stub := shim.NewMockStub("testConfigState", nil)
	stub.MockTransactionStart("saveConfiguration")
	configManager := NewConfigManager(stub)
	if configManager == nil {
		t.Fatal("Cannot instantiate config manager")
	}
	key := api.ConfigKey{MspID: ""}
	if err := configManager.Delete(key); err == nil {
		t.Fatalf("Expected error 'Cannot create key using mspID: , peerID , appName'")
	}
	stub.MockTransactionEnd("saveConfiguration")
}

func TestCreateConfigKey(t *testing.T) {
	//(mspID string, peerID string, appName string) (api.ConfigKey, error) {
	if _, err := createConfigKey("", "asv", "aaa"); err == nil {
		t.Fatalf("Expected error ")
	}
	if _, err := createConfigKey("safsdf", "", "aaa"); err == nil {
		t.Fatalf("Expected error ")
	}
	if _, err := createConfigKey("sdfsdf", "asv", ""); err == nil {
		t.Fatalf("Expected error ")
	}

}

func TestUnmarshalConfig(t *testing.T) {
	if _, err := unmarshalConfig(nil); err == nil {
		t.Fatalf("Expected error 'No configuration passed to unmarshaller'")
	}
	var config []byte
	if _, err := unmarshalConfig(config); err == nil {
		t.Fatalf("Expected error 'No configuration passed to unmarshaller'")
	}
}
