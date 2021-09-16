package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/onflow/flow-go/cmd/bootstrap/utils"
	"github.com/onflow/flow-go/follower"
	"github.com/onflow/flow-go/model/flow"

	// crypto needs to be locally installed and compiled
	"github.com/onflow/flow-go/crypto"
)

const (

	// directory to store chain state
	defaultDataDir = "/tmp/data"

	// downloaded genesis data
	defaultBootstrapDir = "./bootstrap"

	// consensus followers own address
	defaultLocalBindAddr = "0.0.0.0:0"

	// upstream access node (bootstrap peer)
	defaultAccessNodeHostname            = "access-001.canary8.nodes.onflow.org"
	defaultAccessNodeLibp2pPort          = 3569
	defaultAccessNodeNetworkingPublicKey = "210c5aae4b72feb1ee16a84d165e386350e4556e731c132f729c761fbf686f73c4927cbcbb22e464da67354a3ce647673bdb733dca129e3df96d82fb8ae94c00"
)

func main() {

	var err error
	dataDir := getEnv("DATA_DIR", defaultDataDir)
	bootstrapDir := getEnv("BOOTSTRAP_DIR", defaultBootstrapDir)
	localBindAddr := getEnv("LOCAL_BIND_ADDRESS", defaultLocalBindAddr)
	accessNodeHostname := getEnv("ACCESS_NODE_HOSTNAME", defaultAccessNodeHostname)
	accessNodeNetworkingPublicKey := getEnv("ACCESS_NODE_NETWORKING_PUBLIC_KEY", defaultAccessNodeNetworkingPublicKey)

	accessNodeLibp2pPort := defaultAccessNodeLibp2pPort
	port := getEnv("ACCESS_NODE_LIBP2P_PORT", "")
	if port != "" {
		accessNodeLibp2pPort, err = strconv.Atoi(port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid value for ACCESS_NODE_LIBP2P_PORT %v\n", port)
			fatalError(err)
		}
	}

	// generate the networking key
	myKey, err := generateKey()
	if err != nil {
		fatalError(err)
	}

	// upstream bootstrap peer public key
	bootstrapPeerKey, err := decodePublicKey(accessNodeNetworkingPublicKey)
	if err != nil {
		fatalError(err)
	}

	bootstrapNodeInfo := follower.BootstrapNodeInfo{
		Host:             accessNodeHostname,
		Port:             uint(accessNodeLibp2pPort),
		NetworkPublicKey: bootstrapPeerKey,
	}

	opts := []follower.Option{
		follower.WithDataDir(dataDir),
		follower.WithBootstrapDir(bootstrapDir),
	}

	// initialize the consensus follower
	cf, err := follower.NewConsensusFollower(myKey, localBindAddr, []follower.BootstrapNodeInfo{bootstrapNodeInfo}, opts...)
	if err != nil {
		fatalError(err)
	}

	cf.AddOnBlockFinalizedConsumer(OnBlockFinalizedConsumer)

	// initialize signal catcher
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer close(sig)

	// run the observer node
	ctx, cancel := context.WithCancel(context.Background())
	go cf.Run(ctx)

	<-sig

	cancel()
}

// decode a public key from a hex string
func decodePublicKey(hexKey string) (crypto.PublicKey, error) {
	bz, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	return crypto.DecodePublicKey(crypto.ECDSAP256, bz)
}

// generate a new network key
func generateKey() (crypto.PrivateKey, error) {
	seed := make([]byte, crypto.KeyGenSeedMinLenECDSASecp256k1)
	n, err := rand.Read(seed)
	if err != nil || n != crypto.KeyGenSeedMinLenECDSASecp256k1 {
		return nil, err
	}
	return utils.GenerateUnstakedNetworkingKey(seedFixture(n))
}

func seedFixture(n int) []byte {
	var seed = make([]byte, n)
	_, _ = rand.Read(seed)
	return seed
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) > 0 {
		return value
	}
	return fallback
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

// handle new finalized block events
func OnBlockFinalizedConsumer(finalizedBlockID flow.Identifier) {
	fmt.Printf(">>>>>>>>>>>>>>>>>>>> Received finalized block: %s\n", finalizedBlockID.String())
}
