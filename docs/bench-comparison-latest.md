# Benchmark Comparison
Generated: 2026-04-12 02:10 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                   68.96µ ± 1%
Service_List_Empty-4               37.84µ ± 3%
Service_List_1000Items-4           2.446m ± 1%
Service_Get-4                      41.46µ ± 3%
Service_Update-4                   104.6µ ± 1%
Service_Delete-4                   62.39µ ± 1%
geomean                            109.5µ

                         │ docs/bench-last.txt │
                         │        B/op         │
Service_Create-4                    680.0 ± 0%
Service_List_Empty-4                720.0 ± 0%
Service_List_1000Items-4          956.0Ki ± 0%
Service_Get-4                     1.102Ki ± 0%
Service_Update-4                  1.797Ki ± 0%
Service_Delete-4                    144.0 ± 0%
geomean                           2.234Ki

                         │ docs/bench-last.txt │
                         │      allocs/op      │
Service_Create-4                    27.00 ± 0%
Service_List_Empty-4                20.00 ± 0%
Service_List_1000Items-4           19.78k ± 0%
Service_Get-4                       39.00 ± 0%
Service_Update-4                    67.00 ± 0%
Service_Delete-4                    7.000 ± 0%
geomean                             76.18

cpu: Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4            49.02µ ±  1%
Service_List_Empty-4        28.68µ ±  2%
Service_List_1000Items-4    3.520m ±  1%
Service_Get-4               36.36µ ± 15%
Service_Update-4            75.56µ ±  3%
Service_Delete-4            41.20µ ±  2%
geomean                     90.80µ

                         │ bench-raw.txt │
                         │     B/op      │
Service_Create-4              680.0 ± 0%
Service_List_Empty-4          720.0 ± 0%
Service_List_1000Items-4    955.9Ki ± 0%
Service_Get-4               1.102Ki ± 0%
Service_Update-4            1.797Ki ± 0%
Service_Delete-4              144.0 ± 0%
geomean                     2.234Ki

                         │ bench-raw.txt │
                         │   allocs/op   │
Service_Create-4              27.00 ± 0%
Service_List_Empty-4          20.00 ± 0%
Service_List_1000Items-4     19.78k ± 0%
Service_Get-4                 39.00 ± 0%
Service_Update-4              67.00 ± 0%
Service_Delete-4              7.000 ± 0%
geomean                       76.18
