# go-pcre2
(Work In Progress) PCRE2 binding for Go

## Benchmarks

```
shoebill% go test -v -run=none -benchmem -benchtime=5s -bench .
BenchmarkGoRegexpMatch-4                  300000         27113 ns/op       46304 B/op         89 allocs/op
BenchmarkPCRE2RegexpMatch-4               500000         14108 ns/op        2536 B/op         64 allocs/op
BenchmarkGoFindAllIndex-4                 300000         27482 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllIndex-4              500000         18518 ns/op        3112 B/op         64 allocs/op
BenchmarkGoFindAllSubmatchIndex-4         300000         35004 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllSubmatchIndex-4      500000         20108 ns/op        3688 B/op         73 allocs/op
ok      github.com/lestrrat/go-pcre2    54.662s
```
