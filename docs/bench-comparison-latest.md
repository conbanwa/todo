# Benchmark Comparison
Generated: 2026-08-23 01:06 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt            │
                         │       sec/op        │    sec/op     vs base               │
Service_Create-4                   70.98µ ± 2%   69.38µ ±  0%  -2.26% (p=0.000 n=10)
Service_List_Empty-4               40.10µ ± 2%   38.25µ ±  3%  -4.60% (p=0.005 n=10)
Service_List_1000Items-4           2.489m ± 0%   2.446m ±  1%  -1.72% (p=0.000 n=10)
Service_Get-4                      42.20µ ± 1%   41.99µ ±  1%       ~ (p=0.353 n=10)
Service_Update-4                   108.4µ ± 1%   108.0µ ± 14%       ~ (p=0.739 n=10)
Service_Delete-4                   64.31µ ± 2%   59.94µ ± 14%       ~ (p=0.143 n=10)
geomean                            113.0µ        109.9µ        -2.74%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   956.0Ki ± 0%       ~ (p=0.078 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       +0.00%
¹ all samples are equal

                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │      allocs/op      │  allocs/op   vs base                 │
Service_Create-4                    27.00 ± 0%    27.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                20.00 ± 0%    20.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4           19.78k ± 0%   19.78k ± 0%       ~ (p=0.628 n=10)
Service_Get-4                       39.00 ± 0%    39.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                    67.00 ± 0%    67.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    7.000 ± 0%    7.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                             76.18         76.18       +0.00%
¹ all samples are equal
