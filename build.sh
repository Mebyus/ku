set -o errexit  # Exit on most errors
set -o errtrace # Make sure any error trap is inherited
set -o nounset  # Disallow expansion of unset variables
set -o pipefail # Use last non-zero exit code in a pipeline

go fmt ./...
go test ./internal/...

CGO_ENABLED=0 go build -o bin/ku ./cmd/ku
