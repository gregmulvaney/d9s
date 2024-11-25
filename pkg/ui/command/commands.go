package command

// No way to make an immutable array.
var Commands = [...]string{
	"containers",
	"networks",
	"secrets",
	"volumes",
	"q",
}
