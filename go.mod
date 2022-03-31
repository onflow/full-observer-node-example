module github.com/onflow/full-observer-node-example

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/onflow/flow-go v0.24.10
	github.com/onflow/flow-go/crypto v0.24.10
)

// crypto needs to be locally installed and compiled
replace github.com/onflow/flow-go => ./flow-go

replace github.com/onflow/flow-go/crypto => ./flow-go/crypto
