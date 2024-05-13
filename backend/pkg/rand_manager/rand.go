package rand_manager

import (
	"fmt"
	"math/rand"
)

//go:generate mockgen -source rand.go -destination ../../mocks/pkg/rand_manager/rand.go

type RandManager interface {
	Int63n(n int64) int64
	Intn(n int) int
	// возвращает числовой код заданной длины
	RandomIntCode(len int) string
}

type randManager struct {
	r *rand.Rand
}

func NewRandManager(seed int64) RandManager {
	return &randManager{
		r: rand.New(rand.NewSource(seed)),
	}
}

func (r *randManager) Intn(n int) int {
	return r.r.Intn(n)
}

func (r *randManager) Int63n(n int64) int64 {
	return r.r.Int63n(n)
}

func (r *randManager) RandomIntCode(len int) string {
	var code string
	for i := 0; i < len; i++ {
		code = fmt.Sprintf("%s%d", code, r.r.Intn(10))
	}
	return code
}
