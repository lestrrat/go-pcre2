# go-pcre2
(Work In Progress) PCRE2 binding for Go

## Benchmarks

```
shoebill% go test -v -run=none -benchmem -benchtime=5s -bench .
PASS
BenchmarkGoRegexpMatch-4                      300000         27505 ns/op       46400 B/op         95 allocs/op
BenchmarkPCRE2RegexpMatch-4                   500000         14602 ns/op        2632 B/op         70 allocs/op
BenchmarkGoRegexpMatchString-4                300000         29537 ns/op       46304 B/op         89 allocs/op
BenchmarkPCRE2RegexpMatchString-4             500000         13864 ns/op        2536 B/op         64 allocs/op
BenchmarkGoFindAllIndex-4                     300000         28054 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllIndex-4                  500000         17762 ns/op        3112 B/op         64 allocs/op
BenchmarkGoFindAllStringIndex-4               300000         27144 ns/op       43712 B/op         49 allocs/op
BenchmarkPCRE2FindAllStringIndex-4            500000         18827 ns/op        3016 B/op         61 allocs/op
BenchmarkGoFindSubmatchIndex-4                300000         21354 ns/op       42801 B/op         43 allocs/op
BenchmarkPCRE2FindSubmatchIndex-4             500000         12326 ns/op        2776 B/op         55 allocs/op
BenchmarkGoFindStringSubmatchIndex-4          300000         21164 ns/op       42705 B/op         40 allocs/op
BenchmarkPCRE2FindStringSubmatchIndex-4      1000000         11843 ns/op        2680 B/op         52 allocs/op
BenchmarkGoFindAllSubmatchIndex-4             300000         28490 ns/op       43808 B/op         52 allocs/op
BenchmarkPCRE2FindAllSubmatchIndex-4          300000         20977 ns/op        3688 B/op         73 allocs/op
BenchmarkGoFindAllStringSubmatchIndex-4       300000         27409 ns/op       43712 B/op         49 allocs/op
BenchmarkPCRE2FindAllStringSubmatchIndex-4    500000         18764 ns/op        3592 B/op         70 allocs/op
ok      github.com/lestrrat/go-pcre2    137.179s
```
