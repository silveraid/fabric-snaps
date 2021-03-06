/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	sdkConfigApi "github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	sdkApi "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	transactionsnapApi "github.com/securekey/fabric-snaps/transactionsnap/api"
)

//GetNewClientWrapper returns wrapper mock object of client
func GetNewClientWrapper(fcClient transactionsnapApi.Client) *MockClient {
	return &MockClient{fcClient: fcClient}
}

//MockClient wrapper for client.Client which can be manipulated for desired results for tests
type MockClient struct {
	fcClient transactionsnapApi.Client
}

// NewChannel registers a channel object with the fabric client
// this object represents a channel on the fabric network
// @param {string} name of the channel
// @returns {Channel} channel object
// @returns {error} error, if any
func (c *MockClient) NewChannel(name string) (sdkApi.Channel, error) {
	return c.fcClient.NewChannel(name)
}

// GetChannel returns a channel object that has been added to the fabric client
// @param {string} name of the channel
// @returns {Channel} channel that was requested
// @returns {error} error, if any
func (c *MockClient) GetChannel(name string) (sdkApi.Channel, error) {
	return c.fcClient.GetChannel(name)
}

// EndorseTransaction request endorsement from the peers on this channel
// for a transaction with the given parameters
// @param {Channel} channel on which we want to transact
// @param {EndorseTxRequest} request identifies the chaincode to invoke
// @returns {[]TransactionProposalResponse} responses from endorsers
// @returns {error} error, if any
func (c *MockClient) EndorseTransaction(channel sdkApi.Channel, request *transactionsnapApi.EndorseTxRequest) (
	[]*apitxn.TransactionProposalResponse, error) {
	return c.fcClient.EndorseTransaction(channel, request)
}

// CommitTransaction submits the given endorsements on the specified channel for
// commit
// @param {Channel} channel on which the transaction is taking place
// @param {[]TransactionProposalResponse} responses from endorsers
// @param {bool} register for Tx event
// @returns {error} error, if any
func (c *MockClient) CommitTransaction(channel sdkApi.Channel, txres []*apitxn.TransactionProposalResponse, register bool) error {
	return c.fcClient.CommitTransaction(channel, txres, register)
}

// QueryChannels joined by the given peer
// @param {Peer} The peer to query
// @returns {[]string} list of channels
// @returns {error} error, if any
func (c *MockClient) QueryChannels(peer sdkApi.Peer) ([]string, error) {
	return c.fcClient.QueryChannels(peer)
}

// VerifyTxnProposalSignature verify TxnProposalSignature against msp
// @param {Channel} channel on which the transaction is taking place
// @param {[]byte} Txn Proposal
// @returns {error} error, if any
func (c *MockClient) VerifyTxnProposalSignature(channel sdkApi.Channel, bytes []byte) error {
	return c.fcClient.VerifyTxnProposalSignature(channel, bytes)
}

// SetSelectionService is used to inject a selection service for testing
// @param {SelectionService} SelectionService
func (c *MockClient) SetSelectionService(service transactionsnapApi.SelectionService) {
	c.fcClient.SetSelectionService(service)
}

// GetSelectionService returns the SelectionService
func (c *MockClient) GetSelectionService() transactionsnapApi.SelectionService {
	return c.fcClient.GetSelectionService()
}

// GetEventHub returns the GetEventHub
// @returns {EventHub} EventHub
// @returns {error} error, if any
func (c *MockClient) GetEventHub() (sdkApi.EventHub, error) {
	return c.fcClient.GetEventHub()
}

// Hash message
// @param {[]byte} message to hash
// @returns {[[]byte} hash
// @returns {error} error, if any
func (c *MockClient) Hash(message []byte) ([]byte, error) {
	return c.fcClient.Hash(message)
}

// InitializeChannel returns nil for tests assuming that give channel is already initialized
// @param {Channel} Channel that needs to be initialized
// @returns {error} error, if any
func (c *MockClient) InitializeChannel(channel sdkApi.Channel) error {
	//return c.fcClient.InitializeChannel(channel)
	return nil
}

// GetConfig get client config
// @returns {Config} config
func (c *MockClient) GetConfig() sdkConfigApi.Config {
	return c.fcClient.GetConfig()
}

// GetSigningIdentity returns the signingIdentity (user) context from the client
// @retruns {sdkApi.IdentityContext} sdkApi.IdentityContext
func (c *MockClient) GetSigningIdentity() sdkApi.IdentityContext {
	return c.fcClient.GetSigningIdentity()
}
