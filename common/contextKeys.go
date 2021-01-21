package common

import (
	"github.com/berry-data/BerryMiner/util"
)

var (
	//ClientContextKey is the key used to set the eth client on tracker contexts
	ClientContextKey = util.NewKey("common", "ETHClient")

	//DBContextKey is the shared context key where a DB instance can be found in a context
	DBContextKey = util.NewKey("common", "DB")

	//Berry Contract Address
	ContractAddress = util.NewKey("common", "contractAddress")

	//MasterContractContextKey is the shared context key to get shared master berry contract instance
	MasterContractContextKey = util.NewKey("common", "masterContract")

	//TransactorContractContextKey is the shared context key to get shared transactor berry contract instance
	TransactorContractContextKey = util.NewKey("common", "transactorContract")

	//DataProxyKey used to access the local or remote data server proxy
	DataProxyKey = util.NewKey("common", "DataServerProxy")

	//Ethereum wallet private key
	PrivateKey = util.NewKey("common", "PrivateKey")

	//Ethereum wallet public address
	PublicAddress = util.NewKey("common", "PublicAddress")

)
