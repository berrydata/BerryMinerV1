package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	berryCommon "github.com/berry-data/BerryMiner/common"
	"github.com/berry-data/BerryMiner/contracts"
	"github.com/berry-data/BerryMiner/db"
	"github.com/berry-data/BerryMiner/util"
)

var zero = big.NewInt(0)
var top50Logger = util.NewLogger("tracker", "Top50Tracker")

//Top50Tracker concrete tracker type
type Top50Tracker struct {
}

func (b *Top50Tracker) String() string {
	return "Top50Tracker"
}

//Exec implementation for tracker
func (b *Top50Tracker) Exec(ctx context.Context) error {

	//cast client using type assertion since context holds generic interface{}
	DB := ctx.Value(berryCommon.DBContextKey).(db.DB)

	instance := ctx.Value(berryCommon.MasterContractContextKey).(*contracts.BerryMaster)

	top50Logger.Debug("Querying for top50 request IDs...")
	top50, err := instance.GetRequestQ(nil)
	if err != nil {
		fmt.Println("top50 get error")
		return err
	}
	rIDs := []byte{}
	for i := range top50 {
		reqID := top50[i]
		if reqID.Cmp(zero) == 0 {
			//top50Logger.Info("Skipping zero-value request id")
			continue
		}
		queryMetadata, err := DB.Get(fmt.Sprintf("%s%d", db.QueryMetadataPrefix, reqID.Uint64()))
		if err != nil {
			top50Logger.Error("Problem reading query meta from DB: %v\n", err)
			continue
		}

		if queryMetadata == nil || len(queryMetadata) == 0 {
			var meta *IDSpecifications

			//did not find meta stored locally, let's grab it from on-chain then
			top50Logger.Debug("Pulling query metadata from on-chain with id: %v\n", reqID)
			meta, err = GetSpecs(ctx, uint(reqID.Uint64()))
			if err != nil {
				top50Logger.Error("Problem pulling query metadata from on-chain: %v\n", err)
				continue
			}

			if meta == nil {
				top50Logger.Error("Could not resolve request with id: %v\n", reqID)
				continue
			}
			top50Logger.Debug("Returned query meta: %+v\n", meta)
			jBytes, err := json.Marshal(meta)
			if err != nil {
				top50Logger.Error("Problem marshalling query metadata", err)
				continue
			}
			err = DB.Put(fmt.Sprintf("%s%d", db.QueryMetadataPrefix, reqID), jBytes)
			if err != nil {
				top50Logger.Error("Problem storing query metadata", err)
				//but we keep going anyway
			}
		}

		//TODO: retrieve request meta details here if we've never seen ID before
		top50Logger.Debug("Found top50 ID: %v\n", top50[i])
		rIDs = append(rIDs, top50[i].Bytes()...)

	}
	return DB.Put(db.Top50Key, rIDs)
}
