# Benchmark Comparison
Generated: 2026-05-31 03:00 UTC

## Changes vs previous run
goos: linux
goarch: amd64
pkg: github.com/conbanwa/todo/internal/transport
cpu: AMD EPYC 9V74 80-Core Processor                
                         │ docs/bench-last.txt │            bench-raw.txt             │
                         │       sec/op        │    sec/op     vs base                │
Service_Create-4                   58.68µ ± 1%   58.49µ ±  1%        ~ (p=0.247 n=10)
Service_List_Empty-4               28.44µ ± 2%   33.23µ ± 14%  +16.84% (p=0.005 n=10)
Service_List_1000Items-4           2.296m ± 1%   2.296m ±  1%        ~ (p=0.481 n=10)
Service_Get-4                      29.84µ ± 1%   29.55µ ±  1%        ~ (p=0.105 n=10)
Service_Update-4                   78.61µ ± 2%   79.77µ ±  1%        ~ (p=0.075 n=10)
Service_Delete-4                   49.68µ ± 3%   50.17µ ±  2%        ~ (p=0.579 n=10)
geomean                            87.43µ        89.90µ         +2.83%

                         │ docs/bench-last.txt │             bench-raw.txt             │
                         │        B/op         │     B/op      vs base                 │
Service_Create-4                    680.0 ± 0%     680.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_Empty-4                720.0 ± 0%     720.0 ± 0%       ~ (p=1.000 n=10) ¹
Service_List_1000Items-4          955.9Ki ± 0%   955.9Ki ± 0%       ~ (p=0.839 n=10)
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
