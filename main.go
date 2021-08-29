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
const accessNodeHostname = "access-001.canary7.nodes.onflow.org"
const accessNodeLibp2pPort = 3569
const accessNodeNetworkingPublicKey = "\"e0c141fc192cfdced2a3c8d28ab016a31ae3165d0df43e59200187e36ba079d55e7d38cf7a082c7550af716bca27238ef388260e53257802869120f7df2f4dfb\""

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

	time.Sleep(30 * time.Second)
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
	fmt.Printf("Received finalized block: %s\n", finalizedBlockID.String())
}
