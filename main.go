// RUN
// $ LD_LIBRARY_PATH=ldk_node ./ldk-node-go

package main

import (
	"fmt"
	"ldk-node-go/ldk_node"
	"os"
)

const (
	walletPath = "./wallet"
	seedFile   = "./wallet/seed.bin"

	usingNetwork = ldk_node.NetworkTestnet
	// esploraServer = "https://blockstream.info/testnet/api"
	esploraServer = "https://mempool.space/testnet/api/"
	gossipSource  = "https://rapidsync.lightningdevkit.org/testnet/snapshot"
	logPath       = "./logs"
	expiryDelta   = 40
)

// https://docs.rs/ldk-node/0.1.0/ldk_node/all.html
func main() {
	var err error
	var node *ldk_node.LDKNode

	logDir := logPath

	// https://docs.rs/ldk-node/0.1.0/ldk_node/struct.Config.html
	// https://github.com/lightningdevkit/ldk-node/blob/0c137264975e02757cf2b4a17de116d12a8c8296/src/lib.rs#L161
	// // Config defaults
	// const DEFAULT_STORAGE_DIR_PATH: &str = "/tmp/ldk_node/";
	// const DEFAULT_NETWORK: Network = Network::Bitcoin;
	// const DEFAULT_CLTV_EXPIRY_DELTA: u32 = 144;
	// const DEFAULT_BDK_WALLET_SYNC_INTERVAL_SECS: u64 = 80;
	// const DEFAULT_LDK_WALLET_SYNC_INTERVAL_SECS: u64 = 30;
	// const DEFAULT_FEE_RATE_CACHE_UPDATE_INTERVAL_SECS: u64 = 60 * 10;
	// const DEFAULT_PROBING_LIQUIDITY_LIMIT_MULTIPLIER: u64 = 3;
	// const DEFAULT_LOG_LEVEL: LogLevel = LogLevel::Debug;
	config := ldk_node.Config{
		StorageDirPath: walletPath,
		LogDirPath:     &logDir,
		Network:        usingNetwork,
		// ListeningAddress:       &listenAddr,
		DefaultCltvExpiryDelta:        expiryDelta,
		OnchainWalletSyncIntervalSecs: 10,
		WalletSyncIntervalSecs:        30,
		// FeeRateCacheUpdateIntervalSecs:  600,
		// TrustedPeers0conf:               []ldk_node.PublicKey{},
		// ProbingLiquidityLimitMultiplier: 3,
		LogLevel: ldk_node.LogLevelTrace,
	}
	builder := ldk_node.BuilderFromConfig(config)
	builder.SetEsploraServer(esploraServer)
	builder.SetGossipSourceRgs(gossipSource)

	// set entropy
	if _, err = os.Stat(seedFile); err != nil {
		if _, err := os.Stat(walletPath); err == nil {
			fmt.Printf("wallet already exist: %s\n", walletPath)
			return
		}
		fmt.Printf("new wallet\n")
	}
	builder.SetEntropySeedPath(seedFile)
	node, err = builder.Build()
	if err != nil {
		fmt.Printf("builder.Build: %v\n", err)
		return
	}

	// start LDK/BDK
	node.Start()
	node.SyncWallets()

	fmt.Printf("node_id: %s\n\n", node.NodeId())

	// https://1ml.com/testnet/node/02312627fdf07fbdd7e5ddb136611bdde9b00d26821d14d94891395452f67af248
	nodeID := ldk_node.PublicKey("02312627fdf07fbdd7e5ddb136611bdde9b00d26821d14d94891395452f67af248")
	nodeAddr := ldk_node.NetAddress("23.237.77.12:9735")
	err = node.Connect(nodeID, nodeAddr, false)
	if err != nil {
		fmt.Printf("node.Connect: %v\n", err)
		return
	}
	peers := node.ListPeers()
	for i, peer := range peers {
		fmt.Printf("[%d]\n", i)
		fmt.Printf("  node_id: %s\n", peer.NodeId)
		fmt.Printf("  address: %s\n", peer.Address)
		fmt.Printf("  connected: %v\n", peer.IsConnected)
	}
	fmt.Printf("\n")

	channels := node.ListChannels()
	for i, ch := range channels {
		fmt.Printf("[%d]\n", i)
		fmt.Printf("  node_id: %s\n", ch.CounterpartyNodeId)
		fmt.Printf("  channel_point: %s\n", ch.FundingTxo.Txid)
		fmt.Printf("  chan_id: %s\n", ch.ChannelId)
		fmt.Printf("  value: %d sat\n", ch.ChannelValueSats)
		fmt.Printf("  balance: %d msat\n", ch.BalanceMsat)
		fmt.Printf("  outbound_capacity: %d msat\n", ch.OutboundCapacityMsat)
		fmt.Printf("  inbound_capacity: %d msat\n", ch.InboundCapacityMsat)
		fmt.Printf("  confs: %d\n", *ch.Confirmations)
		fmt.Printf("  outbound: %v\n", ch.IsOutbound)
		fmt.Printf("  channel_ready: %v\n", ch.IsChannelReady)
		fmt.Printf("  usable: %v\n", ch.IsUsable)
		fmt.Printf("  public: %v\n", ch.IsPublic)
	}
	fmt.Printf("\n")

	totalAmount, err := node.TotalOnchainBalanceSats()
	if err != nil {
		fmt.Printf("node.TotalOnchainBalanceSats: %v\n", err)
		return
	}
	fmt.Printf("totalAmount: %d\n", totalAmount)

	spendableAmount, err := node.SpendableOnchainBalanceSats()
	if err != nil {
		fmt.Printf("node.SpendableOnchainBalanceSats: %v\n", err)
		return
	}
	fmt.Printf("spendableAmount: %d\n", spendableAmount)

	if spendableAmount < 20000 && len(channels) == 0 {
		fundingAddress, err := node.NewOnchainAddress()
		if err != nil {
			fmt.Printf("node.NewOnchainAddress: %v\n", err)
			return
		}
		fmt.Printf("fundingAddress: %s\n", fundingAddress)
		return
	}

	if len(channels) == 0 {
		err = node.ConnectOpenChannel(nodeID, nodeAddr, 20000, nil, nil, false)
		if err != nil {
			fmt.Printf("node.ConnectOpenChannel: %v\n", err)
			return
		}
		event := node.WaitNextEvent()
		fmt.Printf("event: %v\n", event)
		node.EventHandled()
	}

	fmt.Printf("done.\n")
}
