# SRTP-ChainCode
ZJU SRTP ChainCode

### 结构体定义
```
type TransferRecord struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Returned bool   `json:"Returned"`
}

type Asset struct {
	// 图片id，hash值，拥有者，当前可用者，交换记录
	AssetID         string           `json:"AssetID"`
	Hash            string           `json:"Hash"`
	Owner           string           `json:"Owner"`
	CurrentHolder   []string         `json:"CurrentHolder"`
	TransferRecords []TransferRecord `json:"TransferRecords"`
}
```

### 当前实现了以下函数
+ InitLedger 
+ CreateAsset(imageId, hash, owner)
+ ReadAsset(imageId)
+ TransferAsset(imageId, from, to)
+ ReturnAsset(imageId, from, to)
+ DeleteAsset(imageId)
+ GetAllAssets()