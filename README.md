# SRTP-ChainCode
ZJU SRTP ChainCode

### 数据结构
```
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
```


### APIs
+ GetID(ctx contractapi.TransactionContextInterface, privateCollectionName string) (string, error)
  + 获取当前private collection有多少数据

+ SetPrivateData(ctx contractapi.TransactionContextInterface, fileBytes []byte, orgCollectionName string, id string, uploadTime string) error
  + 将数据存入私有数据库中

+ SetPublicData(ctx contractapi.TransactionContextInterface, pid string, fileBytes []byte, uploadTime string) error
  + 将数据存入公有数据库中

+ DeletePrivateData(ctx contractapi.TransactionContextInterface, privateCollectionName string, id string) error
  + 删除私有数据


+ DeletePublicData(ctx contractapi.TransactionContextInterface, pid string) error
  + 删除公有数据

+ CommitTransaction(ctx contractapi.TransactionContextInterface, pid string, to string, hash string)
  + 提交交易
