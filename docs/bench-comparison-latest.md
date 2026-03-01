# Benchmark Comparison
Generated: 2026-03-01 01:58 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │           bench-raw.txt            │
                         │       sec/op        │   sec/op     vs base               │
Service_Create-4                   72.05µ ± 2%   74.15µ ± 1%  +2.92% (p=0.003 n=10)
Service_List_Empty-4               40.37µ ± 1%   40.62µ ± 1%       ~ (p=0.494 n=10)
Service_List_1000Items-4           2.422m ± 0%   2.407m ± 1%  -0.62% (p=0.019 n=10)
Service_Get-4                      41.89µ ± 2%   42.35µ ± 1%  +1.11% (p=0.043 n=10)
Service_Update-4                   107.6µ ± 1%   109.0µ ± 2%  +1.34% (p=0.003 n=10)
Service_Delete-4                   62.14µ ± 1%   63.43µ ± 3%  +2.07% (p=0.000 n=10)
geomean                            112.0µ        113.4µ       +1.24%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   955.9Ki ± 0%  -0.00% (p=0.000 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       -0.00%
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
