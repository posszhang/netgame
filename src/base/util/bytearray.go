package util

import (
	"encoding/binary"
	"io"
)

type ByteArray struct {
	//数据区
	data []byte
	//小端，默认大端
	endian binary.ByteOrder
	//ByteArray 对象的长度（以字节为单位）。
	length uint32
	//读的位置
	rp uint32
	//写的位置
	wp uint32
}

func NewByteArray() *ByteArray {

	bytes := &ByteArray{
		data:   make([]byte, 1024),
		endian: binary.BigEndian,
		length: 1024,
		rp:     0,
		wp:     0,
	}

	return bytes
}

func (bytes *ByteArray) Little() {
	bytes.endian = binary.LittleEndian
}

func (bytes *ByteArray) RdSize() uint32 {

	if bytes.wp >= bytes.rp {
		return bytes.wp - bytes.rp
	}

	s := bytes.length - bytes.rp
	s += bytes.wp

	return s
}

// 目前剩余可写的大小
func (bytes *ByteArray) WrSize() uint32 {

	// ---r+++w---
	if bytes.wp >= bytes.rp {

		left := bytes.length - bytes.wp
		left += bytes.rp
		return left
	}

	// +++w---r+++
	return bytes.rp - bytes.wp
}

func (bytes *ByteArray) Read(buf []byte) (int, error) {

	size := uint32(len(buf))

	if bytes.RdSize() < size {
		return 0, io.EOF
	}

	//直接返回
	if bytes.rp+size <= uint32(len(bytes.data)) {
		n := copy(buf, bytes.data[bytes.rp:])
		bytes.rp += uint32(n)
		return n, nil
	}

	n := copy(buf, bytes.data[bytes.rp:])
	left := size - uint32(n)
	n += copy(buf[n:], bytes.data[0:left])

	bytes.rp = left

	return n, nil
}

func (bytes *ByteArray) Write(buf []byte) (int, error) {

	ws := uint32(len(buf))

	if bytes.wp+ws <= bytes.length {

		n := copy(bytes.data[bytes.wp:], buf)
		bytes.wp += uint32(n)

		return n, nil
	}

	left := bytes.WrSize()

	//环形buf写
	if left >= ws {

		w1 := copy(bytes.data[bytes.wp:], buf)
		w2 := copy(bytes.data[0:], buf[w1:])

		bytes.wp = uint32(w2)

		return w1 + w2, nil
	}

	//重设缓冲
	bytes.resize((bytes.length + ws) * 2)
	return bytes.Write(buf)
}

func (bytes *ByteArray) ReadBool() (ret bool, err error) {
	err = binary.Read(bytes, bytes.endian, &ret)
	return ret, err
}

func (bytes *ByteArray) ReadByte() (ret byte, err error) {
	err = binary.Read(bytes, bytes.endian, &ret)
	return ret, err
}

func (bytes *ByteArray) ReadInt16() (ret int16, err error) {
	err = binary.Read(bytes, bytes.endian, &ret)
	return ret, err
}

func (bytes *ByteArray) ReadInt32() (ret int32, err error) {
	err = binary.Read(bytes, bytes.endian, &ret)
	return ret, err
}

func (bytes *ByteArray) ReadInt64() (ret int64, err error) {
	err = binary.Read(bytes, bytes.endian, &ret)
	return ret, err
}

func (bytes *ByteArray) ReadUTF8() string {
	return ""
}

func (bytes *ByteArray) WriteBool(value bool) {
	binary.Write(bytes, bytes.endian, &value)
}

func (bytes *ByteArray) WriteByte(value byte) {

	binary.Write(bytes, bytes.endian, &value)
}

func (bytes *ByteArray) WriteInt16(value int16) {

	binary.Write(bytes, bytes.endian, &value)
}

func (bytes *ByteArray) WriteInt32(value int32) {

	binary.Write(bytes, bytes.endian, &value)
}

func (bytes *ByteArray) WriteInt64(value int64) {

	binary.Write(bytes, bytes.endian, &value)
}

func (bytes *ByteArray) WriteUTF8(value string) {

}

func (bytes *ByteArray) resize(size uint32) {

	if bytes.length >= size {
		return
	}

	buf := make([]byte, size)

	if bytes.wp >= bytes.rp {

		// 无环

		n := copy(buf, bytes.data[bytes.rp:bytes.wp])
		bytes.wp = uint32(n)

	} else {

		w1 := copy(buf, bytes.data[bytes.rp:])
		w2 := copy(buf[w1:], bytes.data[0:bytes.wp])
		bytes.wp = uint32(w1 + w2)
	}

	bytes.data = buf
	bytes.length = size
	bytes.rp = 0
}
