# Benchmark Comparison
Generated: 2026-09-06 02:32 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt            │
                         │       sec/op        │    sec/op     vs base               │
Service_Create-4                  69.38µ ±  0%   70.06µ ±  1%  +0.98% (p=0.001 n=10)
Service_List_Empty-4              38.25µ ±  3%   39.30µ ±  3%  +2.73% (p=0.009 n=10)
Service_List_1000Items-4          2.446m ±  1%   2.458m ±  0%  +0.51% (p=0.023 n=10)
Service_Get-4                     41.99µ ±  1%   42.77µ ±  1%  +1.84% (p=0.019 n=10)
Service_Update-4                  108.0µ ± 14%   110.1µ ± 15%       ~ (p=0.218 n=10)
Service_Delete-4                  59.94µ ± 14%   64.91µ ±  2%       ~ (p=0.089 n=10)
geomean                           109.9µ         112.9µ        +2.68%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          956.0Ki ± 0%   955.9Ki ± 0%       ~ (p=0.148 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       -0.00%
¹ all samples are equal

                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │      allocs/op      │  allocs/op   vs base                 │
Service_Create-4                    27.00 ± 0%    27.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                20.00 ± 0%    20.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4           19.78k ± 0%   19.78k ± 0%       ~ (p=0.087 n=10)
Service_Get-4                       39.00 ± 0%    39.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                    67.00 ± 0%    67.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    7.000 ± 0%    7.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                             76.18         76.18       +0.00%
¹ all samples are equal
