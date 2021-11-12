package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/getevo/evo-ng/internal/generic"
	"github.com/hart87/GoFlake/generator"
	"hash"
	"hash/fnv"
)



// FNV32 hashes using fnv32 algorithm
func FNV32(v interface{}) uint32 {
	algorithm := fnv.New32()
	return uint32Hasher(algorithm, generic.String(v))
}

// FNV32a hashes using fnv32a algorithm
func FNV32a(v interface{}) uint32 {
	algorithm := fnv.New32a()
	return uint32Hasher(algorithm, generic.String(v))
}

// FNV64 hashes using fnv64 algorithm
func FNV64(v interface{}) uint64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, generic.String(v))
}

// FNV64a hashes using fnv64a algorithm
func FNV64a(v interface{}) uint64 {
	algorithm := fnv.New64a()
	return uint64Hasher(algorithm, generic.String(v))
}

// MD5 hashes using md5 algorithm
func MD5(v interface{}) string {
	algorithm := md5.New()
	return stringHasher(algorithm, generic.String(v))
}

// SHA1 hashes using sha1 algorithm
func SHA1(v interface{}) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, generic.String(v))
}

// SHA256 hashes using sha256 algorithm
func SHA256(v interface{}) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, generic.String(v))
}

// SHA512 hashes using sha512 algorithm
func SHA512(v interface{}) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, generic.String(v))
}



func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func uint32Hasher(algorithm hash.Hash32, text string) uint32 {
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	algorithm.Write([]byte(text))
	return algorithm.Sum64()
}

func UUID() string{
	return generator.GenerateIdentifier()
}