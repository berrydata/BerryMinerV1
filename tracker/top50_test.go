package tracker

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	berryCommon "github.com/berry-data/BerryMiner/common"
	"github.com/berry-data/BerryMiner/config"
	berry "github.com/berry-data/BerryMiner/contracts"
	"github.com/berry-data/BerryMiner/db"
	"github.com/berry-data/BerryMiner/rpc"
)

func TestTop50(t *testing.T) {
	startBal := big.NewInt(356000)
	cfg := config.GetConfig()

	top50 := make([]*big.Int, 51)
	mockQueryParams := &rpc.MockQueryMeta{QueryString: "json(https://api.binance.com/api/v1/klines?symbol=ETHBTC&interval=1d&limit=1).0.4", Granularity: 1000}
	paramsMap := make(map[uint]*rpc.MockQueryMeta)
	for i := range top50 {
		top50[i] = big.NewInt(int64(i + 51))
		paramsMap[uint(i+51)] = mockQueryParams
	}

	opts := &rpc.MockOptions{ETHBalance: startBal, Nonce: 1, GasPrice: big.NewInt(700000000),
		TokenBalance: big.NewInt(0), Top50Requests: top50, QueryMetadata: paramsMap}
	client := rpc.NewMockClientWithValues(opts)

	contractAddress := common.HexToAddress(cfg.ContractAddress)
	masterInstance, err := berry.NewBerryMaster(contractAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	DB, err := db.Open(filepath.Join(os.TempDir(), "test_top50"))
	if err != nil {
		t.Fatal(err)
	}
	tracker := &Top50Tracker{}
	ctx := context.WithValue(context.Background(), berryCommon.ClientContextKey, client)
	ctx = context.WithValue(ctx, berryCommon.DBContextKey, DB)
	ctx = context.WithValue(ctx, berryCommon.MasterContractContextKey, masterInstance)
	err = tracker.Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
	v, err := DB.Get(db.Top50Key)
	fmt.Println("V", v)
	fmt.Println("newV", v[0])
	if err != nil {
		t.Fatal(err)
	}
	bigs := strings.Split(string(v), ",")
	if len(bigs) == 0 {
		t.Fatal("Expected CSV list of top50 ids")
	}
}
