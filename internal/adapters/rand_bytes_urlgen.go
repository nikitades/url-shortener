package adapters

import "math/rand"

type RandBytesUrlGenerator struct{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func (r *RandBytesUrlGenerator) Generate(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func NewRandBytesUrlGenerator() *RandBytesUrlGenerator {
	return &RandBytesUrlGenerator{}
}
