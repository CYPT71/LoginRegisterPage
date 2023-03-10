package domain

var Permissions = map[string]uint64{
	"user":  1 << 0,
	"mod":   1 << 1,
	"dev":   1 << 2,
	"owner": 1 << 3,
}
