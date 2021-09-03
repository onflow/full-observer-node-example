package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/onflow/flow-go/cmd/bootstrap/utils"
	"github.com/onflow/flow-go/follower"
	"github.com/onflow/flow-go/model/encodable"
	"github.com/onflow/flow-go/model/flow"

	// crypto needs to be locally installed and compiled
	"github.com/onflow/flow-go/crypto"
)

// downloaded genesis data
// mkdir bootstrap
// gsutil -m cp -r "gs://flow-genesis-bootstrap/mainnet-12/*" bootstrap/.
const bootstrapDir = "./bootstrap"

// directory to store chain state
const dataDir = "/tmp/data"

// consensus followers own address
const localBindAddr = "0.0.0.0:0"

// upstream access node (bootstrap peer)
const accessNodeHostname = "access-001.canary8.nodes.onflow.org"
const accessNodeLibp2pPort = 3569
const accessNodeNetworkingPublicKey = "\"210c5aae4b72feb1ee16a84d165e386350e4556e731c132f729c761fbf686f73c4927cbcbb22e464da67354a3ce647673bdb733dca129e3df96d82fb8ae94c00\""

func main() {

	// generate a key
	myKey, err := generateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// upstream bootstrap peer public key (dervied from accessNodeNetworkingPublicKey string)
	var bootstrapPeerKey encodable.NetworkPubKey
	err = json.Unmarshal([]byte(accessNodeNetworkingPublicKey), &bootstrapPeerKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	bootstrapNodeInfo := follower.BootstrapNodeInfo{
		Host:             accessNodeHostname,
		Port:             accessNodeLibp2pPort,
		NetworkPublicKey: bootstrapPeerKey.PublicKey,
	}

	opts := []follower.Option{
		follower.WithDataDir(dataDir),
		follower.WithBootstrapDir(bootstrapDir),
	}

	follower, err := follower.NewConsensusFollower(myKey, localBindAddr, []follower.BootstrapNodeInfo{bootstrapNodeInfo}, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	follower.AddOnBlockFinalizedConsumer(OnBlockFinalizedConsumer)

	ctx, cancel := context.WithCancel(context.Background())
	go follower.Run(ctx)

	time.Sleep(5 * time.Minute)
	cancel()
}

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

func OnBlockFinalizedConsumer(finalizedBlockID flow.Identifier) {
	fmt.Printf(">>>>>>>>>>>>>>>>>>>> Received finalized block: %s\n", finalizedBlockID.String())
}
