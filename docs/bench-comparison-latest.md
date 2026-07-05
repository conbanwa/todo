# Benchmark Comparison
Generated: 2026-07-05 02:40 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt            │
                         │       sec/op        │   sec/op     vs base                │
Service_Create-4                   69.06µ ± 1%   73.38µ ± 2%   +6.26% (p=0.000 n=10)
Service_List_Empty-4               37.23µ ± 4%   42.05µ ± 3%  +12.94% (p=0.000 n=10)
Service_List_1000Items-4           2.424m ± 1%   2.531m ± 1%   +4.41% (p=0.000 n=10)
Service_Get-4                      39.85µ ± 3%   43.11µ ± 3%   +8.20% (p=0.000 n=10)
Service_Update-4                   104.1µ ± 1%   110.3µ ± 1%   +6.01% (p=0.000 n=10)
Service_Delete-4                   61.49µ ± 3%   66.17µ ± 2%   +7.63% (p=0.001 n=10)
geomean                            108.0µ        116.2µ        +7.54%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   955.9Ki ± 0%       ~ (p=0.469 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10)
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       +0.00%
¹ all samples are equal

                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │      allocs/op      │  allocs/op   vs base                 │
Service_Create-4                    27.00 ± 0%    27.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                20.00 ± 0%    20.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4           19.78k ± 0%   19.78k ± 0%       ~ (p=1.000 n=10) ¹
Service_Get-4                       39.00 ± 0%    39.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                    67.00 ± 0%    67.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    7.000 ± 0%    7.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                             76.18         76.18       +0.00%
¹ all samples are equal
