package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var DEBUG string = os.Getenv("DEBUG")

var K = [64]uint32{
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
	0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
	0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
	0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
	0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
	0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
	0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
	0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
	0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
}

var s = [64]uint32{
	7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
	5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
	4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
	6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
}

func md5(input *os.File) string {
	var (
		a0 uint32 = 0x67452301
		b0 uint32 = 0xefcdab89
		c0 uint32 = 0x98badcfe
		d0 uint32 = 0x10325476
	)

	var (
		byte_block        = make([]byte, 0, 64)
		byte_count uint64 = 0
	)

	processInput(input,
		&a0, &b0, &c0, &d0, &byte_block, &byte_count)
	hash := outputHash(&a0, &b0, &c0, &d0)
	println(hash)
	return hash
}

func DebugMsg(s string) {
	if DEBUG == string('1') {
		fmt.Println("[DEBUG] " + s)
	}
}

func leftrotate(x uint32, n uint32) uint32 {
	return ((x << n) | (x >> (32 - n)))
}

func outputHash(a0, b0, c0, d0 *uint32) string {
	DebugMsg(fmt.Sprintf("a0=%d b0=%d c0=%d d0=%d", a0, b0, c0, d0))

	buf1 := make([]byte, 4)
	buf2 := make([]byte, 4)
	buf3 := make([]byte, 4)
	buf4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf1, *a0)
	binary.LittleEndian.PutUint32(buf2, *b0)
	binary.LittleEndian.PutUint32(buf3, *c0)
	binary.LittleEndian.PutUint32(buf4, *d0)
	return fmt.Sprintf("%x%x%x%x", buf1, buf2, buf3, buf4)
}

func processBlock(block []byte, a0, b0, c0, d0 *uint32, byte_block *[]byte, byte_count *uint64) {
	// DebugMsg(fmt.Sprintf("process block %v %x", block, len(block)))

	M := make([]uint32, 0, 16)

	for i := 0; i < 16; i++ {
		j := i * 4
		M = append(M,
			uint32(block[j])|
				uint32(block[j+1])<<8|
				uint32(block[j+2])<<16|
				uint32(block[j+3])<<24)
	}

	//DebugMsg(fmt.Sprintf("M=%s", M))

	A := *a0
	B := *b0
	C := *c0
	D := *d0

	for i := uint32(0); i < 64; i++ {

		var F, g uint32

		if i <= 15 {
			F = (B & C) | ((^B) & D)
			g = i
		} else if i <= 31 {
			F = (D & B) | ((^D) & C)
			g = (5*i + 1) % 16
		} else if i <= 47 {
			F = B ^ C ^ D
			g = ((3 * i) + 5) % 16
		} else if i <= 63 {
			F = C ^ (B | (^D))
			g = (7 * i) % 16
		}

		F = F + A + K[i] + uint32(M[g])
		A = D
		D = C
		C = B
		B = B + leftrotate(F, s[i])

	}

	*a0 = *a0 + A
	*b0 = *b0 + B
	*c0 = *c0 + C
	*d0 = *d0 + D
}

func processByte(c byte, a0, b0, c0, d0 *uint32, byte_block *[]byte, byte_count *uint64) {
	DebugMsg(fmt.Sprintf("process byte %d", c))

	*byte_block = append(*byte_block, c)
	if len(*byte_block) == 64 {
		processBlock(*byte_block, a0, b0, c0, d0, byte_block, byte_count)
		*byte_block = (*byte_block)[:0] // empties the slice
	}
}

func processInput(input *os.File, a0, b0, c0, d0 *uint32, byte_block *[]byte, byte_count *uint64) {
	reader := bufio.NewReader(input)

	chunkSize := 4096
	buf := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buf)

		if n > 0 {
			DebugMsg(fmt.Sprintf("read %x byte(s)\n", len(buf[:n])))
			for i := range n {
				processByte(buf[i], a0, b0, c0, d0, byte_block, byte_count)
				*byte_count += 1
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

	}
	processByte(0x80, a0, b0, c0, d0, byte_block, byte_count)
	for len(*byte_block) < 56 {
		processByte(0, a0, b0, c0, d0, byte_block, byte_count)
	}

	// converting the length of the message to little endian binary and append it to the byte_block
	lenbuf := [8]byte{}
	bitlen := uint64(*byte_count) * 8
	binary.LittleEndian.PutUint64(lenbuf[:], bitlen)
	for _, bt := range lenbuf {
		processByte(bt, a0, b0, c0, d0, byte_block, byte_count)
	}

	// manual way to convert the length of the message to little endian binary and append it to the byte_block
	// low_bits := uint32(byte_count << 3)
	// high_bits := uint32(byte_count >> (32 - 3))

	// for i := 0; i < 4; i++ {
	// 	bt := byte(low_bits >> (i * 8))
	// 	// DebugMsg(fmt.Sprintf("low_bits[%x]=%v", i, bt))
	// 	processByte(bt)
	// }

	// for i := 0; i < 4; i++ {
	// 	bt := byte(high_bits >> (i * 8))
	// 	// DebugMsg(fmt.Sprintf("high_bits[%x]=%v", i, bt))
	// 	processByte(bt)
	// }

	if len(*byte_block) != 0 {
		DebugMsg("something is wrong")
		os.Exit(1)
	}
}

func main() {
	args := os.Args

	if len(args) > 1 {
		DebugMsg("reading data from file")
		fp, err := os.Open(args[1])
		if err != nil {
			log.Panic(err)
		}
		md5(fp)
	} else {
		DebugMsg("reading data from input")
		md5(os.Stdin)
	}
}
