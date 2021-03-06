/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/securekey/fabric-snaps/eventserver/pkg/channelutil"
	eventapi "github.com/securekey/fabric-snaps/eventservice/api"
	"github.com/securekey/fabric-snaps/eventservice/pkg/localservice"
)

var logger = shim.NewLogger("EventConsumerSnap")

const (
	// Available function names
	putFunc                       = "put"
	registerBlockFunc             = "registerblock"
	unregisterBlockFunc           = "unregisterblock"
	getBlockeventsFunc            = "getblockevents"
	deleteBlockEventsFunc         = "deleteblockevents"
	registerFilteredBlockFunc     = "registerfilteredblock"
	unregisterFilteredBlockFunc   = "unregisterfilteredblock"
	getFilteredBlockEventsFunc    = "getfilteredblockevents"
	deleteFilteredBlockEventsFunc = "deletefilteredblockevents"
	registerCCFunc                = "registercc"
	unregisterCCFunc              = "unregistercc"
	getCCEventsFunc               = "getccevents"
	deleteCCEventsFunc            = "deleteccevents"
	registerTxFunc                = "registertx"
	unregisterTxFunc              = "unregistertx"
	getTxEventsFunc               = "gettxevents"
	deleteTxEventsFunc            = "deletetxevents"
)

// funcMap is a map of functions by function name
type funcMap map[string]func(shim.ChaincodeStubInterface, []string) pb.Response

// eventConsumerSnap is used in the EventSnap BDD test to test the features of the EventSnap
type eventConsumerSnap struct {
	functions           funcMap
	regmutex            sync.RWMutex
	eventmutex          sync.RWMutex
	blockRegistrations  map[string]eventapi.Registration
	fblockRegistrations map[string]eventapi.Registration
	ccRegistrations     map[string]eventapi.Registration
	txRegistrations     map[string]eventapi.Registration
	blockEvents         map[string][]*eventapi.BlockEvent
	fblockEvents        map[string][]*eventapi.FilteredBlockEvent
	ccEvents            map[string][]*eventapi.CCEvent
	txEvents            map[string][]*eventapi.TxStatusEvent
}

// New chaincode implementation
func New() shim.Chaincode {
	s := &eventConsumerSnap{
		functions:           make(funcMap),
		blockRegistrations:  make(map[string]eventapi.Registration),
		fblockRegistrations: make(map[string]eventapi.Registration),
		ccRegistrations:     make(map[string]eventapi.Registration),
		txRegistrations:     make(map[string]eventapi.Registration),
		blockEvents:         make(map[string][]*eventapi.BlockEvent),
		fblockEvents:        make(map[string][]*eventapi.FilteredBlockEvent),
		ccEvents:            make(map[string][]*eventapi.CCEvent),
		txEvents:            make(map[string][]*eventapi.TxStatusEvent),
	}

	s.functions[registerBlockFunc] = s.registerBlockEvents
	s.functions[unregisterBlockFunc] = s.unregisterBlockEvents
	s.functions[getBlockeventsFunc] = s.getBlockEvents
	s.functions[deleteBlockEventsFunc] = s.deleteBlockEvents
	s.functions[registerFilteredBlockFunc] = s.registerFilteredBlockEvents
	s.functions[unregisterFilteredBlockFunc] = s.unregisterFilteredBlockEvents
	s.functions[getFilteredBlockEventsFunc] = s.getFilteredBlockEvents
	s.functions[deleteFilteredBlockEventsFunc] = s.deleteFilteredBlockEvents
	s.functions[registerCCFunc] = s.registerCCEvents
	s.functions[unregisterCCFunc] = s.unregisterCCEvents
	s.functions[getCCEventsFunc] = s.getCCEvents
	s.functions[deleteCCEventsFunc] = s.deleteCCEvents
	s.functions[registerTxFunc] = s.registerTxEvents
	s.functions[unregisterTxFunc] = s.unregisterTxEvents
	s.functions[getTxEventsFunc] = s.getTxEvents
	s.functions[deleteTxEventsFunc] = s.deleteTxEvents
	s.functions[putFunc] = s.put

	return s
}

// Init registers for events
func (s *eventConsumerSnap) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke invokes various test functions
func (s *eventConsumerSnap) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	functionName, args := stub.GetFunctionAndParameters()
	if functionName == "" {
		return shim.Error("Function name is required")
	}

	function, valid := s.functions[functionName]
	if !valid {
		return shim.Error(fmt.Sprintf("Invalid invoke function [%s]. Expecting one of: %v", functionName, s.functions))
	}

	return function(stub, args)
}

