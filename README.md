![Logy Logo](https://user-images.githubusercontent.com/5354910/223832463-ed73a30a-7edf-4d09-886e-945ea9157b1c.png)

# Logy

## Performance

Log a message without context fields

| Package                 |    Time       | Objects Allocated |
| :-----------------------|  :----------: | :---------------: |
| :star: logy             |   27.99 ns/op |    0 allocs/op    |
| :star: logy(formatting) |  883.8 ns/op  |    7 allocs/op    |
| :zap: exp/slog          |   38.08 ns/op |    0 allocs/op    |
| zerolog                 |   37.49 ns/op |    0 allocs/op    |
| zerolog(formatting)     |   3030 ns/op  |    108 allocs/op  |
| zap                     |   98.30 ns/op |    0 allocs/op    |
| zap sugar               |   110.9 ns/op |    1 allocs/op    |
| zap sugar (formatting)  |   3369 ns/op  |    108 allocs/op  |
| go-kit                  |  248.5 ns/op  |    9 allocs/op    |
| log15                   |  2490 ns/op   |    20 allocs/op   |
| apex/log                |  1139 ns/op   |    6 allocs/op    |
| logrus                  |  1831 ns/op   |   23 allocs/op    |

Log a message with a logger that already has 10 fields of context:

| Package                 |    Time       | Objects Allocated |
| :-----------------------|  :----------: | :---------------: |
| :star: logy             |   61.43 ns/op |    0 allocs/op    |
| :star: logy(formatting) | 1208.0 ns/op  |    7 allocs/op    |
| :zap: exp/slog          |   266.3 ns/op |    0 allocs/op    |
| zerolog                 |   44.84 ns/op |    0 allocs/op    |
| zerolog(formatting)     | 3103.0 ns/op  |    108 allocs/op  |
| zap                     |   92.50 ns/op |    0 allocs/op    |
| zap sugar               |   113.7 ns/op |    1 allocs/op    |
| zap sugar (formatting)  |   3355 ns/op  |    108 allocs/op  |
| go-kit                  |  3628 ns/op   |   66 allocs/op    |
| log15                   | 12532 ns/op   |   130 allocs/op   |
| apex/log                | 14494 ns/op   |   53 allocs/op    |
| logrus                  | 16246 ns/op   |   68 allocs/op    |
