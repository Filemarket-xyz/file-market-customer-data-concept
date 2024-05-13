package hash

//go:generate  mockgen -source hash.go -destination ../../mocks/pkg/hash/hash.go
type HashManager interface {
	HashSha256(s string) string
	TriplePassHash(nakedPass string) string
}

type hashManager struct {
}

func NewHashManager() HashManager {
	return &hashManager{}
}
