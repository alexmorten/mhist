# ppp-mhist
## simple measurement history logger
This is a very simple measurement database, that receives measurements (consisting of name, value and optionally a timestamp) through tcp or http. If you don't send a timestamp with the measurement, the current time is used (there are rarely reasons to send a different timestamp).
The latest measurements are stored in memory for fast access and all measurements are also stored on disk for permanent storage.

For realtime updates you can subscribe to mhist with tcp and for historical access you can retrieve measurements with http.

Mhist also supports barebones data-replication to other instances of itself (the adresses of which have to be known beforehand).

### assumptions
- measurements are received by mhist in the order they are generated
- there are only two types of measurements: `numerical`, sent to mhist as numbers, and `categorical`, sent to mhist as strings
- measurement types don't change for a certain measurement name.
- measurements are taken in regular intervals.
- it is known in advance how much memory and diskspace can be used by mhist.
- for 'recent' measurements really fast access is more important than perfect accuracy.
- when retrieving measurements you want to retrieve measurements of all names more often than just certain names.

### setup

assuming you have a working go installation:
clone this repo into `$GOPATH/src/github.com/alexmorten/mhist`.

Run `make install dep run` in the `ppp-profiler` directory.

To see how to change the default configuration, run ` go run main/main.go -h`

## endpoints

- `/`
  - `POST` send measurement to mhist as json with `name: string` and `value: number|string`.
  - `GET` get recorded measurements with the following optional query params:
    - `start` & `end` points in time as unix-timestamps in nanoseconds, defining what timestamp of measurements to filter for.
    - `granularity` minimum [duration](https://golang.org/pkg/time/#ParseDuration) between measurements (i.e. with a granularity of `1s` all measurements returned will have at least 1 second between them)
    - `names` comma separated list of names of measurements. Measurements that are not in the list will not be returned
- `/meta` get a list of stored measurement names and their types.

### todos

- [ ] refactor package layout. Types used all over the place should be in a package like `models` instead of being defined in the `mhist` package.
- [ ] historical access should be also possible over tcp. For example streaming all measurements starting from a certain timestamp. This would also enable:
- [ ] starting a new mhist instance that grabs all data that was already received by another instance and also gets all replicated date from that point forward.
