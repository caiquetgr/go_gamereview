module github.com/caiquetgr/go_gamereview

go 1.20

require (
	github.com/gin-gonic/gin v1.8.2
	github.com/uptrace/bun/dialect/pgdialect v1.1.10
	github.com/uptrace/bun/driver/pgdriver v1.1.10
)

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.0.2
	github.com/google/go-cmp v0.5.9
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/google/uuid v1.3.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/uptrace/bun v1.1.10
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	mellium.im/sasl v0.3.1 // indirect
)

replace (
	github.com/cucumber/godog => github.com/laurazard/godog v0.0.0-20220922095256-4c4b17abdae7

	// For k8s dependencies, we use a replace directive, to prevent them being
	// upgraded to the version specified in containerd, which is not relevant to the
	// version needed.
	// See https://github.com/docker/buildx/pull/948 for details.
	// https://github.com/docker/buildx/blob/v0.8.1/go.mod#L62-L64
	k8s.io/api => k8s.io/api v0.22.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.4
	k8s.io/client-go => k8s.io/client-go v0.22.4
)
