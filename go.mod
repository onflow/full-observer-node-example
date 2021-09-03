module github.com/onflow/unstaked-consensus-follower

go 1.16

require (
	github.com/onflow/flow-go v0.21.1-0.20210903195644-6a670b6fce93
	github.com/onflow/flow-go/crypto v0.18.0
)

// crypto needs to be locally installed and compiled
replace github.com/onflow/flow-go/crypto => ./crypto
