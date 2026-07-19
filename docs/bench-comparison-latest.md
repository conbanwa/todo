# Benchmark Comparison
Generated: 2026-07-19 02:14 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt            │
                         │       sec/op        │    sec/op     vs base               │
Service_Create-4                   72.15µ ± 1%   72.12µ ±  1%       ~ (p=0.684 n=10)
Service_List_Empty-4               39.97µ ± 3%   39.90µ ±  2%       ~ (p=0.436 n=10)
Service_List_1000Items-4           2.426m ± 0%   2.500m ±  1%  +3.06% (p=0.000 n=10)
Service_Get-4                      42.34µ ± 1%   43.79µ ± 20%  +3.42% (p=0.000 n=10)
Service_Update-4                   109.0µ ± 1%   109.1µ ±  2%       ~ (p=0.631 n=10)
Service_Delete-4                   64.97µ ± 2%   64.35µ ±  1%  -0.96% (p=0.035 n=10)
geomean                            113.1µ        114.1µ        +0.89%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          956.0Ki ± 0%   955.9Ki ± 0%       ~ (p=0.304 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       -0.00%
¹ all samples are equal

                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │      allocs/op      │  allocs/op   vs base                 │
Service_Create-4                    27.00 ± 0%    27.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                20.00 ± 0%    20.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4           19.78k ± 0%   19.78k ± 0%       ~ (p=0.474 n=10)
Service_Get-4                       39.00 ± 0%    39.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                    67.00 ± 0%    67.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    7.000 ± 0%    7.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                             76.18         76.18       +0.00%
¹ all samples are equal
