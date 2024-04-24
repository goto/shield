module github.com/goto/shield

go 1.21

require (
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/abbot/go-http-auth v0.4.0
	github.com/authzed/authzed-go v0.7.1-0.20221109204547-1aa903788b3b
	github.com/authzed/grpcutil v0.0.0-20230908193239-4286bb1d6403
	github.com/authzed/spicedb v1.15.0
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/envoyproxy/protoc-gen-validate v1.0.2
	github.com/ghodss/yaml v1.0.0
	github.com/golang-migrate/migrate/v4 v4.16.0
	github.com/golang/protobuf v1.5.3
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/goto/salt v0.3.6-0.20240424230902-5eb203bb515d
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0
	github.com/jackc/pgconn v1.14.0
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa
	github.com/jackc/pgx/v4 v4.18.1
	github.com/jhump/protoreflect v1.14.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.7
	github.com/mcuadros/go-defaults v1.2.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/newrelic/go-agent v3.20.2+incompatible
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.7.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/xid v1.4.0
	github.com/spf13/afero v1.9.3
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/gjson v1.14.4
	go.opencensus.io v0.24.0
	go.uber.org/zap v1.24.0
	gocloud.dev v0.28.0
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0
	golang.org/x/net v0.17.0
	golang.org/x/oauth2 v0.13.0
	google.golang.org/genproto/googleapis/api v0.0.0-20231002182017-d307bd883b97
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/oklog/run v1.1.0 // indirect

require (
	cloud.google.com/go v0.110.8 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.2 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/aymanbagabas/go-osc52 v1.2.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/briandowns/spinner v1.20.0 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/charmbracelet/glamour v0.6.0 // indirect
	github.com/cli/safeexec v1.0.1 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/frankban/quicktest v1.14.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/cel-go v0.13.0 // indirect
	github.com/google/pprof v0.0.0-20221219190121-3cb0bae90811 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.4 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jeremywohl/flatten v1.0.1 // indirect
	github.com/jzelinskie/stringz v0.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.21 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.13.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.1.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rs/zerolog v1.29.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/schollz/progressbar/v3 v3.13.0 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/yuin/goldmark v1.5.3 // indirect
	github.com/yuin/goldmark-emoji v1.0.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.128.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
