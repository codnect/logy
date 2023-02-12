module github.com/procyon-projects/logy/benchmarks

go 1.17

replace github.com/procyon-projects/logy => ../

require (
	github.com/apex/log v1.9.0
	github.com/go-kit/log v0.2.0
	github.com/procyon-projects/logy v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.28.0
	github.com/sirupsen/logrus v1.9.0
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20221230185412-738e83a70c30
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20200109203555-b30bc20e4fd1
)

require (
	github.com/benbjohnson/clock v1.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gookit/color v1.5.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
