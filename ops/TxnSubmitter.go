package ops

import (
	"context"

	berryCommon "github.com/berry-data/BerryMiner/common"
	"github.com/berry-data/BerryMiner/rpc"
	"github.com/berry-data/BerryMiner/db"
)

//TxnSubmitter just concrete type for txn submitter
type TxnSubmitter struct {
}

//NewSubmitter creates a new TxnSubmitter instance
func NewSubmitter() TxnSubmitter {
	return TxnSubmitter{}
}

//PrepareTransaction relies on rpc package to prepare and submit transactions
func (s TxnSubmitter) PrepareTransaction(ctx context.Context,proxy db.DataServerProxy, ctxName string, callback berryCommon.TransactionGeneratorFN) error {
	return rpc.PrepareContractTxn(ctx,proxy, ctxName, callback)
}
