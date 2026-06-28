# Benchmark Comparison
Generated: 2026-06-28 03:00 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4             69.06µ ± 1%
Service_List_Empty-4         37.23µ ± 4%
Service_List_1000Items-4     2.424m ± 1%
Service_Get-4                39.85µ ± 3%
Service_Update-4             104.1µ ± 1%
Service_Delete-4             61.49µ ± 3%
geomean                      108.0µ

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

cpu: AMD EPYC 9V74 80-Core Processor                
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                   59.44µ ± 1%
Service_List_Empty-4               29.46µ ± 2%
Service_List_1000Items-4           2.357m ± 1%
Service_Get-4                      30.08µ ± 2%
Service_Update-4                   80.64µ ± 1%
Service_Delete-4                   49.85µ ± 2%
geomean                            89.06µ

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
