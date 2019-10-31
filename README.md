# mhist

## simple measurement history logger

This is a very simple measurement database, that receives measurements (consisting of name, value and optionally a timestamp) through tcp or http. If you don't send a timestamp with the measurement, the current time is used (there are rarely reasons to send a different timestamp).
Measurements are stored on disk

For realtime updates you can subscribe to mhist.

### assumptions

- measurements are received by mhist in the order they are generated
- there are only two types of measurements: `numerical` and `categorical`
- measurement types don't change for a certain measurement name.
- it is known in advance how much memory and diskspace can be used by mhist.
- when retrieving measurements you want to retrieve measurements of all names more often than just certain names.

### setup

This uses [go mod for dependency management](https://github.com/golang/go/wiki/Modules)

To see how to change the default configuration, run `go run main/main.go -h`

## endpoints

see the [proto definition](proto/rpc.proto)

### todos

- [ ] add tests for subscription logic
- [ ] add raw measurement type, where the value is just bytes (value is written to value log file, position in file and length written to current "index" file)
- [ ] - add in memory file index, minimising the filesystem list calls
