# Benchmark Comparison
Generated: 2026-04-26 02:17 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4            71.10µ ±  1%
Service_List_Empty-4        40.47µ ±  3%
Service_List_1000Items-4    2.496m ±  2%
Service_Get-4               42.98µ ± 22%
Service_Update-4            109.9µ ± 16%
Service_Delete-4            68.00µ ±  8%
geomean                     115.0µ

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
Service_Create-4                   58.98µ ± 1%
Service_List_Empty-4               28.80µ ± 1%
Service_List_1000Items-4           2.324m ± 1%
Service_Get-4                      30.39µ ± 3%
Service_Update-4                   79.11µ ± 1%
Service_Delete-4                   49.86µ ± 2%
geomean                            88.28µ

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