func (s *eventConsumerSnap) registerBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.regmutex.Lock()
	defer s.regmutex.Unlock()

	if _, ok := s.blockRegistrations[channelID]; ok {
		return shim.Error(fmt.Sprintf("Block registration already exists for channel: %s", channelID))
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Registering for block events on channel %s ...\n", channelID)

	reg, eventch, err := eventService.RegisterBlockEvent()
	if err != nil {
		return shim.Error(fmt.Sprintf("Error registering for block events on channel: %s", channelID))
	}

	s.blockRegistrations[channelID] = reg

	go func() {
		logger.Infof("Listening for block events on channel: %s\n")
		for {
			bevent, ok := <-eventch
			if !ok {
				logger.Infof("Stopped listening for block events on channel %s\n", channelID)
				return
			}
			go func() {
				logger.Infof("Received block event: %v\n", bevent.Block)

				chID, err := channelutil.ChannelIDFromBlock(bevent.Block)
				if err != nil {
					logger.Errorf("Error extracting channel ID from block: %s\n", err)
				} else {
					s.eventmutex.Lock()
					defer s.eventmutex.Unlock()
					s.blockEvents[chID] = append(s.blockEvents[chID], bevent)
				}
			}()
		}
	}()

	return shim.Success(nil)
}

func (s *eventConsumerSnap) unregisterBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.regmutex.Lock()
	defer s.regmutex.Unlock()

	reg, ok := s.blockRegistrations[channelID]
	if !ok {
		// No registrations
		return shim.Success(nil)
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Unregistering block events on channel %s\n", channelID)

	eventService.Unregister(reg)

	delete(s.blockRegistrations, channelID)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) getBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.RLock()
	defer s.eventmutex.RUnlock()

	bytes, err := json.Marshal(s.blockEvents[channelID])
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshalling block events: %s", err))
	}

	return shim.Success(bytes)
}

func (s *eventConsumerSnap) deleteBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.Lock()
	defer s.eventmutex.Unlock()

	delete(s.blockEvents, channelID)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) registerFilteredBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.regmutex.Lock()
	defer s.regmutex.Unlock()

	if _, ok := s.fblockRegistrations[channelID]; ok {
		return shim.Error(fmt.Sprintf("Filtered block registration already exists for channel: %s", channelID))
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Registering for filtered block events on channel %s ...\n", channelID)

	reg, eventch, err := eventService.RegisterFilteredBlockEvent()
	if err != nil {
		return shim.Error(fmt.Sprintf("Error registering for filtered block events on channel: %s", channelID))
	}

	s.fblockRegistrations[channelID] = reg

	go func() {
		logger.Infof("Listening for filtered block events on channel %s\n", channelID)
		for {
			fbevent, ok := <-eventch
			if !ok {
				logger.Infof("Stopped listening for filtered block events on channel %s\n", channelID)
				return
			}
			go func() {
				logger.Infof("Received filtered block event: %v\n", fbevent.FilteredBlock)

				chID, err := channelutil.ChannelIDFromFilteredBlock(fbevent.FilteredBlock)
				if err != nil {
					logger.Errorf("Error extracting channel ID from filtered block: %s\n", err)
				} else {
					s.eventmutex.Lock()
					defer s.eventmutex.Unlock()
					s.fblockEvents[chID] = append(s.fblockEvents[chID], fbevent)
				}
			}()
		}
	}()

	return shim.Success(nil)
}

func (s *eventConsumerSnap) unregisterFilteredBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	reg, ok := s.fblockRegistrations[channelID]
	if !ok {
		// No registrations
		return shim.Success(nil)
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Unregistering filtered block events on channel %s\n", channelID)

	eventService.Unregister(reg)

	delete(s.fblockRegistrations, channelID)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) getFilteredBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.RLock()
	defer s.eventmutex.RUnlock()

	bytes, err := json.Marshal(s.fblockEvents[channelID])
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshalling filtered block events: %s", err))
	}

	return shim.Success(bytes)
}

func (s *eventConsumerSnap) deleteFilteredBlockEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.Lock()
	defer s.eventmutex.Unlock()

	delete(s.fblockEvents, channelID)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) registerCCEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error("Expecting channel ID, CC ID, and event filter")
	}

	channelID := args[0]
	ccID := args[1]
	eventFilter := args[2]

	s.regmutex.Lock()
	defer s.regmutex.Unlock()

	regKey := getCCRegKey(channelID, ccID, eventFilter)
	if _, ok := s.ccRegistrations[regKey]; ok {
		return shim.Error(fmt.Sprintf("CC registration already exists for channel %s, CC %s, and event filter %s", channelID, ccID, eventFilter))
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Registering for CC events on channel %s, CC %s, and event filter %s", channelID, ccID, eventFilter)

	reg, eventch, err := eventService.RegisterChaincodeEvent(ccID, eventFilter)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error registering for CC events on channel %s, CC %s, and event filter %s", channelID, ccID, eventFilter))
	}

	s.ccRegistrations[regKey] = reg

	go func() {
		logger.Infof("Listening for chaincode events on channel %s, CC ID: %s, Event filter: %s\n", channelID, ccID, eventFilter)
		for {
			ccevent, ok := <-eventch
			if !ok {
				logger.Infof("Stopped listening for chaincode events on channel %s, CC ID: %s, Event filter: %s\n", channelID, ccID, eventFilter)
				return
			}
			go func() {
				logger.Infof("Received CC event on channel %s, CC ID: %s, Event: %s, TxID: %s\n", channelID, ccID, ccevent.EventName, ccevent.TxID)
				s.eventmutex.Lock()
				defer s.eventmutex.Unlock()
				s.ccEvents[channelID] = append(s.ccEvents[channelID], ccevent)
			}()
		}
	}()

	return shim.Success(nil)
}

