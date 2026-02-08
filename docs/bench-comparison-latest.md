# Benchmark Comparison
Generated: 2026-02-08 02:12 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 7763 64-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt            │
                         │       sec/op        │    sec/op     vs base               │
Service_Create-4                   70.97µ ± 2%   71.28µ ±  1%       ~ (p=0.436 n=10)
Service_List_Empty-4               39.50µ ± 3%   40.43µ ±  2%  +2.35% (p=0.011 n=10)
Service_List_1000Items-4           2.415m ± 1%   2.412m ±  2%       ~ (p=0.853 n=10)
Service_Get-4                      42.48µ ± 1%   43.02µ ±  2%       ~ (p=0.143 n=10)
Service_Update-4                   110.5µ ± 1%   110.5µ ± 18%       ~ (p=0.853 n=10)
Service_Delete-4                   59.99µ ± 1%   60.36µ ± 15%       ~ (p=0.218 n=10)
geomean                            111.4µ        112.2µ        +0.75%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   956.0Ki ± 0%       ~ (p=0.725 n=10)
Service_Get-4                     1.102Ki ± 0%   1.102Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Update-4                  1.797Ki ± 0%   1.797Ki ± 0%       ~ (p=1.000 n=10) ¹
Service_Delete-4                    144.0 ± 0%     144.0 ± 0%       ~ (p=1.000 n=10) ¹
geomean                           2.234Ki        2.234Ki       +0.00%
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
