module github.com/nonetype/gocdb

go 1.17

replace github.com/nonetype/gocdb/subprocess => ./pkg/subprocess

replace github.com/nonetype/gocdb/cdbController => ./pkg/cdbController

require github.com/nonetype/gocdb/subprocess v0.0.0-00010101000000-000000000000

require github.com/nonetype/gocdb/cdbController v0.0.0-00010101000000-000000000000
