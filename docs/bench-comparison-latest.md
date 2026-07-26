# Benchmark Comparison
Generated: 2026-07-26 02:18 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │
                         │       sec/op        │
Service_Create-4                  72.12µ ±  1%
Service_List_Empty-4              39.90µ ±  2%
Service_List_1000Items-4          2.500m ±  1%
Service_Get-4                     43.79µ ± 20%
Service_Update-4                  109.1µ ±  2%
Service_Delete-4                  64.35µ ±  1%
geomean                           114.1µ

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

cpu: AMD EPYC 9V74 80-Core Processor                
                         │ bench-raw.txt │
                         │    sec/op     │
Service_Create-4             60.89µ ± 1%
Service_List_Empty-4         30.82µ ± 2%
Service_List_1000Items-4     2.346m ± 0%
Service_Get-4                31.43µ ± 1%
Service_Update-4             83.20µ ± 2%
Service_Delete-4             50.95µ ± 2%
geomean                      91.49µ

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
