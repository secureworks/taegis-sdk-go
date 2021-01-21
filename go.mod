module github.com/secureworks/tdr-sdk-go

go 1.13

require (
	github.com/99designs/gqlgen v0.11.3
	github.com/VerticalOps/fakesentry v0.0.0-20200528170726-b873dcced65c
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/certifi/gocertifi v0.0.0-20200211180108-c7c1fbc02894 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/getsentry/sentry-go v0.7.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/nulls v0.4.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/makasim/sentryhook v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.14.0
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/rs/zerolog v1.19.0
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.5.1
	github.com/vektah/gqlparser/v2 v2.1.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	moul.io/http2curl v1.0.0
)

replace golang.org/x/text => golang.org/x/text v0.3.3
