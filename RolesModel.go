package main

var Roles = map[string]byte{
	"user":  1 << 0,
	"mod":   1 << 1,
	"dev":   1 << 2,
	"owner": 1 << 3,
}
