module github.com/onflow/full-observer-node-example

go 1.16

require (
	github.com/onflow/flow-go v0.22.10
	github.com/onflow/flow-go/crypto v0.22.10
)

// crypto needs to be locally installed and compiled
replace github.com/onflow/flow-go/crypto => ./flow-go/crypto
