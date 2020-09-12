package rand

import (
	"bytes"
	"encoding/binary"
	"math"
	mrand "math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brickman-source/golang-utilities/errors"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	//r       = mrand.New(mrand.NewSource(time.Now().UnixNano()))
	r = mrand.New(&lockedSource{src: mrand.NewSource(time.Now().UnixNano()).(mrand.Source64)})

	seqNo   uint64
	Address uint64
)

func init() {
	Address = uint64(getIpAddr())
	//address = getMacAddr()
}

func GetNonceString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func GetShortTimestampSequenceNo() uint64 {
	timeStamp := uint64(time.Now().Unix() - 1466035200)

	sn := atomic.AddUint64(&seqNo, 1)

	a := sn
	i := 0
	for a > 0 {
		a /= 10
		i++
	}

	return  timeStamp * uint64(math.Pow10(i)) + sn
}

func GetShortTimestampRandString() string {
	return GetRandStringWithTimestamp(time.Now().Unix())
}

// 返回：timeStamp(unix epoch),ip address, random number
func ParseShortTimestampRandString(str string) (int64, uint64, int64) {
	parts := strings.Split(str, "0")
	if len(parts) == 3 {
		ts, err := read36N(parts[0])
		if err != nil {
			return 0, 0, 0
		}
		addr, err := read36N(parts[1])
		if err != nil {
			return 0, 0, 0
		}
		sn, err := read36N(parts[2])
		if err != nil {
			return 0, 0, 0
		}

		return int64(ts) + 1466035200, addr, int64(sn)
	}
	return 0, 0, 0
}
func GetRandStringWithTimestamp(epochTime int64) string {
	timeStamp := epochTime - 1466035200
	sn := atomic.AddUint64(&seqNo, 1)

	buf := new(bytes.Buffer)
	write36N(buf, uint64(timeStamp))
	buf.WriteByte('0')
	write36N(buf, Address)
	buf.WriteByte('0')
	write36N(buf, sn)
	return buf.String()
}

func GetNormalTimestampRandString() string {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	sn := atomic.AddUint64(&seqNo, 1)

	buf := new(bytes.Buffer)

	writeN(buf, uint64(year%100), 2)
	writeN(buf, uint64(int(month)), 2)
	writeN(buf, uint64(day), 2)
	writeN(buf, uint64(hour), 2)
	writeN(buf, uint64(min), 2)
	writeN(buf, uint64(sec), 2)
	buf.WriteByte('-')
	writeN(buf, sn, 2)
	buf.WriteByte('-')
	writeN(buf, Address, 2)
	writeN(buf, uint64(r.Intn(0xFF)), 1)

	return buf.String()
}
func GetRandFloat64() float64 {
	mrand.Float64()
	return r.Float64()
}
func GetRandFloat32() float32 {
	return r.Float32()
}

//[min,max)
func GetRandInt64(min int64, max int64) int64 {
	return min + r.Int63n(max-min)
}

//[min,max)
func GetRandInt32(min int32, max int32) int32 {
	return min + r.Int31n(max-min)
}

//[min,max)
func GetRandInt(min int, max int) int {
	return min + r.Intn(max-min)
}
func writeN(buffer *bytes.Buffer, x uint64, width int) {
	u := uint(x)
	if x < 0 {
		buffer.WriteByte('-')
		u = uint(-x)
	}

	var buf [20]byte
	i := len(buf)
	for u >= 10 {
		i--
		q := u / 10
		buf[i] = byte('0' + u - q*10)
		u = q
	}
	i--
	buf[i] = byte('0' + u)

	for w := len(buf) - i; w < width; w++ {
		buffer.WriteByte('0')
	}

	for k := i; k < 20; k++ {
		buffer.WriteByte(buf[k])
	}
}

const digits = "123456789abcdefghijklmnopqrstuvwxyz"
const maxUint64 = (1<<64 - 1)

func write36N(buffer *bytes.Buffer, u uint64) {
	base := 35

	var buf [64 + 1]byte
	i := len(buf)

	b := uint64(base)
	for u >= b {
		i--
		q := u / b
		buf[i] = digits[uint(u-q*b)]
		u = q
	}
	i--
	buf[i] = digits[uint(u)]

	for k := i; k < 65; k++ {
		buffer.WriteByte(buf[k])
	}
}
func read36N(s string) (uint64, error) {
	if len(s) == 0 {
		return 0, errors.New("syntaxError")
	}
	base := 35

	bitSize := 64

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	var cutoff uint64
	switch base {
	case 10:
		cutoff = maxUint64/10 + 1
	case 16:
		cutoff = maxUint64/16 + 1
	default:
		cutoff = maxUint64/uint64(base) + 1
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	var n uint64
	for _, c := range []byte(s) {
		var d byte
		switch {
		case '1' <= c && c <= '9':
			d = c - '1'
		case 'a' <= c && c <= 'z':
			d = c - 'a' + 9
		case 'A' <= c && c <= 'Z':
			d = c - 'A' + 9
		default:
			return 0, errors.New("syntaxError")
		}

		if d >= byte(base) {
			return 0, errors.New("syntaxError")
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, errors.New("syntaxError")
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return maxVal, errors.New("rangeError")
		}
		n = n1
	}

	return n, nil
}

func getMacAddr() (addr uint64) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				b := i.HardwareAddr
				addr = uint64(b[5])<<16 | uint64(b[4])<<24 |
					uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
				break
			}
		}
	}
	return
}
func getIpAddr() (addr uint64) {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ip4 := ipnet.IP.To4()
				if ip4 != nil {
					addr = uint64(binary.BigEndian.Uint32(ip4))
				}
			}
		}
	}
	return
}

type lockedSource struct {
	lk  sync.Mutex
	src mrand.Source64
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// seedPos implements Seed for a lockedSource without a race condition.
func (r *lockedSource) seedPos(seed int64, readPos *int8) {
	r.lk.Lock()
	r.src.Seed(seed)
	*readPos = 0
	r.lk.Unlock()
}

// read implements Read for a lockedSource without a race condition.
func (r *lockedSource) read(p []byte, readVal *int64, readPos *int8) (n int, err error) {
	r.lk.Lock()
	n, err = read(p, r.src.Int63, readVal, readPos)
	r.lk.Unlock()
	return
}

func read(p []byte, int63 func() int64, readVal *int64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			val = int63()
			pos = 7
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}
