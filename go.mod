module github.com/xyths/sss

go 1.13

require (
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sero-cash/go-sero v1.0.3-rc2
	github.com/stretchr/testify v1.4.0
	github.com/xdg/stringprep v1.0.0 // indirect
	github.com/xyths/hs v0.9.3
	github.com/xyths/sero-go v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/text v0.3.2
	gopkg.in/urfave/cli.v2 v2.0.0-20190806201727-b62605953717
)

replace github.com/xyths/sero-go => ../sero-go
