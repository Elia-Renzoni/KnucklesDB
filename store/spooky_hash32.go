/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/



/**
*	This file contains the implementation of the Spooky Hash algorithm
*	This algorithm is useful for creating the hash function.
*	Returns a 32bit output despite the algorithm produces a 128-bit hash
*	The Algotithm implemented is the short version of the SpookyHash
*   SpookyHashShort is ideal for short messages such IP addresses and API endpoints
 */

package store

import "fmt"

type SpookyHash struct {
	hashSeed uint32
}

var sc_const uint64 = uint64(0xdeadbeefdeadbeef)

func NewSpookyHash(seed uint32) *SpookyHash {
	return &SpookyHash{
		hashSeed: seed,
	}
}

func Rot64(x, k uint64) uint64 {
	return (x << k) | (x >> (64 - k))
}

func ShortMix(h0, h1, h2, h3 uint64) (uint64, uint64, uint64, uint64) {
	h2 = Rot64(h2, 50)
	h2 += h3
	h0 ^= h2
	h3 = Rot64(h3, 52)
	h3 += h0
	h1 ^= h3
	h0 = Rot64(h0, 30)
	h0 += h1
	h2 ^= h0
	h1 = Rot64(h1, 41)
	h1 += h2
	h3 ^= h1
	h2 = Rot64(h2, 54)
	h2 += h3
	h0 ^= h2
	h3 = Rot64(h3, 48)
	h3 += h0
	h1 ^= h3
	h0 = Rot64(h0, 38)
	h0 += h1
	h2 ^= h0
	h1 = Rot64(h1, 37)
	h1 += h2
	h3 ^= h1
	h2 = Rot64(h2, 62)
	h2 += h3
	h0 ^= h2
	h3 = Rot64(h3, 34)
	h3 += h0
	h1 ^= h3
	h0 = Rot64(h0, 5)
	h0 += h1
	h2 ^= h0
	h1 = Rot64(h1, 36)
	h1 += h2
	h3 ^= h1
	return h0, h1, h2, h3
}

func ShortEnd(h0, h1, h2, h3 uint64) (uint64, uint64, uint64, uint64) {
	h3 ^= h2
	h2 = Rot64(h2, 15)
	h3 += h2
	h0 ^= h3
	h3 = Rot64(h3, 52)
	h0 += h3
	h1 ^= h0
	h0 = Rot64(h0, 26)
	h1 += h0
	h2 ^= h1
	h1 = Rot64(h1, 51)
	h2 += h1
	h3 ^= h2
	h2 = Rot64(h2, 28)
	h3 += h2
	h0 ^= h3
	h3 = Rot64(h3, 9)
	h0 += h3
	h1 ^= h0
	h0 = Rot64(h0, 47)
	h1 += h0
	h2 ^= h1
	h1 = Rot64(h1, 54)
	h2 += h1
	h3 ^= h2
	h2 = Rot64(h2, 32)
	h3 += h2
	h0 ^= h3
	h3 = Rot64(h3, 25)
	h0 += h3
	h1 ^= h0
	h0 = Rot64(h0, 63)
	h1 += h0
	return h0, h1, h2, h3
}

func U8tou32le(p []byte) uint64 {
	return uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24
}

func U8tou64le(p []byte) uint64 {
	return uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 | uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56
}

func SpookyHashShort(in []byte, hash1, hash2 uint64) (uint64, uint64) {
	a, b := hash1, hash2
	c, d := uint64(sc_const), uint64(sc_const)
	length := len(in)

	remainder := length % 32
	if length >= 16 {
		for l := length; l >= 32; l -= 32 {
			c += U8tou64le(in)
			in = in[8:]
			d += U8tou64le(in)
			in = in[8:]
			a, b, c, d = ShortMix(a, b, c, d)
			a += U8tou64le(in)
			in = in[8:]
			b += U8tou64le(in)
			in = in[8:]
		}

		if remainder >= 16 {
			c += U8tou64le(in)
			in = in[8:]
			d += U8tou64le(in)
			in = in[8:]
			a, b, c, d = ShortMix(a, b, c, d)
			remainder -= 16
		}
	}

	d += uint64(length) << 56
	switch remainder {
	case 15:
		d += uint64(in[14]) << 48
		fallthrough
	case 14:
		d += uint64(in[13]) << 40
		fallthrough
	case 13:
		d += uint64(in[12]) << 32
		fallthrough
	case 12:
		d += U8tou32le(in[8:])
		c += U8tou64le(in)
		break
	case 11:
		d += uint64(in[10]) << 16
		fallthrough
	case 10:
		d += uint64(in[9]) << 8
		fallthrough
	case 9:
		d += uint64(in[8])
		fallthrough
	case 8:
		c += U8tou64le(in)
		break
	case 7:
		c += uint64(in[6]) << 48
		fallthrough
	case 6:
		c += uint64(in[5]) << 40
		fallthrough
	case 5:
		c += uint64(in[4]) << 32
		fallthrough
	case 4:
		c += U8tou32le(in)
		break
	case 3:
		c += uint64(in[2]) << 16
		fallthrough
	case 2:
		c += uint64(in[1]) << 8
		fallthrough
	case 1:
		c += uint64(in[0])
		break
	case 0:
		c += sc_const
		d += sc_const
	default:
		fmt.Printf("remainder=%d\n", remainder)
		panic("SpookyHash")
	}

	a, b, c, d = ShortEnd(a, b, c, d)
	return a, b
}

func (s *SpookyHash) Hash32(key []byte) uint32 {
	hash1, hash2 := uint64(s.hashSeed), uint64(s.hashSeed)
	hash, _ := SpookyHashShort(key, hash1, hash2)
	return uint32(hash)
}
