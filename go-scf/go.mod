module github.com/riba2534/wecomchan/go-scf

go 1.16

require (
	github.com/json-iterator/go v1.1.11
	huaweicloud.com/go-runtime v0.0.0-00010101000000-000000000000
)

replace (
    huaweicloud.com/go-runtime => ./huaweicloud-go-runtime
)
