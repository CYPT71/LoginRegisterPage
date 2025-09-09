module webauthn_api

go 1.24.0

toolchain go1.24.3

require (
        github.com/go-webauthn/webauthn v0.13.4
        github.com/golang-jwt/jwt v3.2.2+incompatible
        github.com/google/uuid v1.6.0
        github.com/joho/godotenv v1.5.1
    github.com/redis/go-redis/v9 v9.7.0
    github.com/swaggo/fiber-swagger v1.3.0
        github.com/swaggo/swag v1.16.6
        gorm.io/driver/postgres v1.6.0
)

replace github.com/redis/go-redis/v9 => ./internal/redismock

require (
       github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-openapi/jsonpointer v0.22.0 // indirect
	github.com/go-openapi/jsonreference v0.21.1 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.24.1 // indirect
	github.com/go-openapi/swag/cmdutils v0.24.0 // indirect
	github.com/go-openapi/swag/conv v0.24.0 // indirect
	github.com/go-openapi/swag/fileutils v0.24.0 // indirect
	github.com/go-openapi/swag/jsonname v0.24.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.24.0 // indirect
	github.com/go-openapi/swag/loading v0.24.0 // indirect
	github.com/go-openapi/swag/mangling v0.24.0 // indirect
	github.com/go-openapi/swag/netutils v0.24.0 // indirect
	github.com/go-openapi/swag/stringutils v0.24.0 // indirect
	github.com/go-openapi/swag/typeutils v0.24.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.24.0 // indirect
	github.com/go-webauthn/x v0.1.24 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	github.com/tinylib/msgp v1.4.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/gofiber/fiber/v2 v2.52.9
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.65.0
	golang.org/x/crypto v0.41.0
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/grpc v1.75.0
	gorm.io/gorm v1.30.3
)
