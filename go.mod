module auth

go 1.18

replace auth_grpc => ./pkg/auth_grpc

require (
	auth_grpc v0.0.0
	github.com/getsentry/sentry-go v0.13.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-stack/stack v1.8.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/rs/zerolog v1.27.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.9.0
	github.com/swaggo/http-swagger v1.3.3
	github.com/swaggo/swag v1.8.1
	go.mongodb.org/mongo-driver v1.10.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.34.0
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/exporters/jaeger v1.9.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
	golang.org/x/crypto v0.0.0-20220824171710-5757bc0c5503
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	google.golang.org/grpc v1.49.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Microsoft/hcsshim v0.9.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/cgroups v1.0.3 // indirect
	github.com/containerd/containerd v1.6.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.13+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-migrate/migrate/v4 v4.15.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.5.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	github.com/testcontainers/testcontainers-go v0.13.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/otel/metric v0.31.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
