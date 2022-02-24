module main/Breadbot

go 1.17

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/tidwall/gjson v1.14.0
	main/gleam v0.0.0-00010101000000-000000000000
	main/testFolder v0.0.0-00010101000000-000000000000
)

require (
	github.com/brianvoe/gofakeit/v6 v6.14.5 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/net v0.0.0-20211116231205-47ca1ff31462 // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
)

replace main/gleam => ../gleam

replace main/testFolder => ../testFolder
