# Benchmark Comparison
Generated: 2026-02-15 01:57 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │           bench-raw.txt            │
                         │       sec/op        │   sec/op     vs base               │
Service_Create-4                  71.28µ ±  1%   72.64µ ± 1%  +1.91% (p=0.000 n=10)
Service_List_Empty-4              40.43µ ±  2%   40.42µ ± 2%       ~ (p=0.393 n=10)
Service_List_1000Items-4          2.412m ±  2%   2.386m ± 0%  -1.05% (p=0.000 n=10)
Service_Get-4                     43.02µ ±  2%   42.19µ ± 1%  -1.93% (p=0.009 n=10)
Service_Update-4                  110.5µ ± 18%   108.1µ ± 1%  -2.15% (p=0.000 n=10)
Service_Delete-4                  60.36µ ± 15%   63.05µ ± 2%       ~ (p=0.481 n=10)
geomean                           112.2µ         112.4µ       +0.17%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          956.0Ki ± 0%   955.9Ki ± 0%  -0.00% (p=0.000 n=10)
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
