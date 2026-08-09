# Benchmark Comparison
Generated: 2026-08-09 01:21 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │           bench-raw.txt            │
                         │       sec/op        │   sec/op     vs base               │
Service_Create-4                  71.62µ ±  1%   70.47µ ± 0%  -1.60% (p=0.000 n=10)
Service_List_Empty-4              40.75µ ±  2%   39.75µ ± 3%  -2.45% (p=0.023 n=10)
Service_List_1000Items-4          2.547m ±  1%   2.479m ± 1%  -2.67% (p=0.000 n=10)
Service_Get-4                     42.39µ ± 24%   42.24µ ± 2%       ~ (p=0.280 n=10)
Service_Update-4                  108.2µ ±  2%   108.7µ ± 1%       ~ (p=0.424 n=10)
Service_Delete-4                  64.23µ ± 13%   63.93µ ± 2%       ~ (p=0.579 n=10)
geomean                           114.0µ         112.6µ       -1.20%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   955.9Ki ± 0%       ~ (p=0.196 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       -0.00%
¹ all samples are equal

                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │      allocs/op      │  allocs/op   vs base                 │
Service_Create-4                    27.00 ± 0%    27.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                20.00 ± 0%    20.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4           19.78k ± 0%   19.78k ± 0%       ~ (p=1.000 n=10)
Service_Get-4                       39.00 ± 0%    39.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                    67.00 ± 0%    67.00 ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    7.000 ± 0%    7.000 ± 0%       ~ (p=1.000 n=10) ¹
geomean                             76.18         76.18       +0.00%
¹ all samples are equal
