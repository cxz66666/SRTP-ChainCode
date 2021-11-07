package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TransferRecord struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Returned bool   `json:"Returned"`
}

type Asset struct {
	// 图片id，hash值，拥有者，
	AssetID         string           `json:"AssetID"`
	Hash            string           `json:"Hash"`
	Owner           string           `json:"Owner"`
	CurrentHolder   []string         `json:"CurrentHolder"`
	TransferRecords []TransferRecord `json:"TransferRecords"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transfer_records := []TransferRecord{
		{From: "Alice", To: "Bob", Returned: false},
		{From: "Bob", To: "Alice", Returned: true},
	}

	assets := []Asset{
		{AssetID: "1", Hash: "319010xxxx", Owner: "Alice", CurrentHolder: []string{"Alice", "Bob"}, TransferRecords: transfer_records},
		{AssetID: "5", Hash: "319010xxxx", Owner: "Bob", CurrentHolder: []string{"Alice", "Bob"}, TransferRecords: transfer_records},
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

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, imageId string, hash string, owner string) error {
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
