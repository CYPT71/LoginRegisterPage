package domain

var Permissions = map[string]uint64{
	"owner": 1 << 0,
	"mod":   1 << 1,
	"dev":   1 << 2,
	"user":  1 << 3,
}
