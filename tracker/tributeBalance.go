package tracker

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	berryCommon "github.com/berry-data/BerryMiner/common"
	"github.com/berry-data/BerryMiner/config"
	berry "github.com/berry-data/BerryMiner/contracts"
	"github.com/berry-data/BerryMiner/db"
	"github.com/berry-data/BerryMiner/rpc"
)

type TributeTracker struct {
}

func (b *TributeTracker) String() string {
	return "TributeTracker"
}

func (b *TributeTracker) Exec(ctx context.Context) error {
	//cast client using type assertion since context holds generic interface{}
	client := ctx.Value(berryCommon.ClientContextKey).(rpc.ETHClient)
	DB := ctx.Value(berryCommon.DBContextKey).(db.DB)

	//get the single config instance
	cfg := config.GetConfig()

	//get address from config
	_fromAddress := cfg.PublicAddress

	//convert to address
	fromAddress := common.HexToAddress(_fromAddress)

	_conAddress := cfg.ContractAddress

	//convert to address
	contractAddress := common.HexToAddress(_conAddress)

	instance, err := berry.NewBerryMaster(contractAddress, client)
	if err != nil {
		fmt.Println("Instance error - TributeBalance")
		return err
	}

	balance, err := instance.BalanceOf(nil, fromAddress)
	balanceInTributes, ok := big.NewFloat(1).SetString(balance.String())
	if !ok {
		fmt.Println("Problem converting tributes.")
		balanceInTributes = big.NewFloat(0)
	}
	decimals, ok := big.NewFloat(1).SetString("1000000000000000000")
	if !ok {
		fmt.Println("Could not create tribute float for computing tributes")
		balanceInTributes = big.NewFloat(0)
	}
	if decimals != nil {
		balanceInTributes = balanceInTributes.Quo(balanceInTributes, decimals)
	}

	//numTributes, _ := balanceInTributes.Float64()
	log.Printf("Tribute Balance: %v (%v tributes)\n", balance, balanceInTributes)
	if err != nil {
		fmt.Println("Balance Retrieval Error - Tribute Balance")
		return err
	}
	enc := hexutil.EncodeBig(balance)
	return DB.Put(db.TributeBalanceKey, []byte(enc))
}
