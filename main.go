package main

import (
	"fmt"
	"ldk-node-go/ldk_node"
)

func main() {
	// mnemonic := ldk_node.GenerateEntropyMnemonic()
	mnemonic := "little toss recycle anchor festival half marine fog life twenty attract hotel gravity crash blossom three town nut term huge start omit love mixed"
	listenAddr := ldk_node.NetAddress("localhost:9735")

	// builder := ldk_node.NewBuilder()
	config := ldk_node.Config{
		StorageDirPath:                 "./",
		Network:                        ldk_node.NetworkTestnet,
		ListeningAddress:               &listenAddr,
		OnchainWalletSyncIntervalSecs:  30,
		WalletSyncIntervalSecs:         30,
		FeeRateCacheUpdateIntervalSecs: 100,
		LogLevel:                       ldk_node.LogLevelDebug,
		DefaultCltvExpiryDelta:         80,
	}
	builder := ldk_node.BuilderFromConfig(config)
	// builder.SetNetwork(ldk_node.NetworkTestnet)
	builder.SetEsploraServer("https://blockstream.info/testnet/api")
	builder.SetGossipSourceRgs("https://rapidsync.lightningdevkit.org/testnet/snapshot")
	builder.SetEntropyBip39Mnemonic(mnemonic, nil)

	node, err := builder.Build()
	if err != nil {
		fmt.Printf("builder.Build: %v\n", err)
		return
	}

	node.Start()
	node.SyncWallets()

	fundingAddress, err := node.NewOnchainAddress()
	if err != nil {
		fmt.Printf("node.NewOnchainAddress: %v\n", err)
		return
	}
	fmt.Printf("fundingAddress: %s\n", fundingAddress)

	// .. fund address ..

	// nodeID := ldk_node.PublicKey("")
	// nodeAddr := ldk_node.NetAddress("localhost:9735")
	// node.ConnectOpenChannel(nodeID, nodeAddr, 20000, nil, nil, false)
}
