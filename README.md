# go-pcre2
(Work In Progress) PCRE2 binding for Go

## Benchmarks

```
shoebill% go test -v -run=none -benchmem -benchtime=5s -bench .
PASS
BenchmarkGoRegexpMatch-4                      300000         27679 ns/op       46304 B/op         89 allocs/op
BenchmarkPCRE2RegexpMatch-4                   500000         14706 ns/op        2536 B/op         64 allocs/op
BenchmarkGoFindAllIndex-4                     200000         34696 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllIndex-4                  300000         20060 ns/op        3112 B/op         64 allocs/op
BenchmarkGoFindAllSubmatchIndex-4             200000         32344 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllSubmatchIndex-4          300000         20491 ns/op        3688 B/op         73 allocs/op
BenchmarkGoFindAllStringSubmatchIndex-4       200000         31611 ns/op       43712 B/op         49 allocs/op
BenchmarkPCRE2FindAllStringSubmatchIndex-4    300000         20700 ns/op        3592 B/op         70 allocs/op
ok      github.com/lestrrat/go-pcre2    55.948s
```
