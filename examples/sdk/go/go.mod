module github.com/alibaba/UnifiedModel/examples/sdk/go

go 1.22.2

toolchain go1.24.2

require (
	gopkg.in/yaml.v3 v3.0.1
	umodel_go_cli v0.0.0
)

replace umodel_go_cli => ../../../sdk/go
