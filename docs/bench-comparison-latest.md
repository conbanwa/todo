# Benchmark Comparison
Generated: 2026-07-12 02:16 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │           bench-raw.txt            │
                         │       sec/op        │   sec/op     vs base               │
Service_Create-4                   73.38µ ± 2%   72.15µ ± 1%  -1.68% (p=0.009 n=10)
Service_List_Empty-4               42.05µ ± 3%   39.97µ ± 3%  -4.93% (p=0.000 n=10)
Service_List_1000Items-4           2.531m ± 1%   2.426m ± 0%  -4.15% (p=0.000 n=10)
Service_Get-4                      43.11µ ± 3%   42.34µ ± 1%  -1.79% (p=0.001 n=10)
Service_Update-4                   110.3µ ± 1%   109.0µ ± 1%  -1.19% (p=0.002 n=10)
Service_Delete-4                   66.17µ ± 2%   64.97µ ± 2%       ~ (p=0.075 n=10)
geomean                            116.2µ        113.1µ       -2.60%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   956.0Ki ± 0%  +0.00% (p=0.011 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10)
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       +0.00%
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
