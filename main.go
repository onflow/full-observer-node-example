package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/onflow/flow-go/cmd/bootstrap/utils"
	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/follower"
	"github.com/onflow/flow-go/model/encodable"
	"github.com/onflow/flow-go/model/flow"
)

// mkdir bootstrap
// gsutil -m cp -r "gs://flow-genesis-bootstrap/mainnet-12/" bootstrap

const dataDir = "/tmp/data"
const bootstrapDir = "./bootstrap"
const accessNodeNetworkingPublicKey =""
const bindAddr = "0.0.0.0:0"

func main() {

	// generate a key
	myKey, err := generateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// upstream bootstrap peer public key
	var bootstrapPeerKey encodable.NetworkPubKey
	err = json.Unmarshal([]byte(accessNodeNetworkingPublicKey), &bootstrapPeerKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	bootstrapNodeInfo := follower.BootstrapNodeInfo{
		Host:             "localhost",
		Port:             3569,
		NetworkPublicKey: bootstrapPeerKey.PublicKey,
	}

	opts := []follower.Option{
		follower.WithDataDir(dataDir),
		follower.WithBootstrapDir(bootstrapDir),
	}

	follower, err := follower.NewConsensusFollower(myKey, bindAddr, []follower.BootstrapNodeInfo{bootstrapNodeInfo}, opts...)
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