func (s *eventConsumerSnap) unregisterCCEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 3 {
		return shim.Error("Expecting channel ID, CC ID, and event filter")
	}

	channelID := args[0]
	ccID := args[1]
	eventFilter := args[2]

	regKey := getCCRegKey(channelID, ccID, eventFilter)
	reg, ok := s.ccRegistrations[regKey]
	if !ok {
		// No registrations
		return shim.Success(nil)
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Unregistering CC events for channel %s, CC %s, and event filter %s", channelID, ccID, eventFilter)

	eventService.Unregister(reg)

	delete(s.ccRegistrations, regKey)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) getCCEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.RLock()
	defer s.eventmutex.RUnlock()

	bytes, err := json.Marshal(s.ccEvents[channelID])
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshalling CC events: %s", err))
	}

	return shim.Success(bytes)
}

func (s *eventConsumerSnap) deleteCCEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.Lock()
	defer s.eventmutex.Unlock()

	delete(s.ccEvents, channelID)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) registerTxEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Expecting channel ID, and Tx ID")
	}

	channelID := args[0]
	txID := args[1]

	s.regmutex.Lock()
	defer s.regmutex.Unlock()

	regKey := getTxRegKey(channelID, txID)
	if _, ok := s.txRegistrations[regKey]; ok {
		return shim.Error(fmt.Sprintf("Tx Status registration already exists for channel %s and TxID %s", channelID, txID))
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Registering for Tx Status events on channel %s and TxID %s", channelID, txID)

	reg, eventch, err := eventService.RegisterTxStatusEvent(txID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error registering for Tx Status events on channel %s and TxID %s", channelID, txID))
	}

	s.txRegistrations[regKey] = reg

	go func() {
		logger.Infof("Listening for Tx Status events on channel %s for Tx: %s\n", channelID, txID)

		txevent, ok := <-eventch
		if !ok {
			logger.Infof("Stopped listening for Tx Status events for Tx: %s\n", txID)
			return
		}
		go func() {
			s.eventmutex.Lock()
			defer s.eventmutex.Unlock()
			logger.Infof("Received Tx Status event - TxID: %s, Status: %s\n", txevent.TxID, txevent.TxValidationCode)
			s.txEvents[channelID] = append(s.txEvents[channelID], txevent)
		}()
	}()

	return shim.Success(nil)
}

func (s *eventConsumerSnap) unregisterTxEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Expecting channel ID and Tx ID")
	}

	channelID := args[0]
	txID := args[1]

	regKey := getTxRegKey(channelID, txID)
	reg, ok := s.txRegistrations[regKey]
	if !ok {
		// No registrations
		return shim.Success(nil)
	}

	eventService := localservice.Get(channelID)
	if eventService == nil {
		return shim.Error(fmt.Sprintf("No local event service for channel: %s", channelID))
	}

	logger.Infof("Unregistering Tx Status events for channel %s and Tx ID %s", channelID, txID)

	eventService.Unregister(reg)

	delete(s.txRegistrations, regKey)

	return shim.Success(nil)
}

func (s *eventConsumerSnap) getTxEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	s.eventmutex.RLock()
	defer s.eventmutex.RUnlock()

	bytes, err := json.Marshal(s.txEvents[channelID])
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshalling Tx Status events: %s", err))
	}

	return shim.Success(bytes)
}

func (s *eventConsumerSnap) deleteTxEvents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Expecting channel ID")
	}

	channelID := args[0]

	delete(s.txEvents, channelID)

	return shim.Success(nil)
}

// put is called to generate events
func (s *eventConsumerSnap) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Expecting key, value, and optional event")
	}

	key := args[0]
	value := args[1]

	var eventName string
	if len(args) > 2 {
		eventName = args[2]
	}

	if err := stub.PutState(key, []byte(value)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting state: %s", err))
	}

	if eventName != "" {
		stub.SetEvent(eventName, nil)
	}

	return shim.Success(nil)
}

func getCCRegKey(channelID, ccID, eventFilter string) string {
	return "cc_" + channelID + "_" + ccID + "_" + eventFilter
}

func getTxRegKey(channelID, txID string) string {
	return "tx_" + channelID + "_" + txID
}

func main() {
}
