package node

import (
	"fmt"
	"os"

	ankrconfig "github.com/Ankr-network/ankr-chain/config"
	"github.com/Ankr-network/ankr-chain/consensus"
	ankrp2p "github.com/Ankr-network/ankr-chain/p2p"
	"github.com/Ankr-network/ankr-chain/store/historystore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	bc "github.com/tendermint/tendermint/blockchain"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmcorelog "github.com/tendermint/tendermint/libs/log"
	tmcorenode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
)

type AnkrNode struct {
	Name string
	Log  tmcorelog.Logger
	Node *tmcorenode.Node
}

type AnkrNodeProvider func(*ankrconfig.AnkrConfig, tmcorelog.Logger) (*AnkrNode, error)

func NewAnkrNode(config *ankrconfig.AnkrConfig, logger tmcorelog.Logger) (*AnkrNode, error) {
	// Generate node PrivKey
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	// Convert old PrivValidator if it exists.
	oldPrivVal := config.OldPrivValidatorFile()
	newPrivValKey := config.PrivValidatorKeyFile()
	newPrivValState := config.PrivValidatorStateFile()
	if _, err := os.Stat(oldPrivVal); !os.IsNotExist(err) {
		oldPV, err := privval.LoadOldFilePV(oldPrivVal)
		if err != nil {
			return nil, fmt.Errorf("Error reading OldPrivValidator from %v: %v\n", oldPrivVal, err)
		}
		logger.Info("Upgrading PrivValidator file",
			"old", oldPrivVal,
			"newKey", newPrivValKey,
			"newState", newPrivValState,
		)
		oldPV.Upgrade(newPrivValKey, newPrivValState)
	}

	ankrChainApp := ankrchain.NewAnkrChainApplication(config.DBDir(), config.DBBackend, ankrtypes.APPName, logger.With("tx", "AnkrChainApp"))

	config.FilterPeers = config.AllowedPeers != ""


	blockStoreDB := dbm.NewDB("blockstore", dbm.DBBackendType(config.TendermintCoreConfig().DBBackend), config.TendermintCoreConfig().DBDir())
	blockStore := bc.NewBlockStore(blockStoreDB)

	stateDB := dbm.NewDB("state", dbm.DBBackendType(config.TendermintCoreConfig().DBBackend), config.TendermintCoreConfig().DBDir())
	state := sm.LoadState(stateDB)
	if state.LastBlockHeight == 1002800 {
		state.AppHash = blockStore.LoadBlock(1002800 + 1).AppHash
	}

	sm.SaveState(stateDB, state)
	stateDB.Close()
	blockStoreDB.Close()

	tmNode, err :=  tmcorenode.NewNode(config.TendermintCoreConfig(),
		privval.LoadOrGenFilePV(newPrivValKey, newPrivValState),
		nodeKey,
		proxy.NewLocalClientCreator(ankrChainApp),
		tmcorenode.DefaultGenesisDocProviderFunc(config.TendermintCoreConfig()),
		tmcorenode.DefaultDBProvider,
		tmcorenode.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)

	if err != nil {
		return nil, err
	}

	historyDBLogger := logger.With("tx", "historydb")
	historyDBLogger.Info("historydb parameter", "dbType", config.HistoryDB.Type, "dbHost", config.HistoryDB.Host, "dbName", config.HistoryDB.Name)
	if config.HistoryDB.Type != "" && config.HistoryDB.Host != "" && config.HistoryDB.Name != "" {
		historyDBService := historystore.NewHistoryStorageService(config.HistoryDB.Type, config.HistoryDB.Host, config.HistoryDB.Name, tmNode.EventBus(), historyDBLogger)
		historyDBService.Start()
	}

	peerFilter := ankrp2p.NewPeerFilter()
	sd         := ankrp2p.NewSeeds()
	peerFilter.Config(config.AllowedPeers)
	sd.Config(config.P2P.Seeds)
	ankrp2p.Init(peerFilter, sd)

	return &AnkrNode{"", logger, tmNode}, err
}

