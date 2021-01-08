module github.com/xyths/sss

go 1.13

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sero-cash/go-sero v1.0.3-rc2
	github.com/shopspring/decimal v1.2.0
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli/v2 v2.3.0
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xyths/hs v0.27.2
	github.com/xyths/sero-go v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/zap v1.16.0
	golang.org/x/text v0.3.3
)

replace github.com/xyths/sero-go => ../sero-go
