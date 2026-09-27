# Benchmark Comparison
Generated: 2026-09-27 03:12 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4             71.25µ ± 2%
Service_List_Empty-4         40.04µ ± 2%
Service_List_1000Items-4     2.460m ± 0%
Service_Get-4                42.78µ ± 1%
Service_Update-4             108.2µ ± 1%
Service_Delete-4             64.05µ ± 2%
geomean                      113.0µ

                         │ bench-raw.txt │
                         │     B/op      │
Service_Create-4              680.0 ± 0%
Service_List_Empty-4          720.0 ± 0%
Service_List_1000Items-4    956.0Ki ± 0%
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

cpu: INTEL(R) XEON(R) PLATINUM 8573C
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                  37.17µ ±  1%
Service_List_Empty-4              21.05µ ± 13%
Service_List_1000Items-4          2.862m ±  2%
Service_Get-4                     23.87µ ±  9%
Service_Update-4                  53.05µ ±  1%
Service_Delete-4                  30.60µ ±  1%
geomean                           66.53µ

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
