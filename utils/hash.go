// Much of this file is directly taken from https://github.com/theArtechnology/fiber-inertia/blob/master/hashDir.go
// That project does not have a license.

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

// hashByte returns the MD5 hash of the given byte slice.
func hashByte(content []byte) string {
	hasher := md5.New()
	hasher.Write(content)

	return hex.EncodeToString(hasher.Sum(content))
}

// HashDir returns the MD5 hash of the given directory.
func HashDir(dir string) string {
	var finHash = dir

	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		sbyte       := []byte(finHash)
		concatBytes := hashByte(sbyte)
		nameByte    := []byte(path)
		nameHash    := hashByte(nameByte)

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		fileBytes := make([]byte, reflect.TypeOf(int32(0)).Size())
		if _, err := io.ReadFull(file, fileBytes); err != nil {
			return err
		}

		fileHash := hashByte(fileBytes)
		finHash = concatBytes + fileHash + nameHash

		return nil
	})

	c := []byte(finHash)
	m := hashByte(c)

	return m
}
