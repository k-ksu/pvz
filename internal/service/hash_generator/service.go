package hash_generator

import "HomeWork_1/pkg/hash"

type HashGenerator struct{}

func NewHashGenerator() HashGenerator {
	return HashGenerator{}
}

func (h *HashGenerator) Generate() string {
	return hash.GenerateHash()
}
