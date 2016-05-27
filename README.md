# LevelDB implementation in pure Go
This is a LevelDB implementation, written in pure Go.

[![GoDoc](https://godoc.org/github.com/kezhuw/leveldb?status.svg)](http://godoc.org/github.com/kezhuw/leveldb)
[![Build Status](https://travis-ci.org/kezhuw/leveldb.svg?branch=master)](https://travis-ci.org/kezhuw/leveldb)

**This project is in early stage of development, and should not be considered as stable. Don't use it in production!**


## Goals
This project aims to write an pure Go implementation for LevelDB. It conforms to LevelDB storage format, but not
implementation details of official. This means that there may be differences in implementation with official.


## Features

- **Clean interface**  
The only exporting package is the top level package. All other packages are internal, and should not used by clients.

- **Concurrent compactions**  
Currently, one level compaction can coexist with memtable compaction. It is possible to do concurrent level compactions
in the future.


## Benchmarks

See [kezhuw/go-leveldb-benchmarks][go-leveldb-benchmarks].

```shell
# go test -driver kezhuw -bench .
# https://github.com/kezhuw/leveldb
BenchmarkOpen-4                  1000000              2988 ns/op
BenchmarkSeekRandom-4            1000000             30365 ns/op
BenchmarkReadHot-4               1000000              7372 ns/op
BenchmarkReadRandom-4            1000000              8609 ns/op
BenchmarkReadMissing-4           1000000             15768 ns/op
BenchmarkReadReverse-4           1000000              1053 ns/op
BenchmarkReadSequential-4        2000000               947 ns/op
BenchmarkWriteRandom-4            200000              8073 ns/op
BenchmarkWriteSequential-4        200000              7148 ns/op
BenchmarkDeleteRandom-4           200000              9972 ns/op
BenchmarkDeleteSequential-4       200000              7756 ns/op
ok      github.com/kezhuw/go-leveldb-benchmarks 105.298s
```

```shell
# go test -driver cpp -bench .
# https://github.com/google/leveldb  [cgo]
BenchmarkOpen-4                 10000000               234 ns/op
BenchmarkSeekRandom-4            1000000             47884 ns/op
BenchmarkReadHot-4                500000              4914 ns/op
BenchmarkReadRandom-4             500000              8081 ns/op
BenchmarkReadMissing-4            500000              7672 ns/op
BenchmarkReadReverse-4           1000000              1946 ns/op
BenchmarkReadSequential-4        1000000              1790 ns/op
BenchmarkWriteRandom-4            200000              6220 ns/op
BenchmarkWriteSequential-4        200000              5479 ns/op
BenchmarkDeleteRandom-4           300000              5689 ns/op
BenchmarkDeleteSequential-4       300000              4974 ns/op
ok      github.com/kezhuw/go-leveldb-benchmarks 158.519s
```

## TODO
- [ ] Source code documentation [WIP]
- [ ] Tests [WIP]
- [ ] Logging
- [ ] Abstract cache interface, so we can share cache among multiple LevelDB instances
- [ ] Reference counting openning file collection, don't rely on GC
- [ ] Statistics
- [x] Benchmarks, See [kezhuw/go-leveldb-benchmarks][go-leveldb-benchmarks].
- [ ] Concurrent level compaction
- [ ] Replace hardcoded constants with configurable options
- [ ] Automatic adjustment of volatile options


## License
The MIT License (MIT). See [LICENSE](LICENSE) for the full license text.


## Contribution
Before going on, make sure you agree on [The MIT License (MIT).](LICENSE).

This project is far from complete, lots of things are missing. If you find any bugs or complement any parts of missing,
throw me a pull request.


## Acknowledgments
* [google/leveldb](https://github.com/google/leveldb) Official implementation written in in C++.  
  My knownledge of LevelDB storage format and implementation details mainly comes from this implementation.

* [syndtr/goleveldb](https://github.com/syndtr/goleveldb) Complete and stable Go implementation.  
  The `DB.Range`, `DB.Prefix`, and `filter.Generator` interface origin from this implementation.

* [golang/leveldb](https://github.com/golang/leveldb) Official but incomplete Go implementation.  
  The `Batch` and `Comparator.AppendSuccessor` origins from this implementation.

[go-leveldb-benchmarks]: https://github.com/kezhuw/go-leveldb-benchmarks

