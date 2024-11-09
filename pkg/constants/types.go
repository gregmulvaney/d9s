package constants

type KeySet struct {
	Key         string
	Description string
}

var Commands = []string{
	"containers",
	"networks",
	"volumes",
	"secrets",
}
