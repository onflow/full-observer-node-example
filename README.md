# Setting Up a Full Observer Node

The full observer node is implemented using the Consensus Follower library. This is an example implementation which connects to testnet.

## Fetching genesis data

### Download snapshot data
```
mkdir -p bootstrap/public-root-information
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap/root-protocol-state-snapshot.json
wget -P bootstrap/public-root-information https://storage.googleapis.com/flow-genesis-bootstrap/full_observer_bootstrap/root-protocol-state-snapshot.json.asc
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

For testnet-27, this should return the following:
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

The following access nodes support connections from observer nodes:

Access-003:
* Host: `access-003.devnet27.nodes.onflow.org`
* Public Key: `b662102f4184fc1caeb2933cf87bba75cdd37758926584c0ce8a90549bb12ee0f9115111bbbb6acc2b889461208533369a91e8321eaf6bcb871a788ddd6bfbf7`

Access-004:
* Host: `access-004.devnet27.nodes.onflow.org`
* Public Key: `0d1523612be854638b985fc658740fa55f009f3cd49b739961ab082dc91b178ed781ef5f66878613b4d34672039150abfd9c8cfdfe48c565bca053fa4db30bec`

## Launching your node

### Install dependencies
Clone the `flow-go` repo. The `crypto` module must be built to use the libraries.
```
git clone https://github.com/onflow/flow-go.git
```

Follow the instructions in the `flow-go` README up to `make install-tools`.

### Launch
The following will launch a full node following `access-003.devnet27.nodes.onflow.org` on testnet
```
mkdir /tmp/data

go build -o server --tags relic main.go

export ACCESS_NODE_HOSTNAME=access-003.devnet27.nodes.onflow.org
export ACCESS_NODE_NETWORKING_PUBLIC_KEY=b662102f4184fc1caeb2933cf87bba75cdd37758926584c0ce8a90549bb12ee0f9115111bbbb6acc2b889461208533369a91e8321eaf6bcb871a788ddd6bfbf7
./server
```
