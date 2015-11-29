# go-pcre2
(Work In Progress) PCRE2 binding for Go

## Benchmarks

```
shoebill% go test -v -run=none -benchmem -benchtime=5s -bench .
PASS
BenchmarkGoRegexpMatch-4          200000         33010 ns/op       46304 B/op         89 allocs/op
BenchmarkPCRE2RegexpMatch-4       500000         14796 ns/op        2536 B/op         64 allocs/op
BenchmarkGoFindAllIndex-4         300000         25526 ns/op       42064 B/op         41 allocs/op
BenchmarkPCRE2FindAllIndex-4      500000         15802 ns/op        2312 B/op         58 allocs/op
shoebill% go test -v -run=none -bench .
```
