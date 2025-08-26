package contracts

import (
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/the-web3/dapplink-vrf/bindings"
	"github.com/the-web3/dapplink-vrf/database"
	"github.com/the-web3/dapplink-vrf/database/event"
	"github.com/the-web3/dapplink-vrf/database/worker"
)

type DappLinkVrf struct {
	DlVrfAbi    *abi.ABI
	DlVrfFilter *bindings.DappLinkVRFFilterer
}

func NewDappLinkVrf() (*DappLinkVrf, error) {
	dappLinkVrfAbi, err := bindings.DappLinkVRFMetaData.GetAbi()
	if err != nil {
		log.Error("get dapplink vrf abi fail", "err", err)
		return nil, err
	}
	dappLinkVRFFilterer, err := bindings.NewDappLinkVRFFilterer(common.Address{}, nil)
	if err != nil {
		log.Error("new dapplink vrf filter fail", "err", err)
		return nil, err
	}
	return &DappLinkVrf{
		DlVrfAbi:    dappLinkVrfAbi,
		DlVrfFilter: dappLinkVRFFilterer,
	}, nil
}

func (dvf *DappLinkVrf) ProcessDappLinkVrfEvent(db *database.DB, dappLinkVrfAddres string, startHeight, endHeight *big.Int) ([]worker.RequestSend, []worker.FillRandomWords, error) {
	var RequestSentList []worker.RequestSend
	var FillRandomWordList []worker.FillRandomWords
	contactFilter := event.ContractEvent{ContractAddress: common.HexToAddress(dappLinkVrfAddres)}
	contractEventList, err := db.ContractEvent.ContractEventsWithFilter(contactFilter, startHeight, endHeight)
	if err != nil {
		log.Error("query contacts event fail", "err", err)
		return RequestSentList, FillRandomWordList, err
	}

	for _, contractEvent := range contractEventList {
		if contractEvent.EventSignature.String() == dvf.DlVrfAbi.Events["RequestSent"].ID.String() {
			rquestSentEvent, err := dvf.DlVrfFilter.ParseRequestSent(*contractEvent.RLPLog)
			if err != nil {
				log.Error("parse request sent fail", "err", err)
				return RequestSentList, FillRandomWordList, err
			}
			log.Info("Request sent event", "RequestId", rquestSentEvent.RequestId, "NumWords", rquestSentEvent.NumWords, "Current", rquestSentEvent.Current)
			rs := worker.RequestSend{
				GUID:       uuid.New(),
				RequestId:  rquestSentEvent.RequestId,
				VrfAddress: rquestSentEvent.Current,
				NumWords:   rquestSentEvent.NumWords,
				Status:     0,
				Timestamp:  uint64(time.Now().Unix()),
			}
			RequestSentList = append(RequestSentList, rs)
		}

		if contractEvent.EventSignature.String() == dvf.DlVrfAbi.Events["FillRandomWords"].ID.String() {
			fillRandomWords, err := dvf.DlVrfFilter.ParseFillRandomWords(*contractEvent.RLPLog)
			if err != nil {
				log.Error("parse fill random fail", "err", err)
				return RequestSentList, FillRandomWordList, err
			}
			log.Info("Fill random words event", "RequestId", fillRandomWords.RequestId, "RandomWords", fillRandomWords.RandomWords)
			var randomWords string
			for _, rword := range fillRandomWords.RandomWords {
				randomWords = rword.String()
			}
			frw := worker.FillRandomWords{
				GUID:        uuid.New(),
				RequestId:   fillRandomWords.RequestId,
				RandomWords: randomWords,
				Timestamp:   uint64(time.Now().Unix()),
			}
			FillRandomWordList = append(FillRandomWordList, frw)
		}
	}
	return RequestSentList, FillRandomWordList, nil
}
