var shell = require("shelljs")
shell.exec("go build ./cmd/api_service && go build ./cmd/api_service-cli")