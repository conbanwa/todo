# Benchmark Comparison
Generated: 2026-05-17 02:43 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4             69.44µ ± 1%
Service_List_Empty-4         38.36µ ± 2%
Service_List_1000Items-4     2.478m ± 0%
Service_Get-4                40.57µ ± 2%
Service_Update-4             105.4µ ± 1%
Service_Delete-4             61.88µ ± 1%
geomean                      109.7µ

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

cpu: AMD EPYC 9V74 80-Core Processor                
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                   59.12µ ± 1%
Service_List_Empty-4               29.15µ ± 2%
Service_List_1000Items-4           2.323m ± 1%
Service_Get-4                      29.97µ ± 2%
Service_Update-4                   79.34µ ± 1%
Service_Delete-4                   49.78µ ± 1%
geomean                            88.29µ

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
