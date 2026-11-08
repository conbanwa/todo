# Benchmark Comparison
Generated: 2026-11-08 03:58 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                   74.15µ ± 1%
Service_List_Empty-4               41.48µ ± 4%
Service_List_1000Items-4           2.473m ± 1%
Service_Get-4                      44.22µ ± 2%
Service_Update-4                   110.7µ ± 5%
Service_Delete-4                   64.15µ ± 1%
geomean                            115.6µ

                         │ docs/bench-last.txt │
                         │        B/op         │
Service_Create-4                    680.0 ± 0%
Service_List_Empty-4                720.0 ± 0%
Service_List_1000Items-4          955.9Ki ± 0%
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

cpu: AMD EPYC 9V45 96-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4             35.17µ ± 1%
Service_List_Empty-4         16.99µ ± 2%
Service_List_1000Items-4     2.079m ± 0%
Service_Get-4                18.30µ ± 2%
Service_Update-4             42.79µ ± 3%
Service_Delete-4             26.30µ ± 1%
geomean                      54.28µ

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
