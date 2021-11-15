package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransferRecord struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Returned bool   `json:"Returned"`
}

type Asset struct {
	// 图片id，hash值，拥有者，
	AssetID          string           `json:"AssetID"`
	Hash             string           `json:"Hash"`
	ImageFingerprint string           `json: "ImageFingerprint"`
	Owner            string           `json:"Owner"`
	CurrentHolder    []string         `json:"CurrentHolder"`
	TransferRecords  []TransferRecord `json:"TransferRecords"`
}

type SmartContract struct {
	contractapi.Contract
}

/*
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transfer_records := []TransferRecord{
		{From: "Alice", To: "Bob", Returned: false},
		{From: "Bob", To: "Alice", Returned: true},
	}

	assets := []Asset{
		{AssetID: "1", Hash: "319010xxxx", Owner: "Alice", ImageFingerprint: "13", CurrentHolder: []string{"Alice", "Bob"}, TransferRecords: transfer_records},
		{AssetID: "5", Hash: "319010xxxx", Owner: "Bob", ImageFingerprint: "123", CurrentHolder: []string{"Alice", "Bob"}, TransferRecords: transfer_records},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.AssetID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, imageId string, hash string, imageFingerprint string, owner string) error {
	exists, err := s.AssetExists(ctx, imageId)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("the asset %s already exists", imageId)
	}

	asset := Asset{
		AssetID:         imageId,
		Hash:            hash,
		Owner:           owner,
		CurrentHolder:   []string{owner},
		TransferRecords: []TransferRecord{},
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(imageId, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, imageId string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(imageId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", imageId)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, imageId string, from string, to string) error {
	asset, err := s.ReadAsset(ctx, imageId)
	if err != nil {
		return err
	}

	if asset.Owner != from {
		return fmt.Errorf("the asset does not belong to %s", from)
	}

	if stringInArray(to, asset.CurrentHolder) {
		return fmt.Errorf("the asset has already belonged to %s", to)
	}

	transferRecord := TransferRecord{
		From:     from,
		To:       to,
		Returned: false,
	}

	asset.CurrentHolder = append(asset.CurrentHolder, to)
	asset.TransferRecords = append(asset.TransferRecords, transferRecord)

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(imageId, assetJSON)
}

func (s *SmartContract) ReturnAsset(ctx contractapi.TransactionContextInterface, imageId string, from string, to string) error {
	asset, err := s.ReadAsset(ctx, imageId)
	if err != nil {
		return err
	}

	transferRecordIdx, err := relationInArray(from, to, asset.TransferRecords)
	if err != nil {
		return err
	}

	asset.TransferRecords[transferRecordIdx].Returned = true
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(imageId, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, imageId string) error {
	exists, err := s.AssetExists(ctx, imageId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", imageId)
	}

	return ctx.GetStub().DelState(imageId)
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, imageId string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(imageId)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func relationInArray(from string, to string, transferRecords []TransferRecord) (int, error) {
	for idx, transferRecord := range transferRecords {
		if transferRecord.From == from && transferRecord.To == to {
			return idx, nil
		}
	}
	return -1, fmt.Errorf("no such relation exist")
}

func stringInArray(targetStr string, stringArray []string) bool {
	for _, str := range stringArray {
		if str == targetStr {
			return true
		}
	}
	return false
}

*/

type PrivateAssetDetails struct {
	ImageID    string `json:"imageID"`
	Content    []byte `json:"content"`
	UploadTime string `json:"uploadTime"`
}

type PublicAssetDetails struct {
	Pid        string `json:"pid"`
	Content    []byte `json:"content"`
	UploadTime string `json:"uploadTime"`
}

type TransactionDetails struct {
	TxnID string `json:"txnID"`
	Pid   string `json:"pid"`
	To    string `json:"to"`
	Hash  string `json:"hash"`
}

