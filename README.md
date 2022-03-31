# Setting Up a Consensus Follower

A full observer node follows the Flow protocol consensus to track new finalized blocks, which is accomplished using the consensus follower library included in [github.com/onflow/flow-go/follower](https://github.com/onflow/flow-go/blob/master/follower/consensus_follower.go). This repo provides a working example of how to use the consensus follower library.

There are other features that make up a full observer node, namely, tracking Flow execution state. For an example of a complete full observer node, see [DPS](https://github.com/optakt/flow-dps).

## Fetching genesis data

### Download snapshot data

##### For testnet
```
mkdir -p bootstrap/public-root-information
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap/root-protocol-state-snapshot.json
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap/root-protocol-state-snapshot.json.asc
```

##### For mainnet

```
mkdir -p bootstrap/public-root-information
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap_mainnet/root-protocol-state-snapshot.json
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap_mainnet/root-protocol-state-snapshot.json.asc
```

### Verify the PGP signature
If this is the first time going through this process, add the `flow-signer@onflow.org` public key.
```
gpg --keyserver keys.openpgp.org --search-keys flow-signer@onflow.org
```

Verify the snapshot file
```
gpg --verify bootstrap/public-root-information/root-protocol-state-snapshot.json.asc
```

```
gpg: assuming signed data in 'bootstrap/public-root-information/root-protocol-state-snapshot.json'
gpg: Signature made Wed Sep 15 11:34:33 2021 PDT
gpg:                using ECDSA key 40CD95717AC463E61EE3B285B718CA310EDB542F
gpg: Good signature from "Flow Team (Flow Full Observer node snapshot verification master key) <flow-signer@onflow.org>" [unknown]
gpg: WARNING: This key is not certified with a trusted signature!
gpg:          There is no indication that the signature belongs to the owner.
Primary key fingerprint: 7D23 8D1A E6D3 2A71 8ECD  8611 CB52 64F7 FD4C DD27
     Subkey fingerprint: 40CD 9571 7AC4 63E6 1EE3  B285 B718 CA31 0EDB 542F
```

## Staked Access Nodes

Below are access nodes that support connections from observer nodes for each network.

Note: while the public keys remain the same, the hostnames change each spork to include the spork name. Substitute `[devnet_spork]` and `[mainnet_spork]` with the appropriate spork name (e.g. `mainnet14`). See [Past Sporks](https://docs.onflow.org/node-operation/past-sporks/) for the current spork for each network.

##### For testnet

Access-003:
* Host: `access-003.[devnet_spork].nodes.onflow.org`
* Public Key: `b662102f4184fc1caeb2933cf87bba75cdd37758926584c0ce8a90549bb12ee0f9115111bbbb6acc2b889461208533369a91e8321eaf6bcb871a788ddd6bfbf7`

Access-004:
* Host: `access-004.[devnet_spork].nodes.onflow.org`
* Public Key: `0d1523612be854638b985fc658740fa55f009f3cd49b739961ab082dc91b178ed781ef5f66878613b4d34672039150abfd9c8cfdfe48c565bca053fa4db30bec`

##### For mainnet

Access-007:
* Host: `access-007.[mainnet_spork].nodes.onflow.org`
* Public Key: `28a0d9edd0de3f15866dfe4aea1560c4504fe313fc6ca3f63a63e4f98d0e295144692a58ebe7f7894349198613f65b2d960abf99ec2625e247b1c78ba5bf2eae`

Access-008:
* Host: `access-008.[mainnet_spork].nodes.onflow.org`
* Public Key: `11742552d21ac93da37ccda09661792977e2ca548a3b26d05f22a51ae1d99b9b75c8a9b3b40b38206b38951e98e4d145f0010f8942fd82ddf0fb1d670202264a`

## Launching your node

### Build
The consensus follower requires the `crypto` modules from `onflow/flow-go`, which must be built locally.

Clone the `flow-go` repo
```
git clone https://github.com/onflow/flow-go.git
```

Follow the instructions in the `flow-go` [README](https://github.com/onflow/flow-go/blob/master/README.md) up to `make install-tools`.

Once `crypto` has been installed, build the observer node.
```
go build --tags relic -o observer main.go
```

### Launch

The following are examples of how to launch observer nodes following staked Access Nodes in different networks.

##### For testnet

Launch the observer node following `access-003.devnet33.nodes.onflow.org` on testnet
```
mkdir /tmp/data

export DEVNET_SPORK=devnet33
export ACCESS_NODE_HOSTNAME=access-003.${DEVNET_SPORK}.nodes.onflow.org
export ACCESS_NODE_NETWORKING_PUBLIC_KEY=b662102f4184fc1caeb2933cf87bba75cdd37758926584c0ce8a90549bb12ee0f9115111bbbb6acc2b889461208533369a91e8321eaf6bcb871a788ddd6bfbf7
./observer
```

##### For mainnet

Launch the observer node following `access-008.mainnet16.nodes.onflow.org` on mainnet
```
mkdir /tmp/data

export MAINNET_SPORK=mainnet16
export ACCESS_NODE_HOSTNAME=access-008.${MAINNET_SPORK}.nodes.onflow.org
export ACCESS_NODE_NETWORKING_PUBLIC_KEY=11742552d21ac93da37ccda09661792977e2ca548a3b26d05f22a51ae1d99b9b75c8a9b3b40b38206b38951e98e4d145f0010f8942fd82ddf0fb1d670202264a
./observer
```
