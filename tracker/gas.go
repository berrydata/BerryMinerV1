package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	berryCommon "github.com/berry-data/BerryMiner/common"
	"github.com/berry-data/BerryMiner/db"
	"github.com/berry-data/BerryMiner/rpc"
)

const GWEI = 1000000000

type GasTracker struct {
}

//GasPriceModel is what ETHGasStation returns from queries. Not all fields are filled in
type GasPriceModel struct {
	Fast    float32 `json:"fast"`
	Fastest float32 `json:"fastest"`
	Average float32 `json:"average"`
}

func (b *GasTracker) String() string {
	return "GasTracker"
}

func (b *GasTracker) Exec(ctx context.Context) error {
	client := ctx.Value(berryCommon.ClientContextKey).(rpc.ETHClient)
	DB := ctx.Value(berryCommon.DBContextKey).(db.DB)

	netId, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}

	var gasPrice *big.Int

	if big.NewInt(1).Cmp(netId) == 0 {
		url := "https://ethgasstation.info/json/ethgasAPI.json"
		req := &FetchRequest{queryURL: url, timeout: time.Duration(15 * time.Second)}
		payload, err := fetchWithRetries(req)
		if err != nil {
			gasPrice, err = client.SuggestGasPrice(context.Background())
		} else {
			gpModel := GasPriceModel{}
			err = json.Unmarshal(payload, &gpModel)
			if err != nil {
				log.Printf("Problem with ETH gas station json: %v\n", err)
				gasPrice, err = client.SuggestGasPrice(context.Background())
			} else {
				gasPrice = big.NewInt(int64(gpModel.Fast / 10))
				gasPrice = gasPrice.Mul(gasPrice, big.NewInt(GWEI))
				log.Println("Using ETHGasStation fast price: ", gasPrice)
			}
		}
	} else {
		gasPrice, err = client.SuggestGasPrice(context.Background())
	}

	enc := hexutil.EncodeBig(gasPrice)
	return DB.Put(db.GasKey, []byte(enc))
}
