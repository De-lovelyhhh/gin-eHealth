package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func padding(src []byte, blocksize int) []byte {
	padnum := blocksize - len(src)%blocksize
	pad := bytes.Repeat([]byte{byte(padnum)}, padnum)
	return append(src, pad...)
}

func unpadding(src []byte) []byte {
	n := len(src)
	if n == 0 {
		return nil
	}
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

//AES加密,CBC
func AesEncrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	src = padding(src, block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block, key)
	blockmode.CryptBlocks(src, src)
	return src
}

//AES解密
func AesDecrypt(src []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src
}

func MakeBlocksFull(src []byte, blockSize int) []byte {

	//1. 获取src的长度， blockSize对于des是8
	length := len(src)

	//2. 对blockSize进行取余数， 4
	remains := length % blockSize

	//3. 获取要填的数量 = blockSize - 余数
	paddingNumber := blockSize - remains //4

	//4. 将填充的数字转换成字符， 4， '4'， 创建了只有一个字符的切片
	//s1 = []byte{'4'}
	s1 := []byte{byte(paddingNumber)}

	//5. 创造一个有4个'4'的切片
	//s2 = []byte{'4', '4', '4', '4'}
	s2 := bytes.Repeat(s1, paddingNumber)

	//6. 将填充的切片追加到src后面
	s3 := append(src, s2...)

	return s3
}
