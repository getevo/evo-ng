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
//  @param v
//  @return uint32
func FNV32(v interface{}) uint32 {
	algorithm := fnv.New32()
	return uint32Hasher(algorithm, generic.Parse(v).String())
}

// FNV32a hashes using fnv32a algorithm
//  @param v
//  @return uint32
func FNV32a(v interface{}) uint32 {
	algorithm := fnv.New32a()
	return uint32Hasher(algorithm, generic.Parse(v).String())
}

// FNV64 hashes using fnv64 algorithm
//  @param v
//  @return uint64
func FNV64(v interface{}) uint64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, generic.Parse(v).String())
}

// FNV64a hashes using fnv64a algorithm
//  @param v
//  @return uint64
func FNV64a(v interface{}) uint64 {
	algorithm := fnv.New64a()
	return uint64Hasher(algorithm, generic.Parse(v).String())
}

// MD5 hashes using md5 algorithm
//  @param v
//  @return string
func MD5(v interface{}) string {
	algorithm := md5.New()
	return stringHasher(algorithm, generic.Parse(v).String())
}

// SHA1 hashes using sha1 algorithm
//  @param v
//  @return string
func SHA1(v interface{}) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, generic.Parse(v).String())
}

// SHA256 hashes using sha256 algorithm
//  @param v
//  @return string
func SHA256(v interface{}) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, generic.Parse(v).String())
}

// SHA512 hashes using sha512 algorithm
//  @param v
//  @return string
func SHA512(v interface{}) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, generic.Parse(v).String())
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

// UUID generates unique identifier
//  @return string
func UUID() string {
	return generator.GenerateIdentifier()
}