func (s *SmartContract) CommitTransaction(ctx contractapi.TransactionContextInterface, pid string, to string, hash string) error {
	transactions, err := s.GetTransactionDataByRange(ctx, "txn", "txp")
	if err != nil {
		return fmt.Errorf("CommitTransaction cannot be performed: Error %v", err)
	}
	txn_id := strconv.Itoa(len(transactions))
	transaction := TransactionDetails{
		TxnID: "txn" + txn_id,
		Pid:   pid,
		To:    to,
		Hash:  hash,
	}

	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}

	JSONBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}
	err = ctx.GetStub().PutState("txn"+txn_id, JSONBytes)
	if err != nil {
		return fmt.Errorf("failed to put transaction details: %v", err)
	}

	return nil
}

func (s *SmartContract) SetPrivateData(ctx contractapi.TransactionContextInterface, fileBytes []byte, orgCollectionName string, id string, uploadTime string) error {

	privateAsset := PrivateAssetDetails{
		ImageID:    id,
		Content:    fileBytes,
		UploadTime: uploadTime,
	}
	JSONBytes, err := json.Marshal(privateAsset)

	// orgCollectionName := "Org1"
	// Save asset details to collection visible to owning organization
	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}
	// TODO 什么时候处理id++
	err = ctx.GetStub().PutPrivateData(orgCollectionName, id, JSONBytes)
	if err != nil {
		return fmt.Errorf("failed to put asset private details: %v", err)
	}

	return nil
}

func (s *SmartContract) GetPrivateDataByRange(ctx contractapi.TransactionContextInterface, privateCollectionName string, start string, end string) ([]*PrivateAssetDetails, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(privateCollectionName, start, end)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()
	var privateAssets []*PrivateAssetDetails
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var privateAsset PrivateAssetDetails
		err = json.Unmarshal(queryResponse.Value, &privateAsset)
		if err != nil {
			return nil, err
		}
		privateAssets = append(privateAssets, &privateAsset)
	}

	return privateAssets, nil
}

func (s *SmartContract) GetID(ctx contractapi.TransactionContextInterface, privateCollectionName string) (string, error) {
	allPrivateData, err := s.GetPrivateDataByRange(ctx, privateCollectionName, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}

	return strconv.Itoa(len(allPrivateData)), nil
}

func (s *SmartContract) SetPublicData(ctx contractapi.TransactionContextInterface, pid string, fileBytes []byte, uploadTime string) error {

	publicAsset := PublicAssetDetails{
		Pid:        pid,
		Content:    fileBytes,
		UploadTime: uploadTime,
	}
	publicAssetJSON, err := json.Marshal(publicAsset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}

	err = ctx.GetStub().PutState(pid, publicAssetJSON)
	if err != nil {
		return fmt.Errorf("SetPublicImage cannot be performed: Error %v", err)
	}

	return nil
}

func (s *SmartContract) DeletePrivateData(ctx contractapi.TransactionContextInterface, privateCollectionName string, id string) error {

	err := ctx.GetStub().DelPrivateData(privateCollectionName, id)
	if err != nil {
		return fmt.Errorf("DeletePrivateData cannot be performed: Error %v", err)
	}

	return nil
}

func (s *SmartContract) DeletePublicData(ctx contractapi.TransactionContextInterface, pid string) error {

	err := ctx.GetStub().DelState(pid)
	if err != nil {
		return fmt.Errorf("DeletePublicData cannot be performed: Error %v", err)
	}
	return nil
}

func (s *SmartContract) GetPublicDataByRange(ctx contractapi.TransactionContextInterface, start string, end string) ([]*PublicAssetDetails, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange(start, end)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var publicAssets []*PublicAssetDetails
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var publicAsset PublicAssetDetails
		err = json.Unmarshal(queryResponse.Value, &publicAsset)
		if err != nil {
			return nil, err
		}
		publicAssets = append(publicAssets, &publicAsset)
	}

	return publicAssets, nil
}

func (s *SmartContract) GetTransactionDataByRange(ctx contractapi.TransactionContextInterface, start string, end string) ([]*TransactionDetails, error) {
	// 应当以txn开头

	resultsIterator, err := ctx.GetStub().GetStateByRange(start, end)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*TransactionDetails
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction TransactionDetails
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}
