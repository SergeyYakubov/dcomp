package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"

	"os"
	"syscall"

	"archive/tar"
	"fmt"
	"io"

	"crypto/sha1"
	"encoding/hex"

	"crypto/md5"
	"crypto/sha256"
	"hash"

	"errors"
	"path/filepath"

	"strings"

	"github.com/spaolacci/murmur3"
	"gopkg.in/yaml.v2"
)

func ReadYaml(fname string, config interface{}) error {
	yamlFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return err
	}
	return nil
}

func CompressString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}

func MapToStruct(m map[string]interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}

const (
	Setuid os.FileMode = 1 << (12 - 1 - iota)
	Setgid
	Sticky
	UserRead
	UserWrite
	UserExecute
	GroupRead
	GroupWrite
	GroupExecute
	OtherRead
	OtherWrite
	OtherExecute
)

func MkdirAllWithCh(path string, perm os.FileMode, uid, gid int) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = MkdirAllWithCh(path[0:j-1], perm, uid, gid)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.

	err = os.Mkdir(path, perm)

	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			err2 := os.Chown(path, uid, gid)
			if err2 != nil {
				return os.Chmod(path, perm)
			}
			return err2
		}
		return err
	}
	err = os.Chown(path, uid, gid)
	if err == nil {
		return os.Chmod(path, perm)
	}
	return err
}

func WriteFile(fname string, b *bytes.Buffer) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, b)
	return err
}

func WriteUnpackedTGZ(dest string, b *bytes.Buffer) error {

	gzf, err := gzip.NewReader(b)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)
	for true {

		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir: // = directory
			os.MkdirAll(dest+name, os.FileMode(header.Mode))
		case tar.TypeReg: // = regular file
			file, err := os.Create(dest + name)
			if err != nil {
				return err
			}
			defer file.Close()

			err = file.Chmod(os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.CopyN(file, tarReader, header.Size)
			if err != nil {
				return err
			}

			/*			data := make([]byte, header.Size)
						_, err := tarReader.Read(data)
						if err != nil {
							return err
						}
						ioutil.WriteFile(dest+name, data, header.Mode)*/
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
	return nil
}

func ReadFile(fname string, compressed bool) (b *bytes.Buffer, err error) {
	b = new(bytes.Buffer)

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//select whether to write via compressor or not
	var w io.Writer
	if compressed {
		gz := gzip.NewWriter(b)
		w = gz
		defer gz.Close()
	} else {
		w = b
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func FileCheckSum(fname, alg string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var hasher hash.Hash

	switch alg {
	case "sha1":
		hasher = sha1.New()
	case "murmur":
		hasher = murmur3.New128()
	case "sha2":
		hasher = sha256.New()
	case "md5":
		hasher = md5.New()
	default:
		return "", errors.New("Unknown algorithm")

	}

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func copyFile(source, dest, mode string, fileInfo os.FileInfo, uid, gid int) error {

	if !CanReadFile(fileInfo, uid, gid) {
		return errors.New("permission denied for " + source)
	}

	// copy only regular files
	if !fileInfo.Mode().IsRegular() {
		return nil
	}

	destPath := filepath.Dir(dest)

	if err := MkdirAllWithCh(destPath, 0777, uid, gid); err != nil {
		return err
	}

	switch mode {
	case "mount":
		return os.Link(source, dest)
	case "copy":
		if err := copyFileContents(source, dest); err != nil {
			return err
		}
		if err := os.Chown(dest, uid, gid); err != nil {
			return err
		}
		return os.Chmod(dest, fileInfo.Mode())

	default:
		return errors.New("CopyFile: Unknown mode")
	}
	return nil
}

func CanReadFile(fileInfo os.FileInfo, uid, gid int) bool {
	fm := fileInfo.Mode()
	if fm&(1<<2) != 0 {
		return true
	} else if (fm&(1<<5) != 0) && (gid ==
		int(fileInfo.Sys().(syscall.Stat_t).Gid)) {
		return true
	} else if (fm&(1<<8) != 0) && (uid ==
		int(fileInfo.Sys().(syscall.Stat_t).Uid)) {
		return true
	}
	return false
}

func CopyPath(source, dest, mode string, uid, gid int) error {
	fileInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return copyDir(source, dest, mode, uid, gid)
	} else {
		return copyFile(source, dest, mode, fileInfo, uid, gid)
	}

}

func copyDir(source, dest, mode string, uid, gid int) error {
	files, err := getDirFiles(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		destfile := GetUploadName(file, source, dest, false)

		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		if err := copyFile(file, destfile, mode, fileInfo, uid, gid); err != nil {
			return err
		}
	}
	return nil
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	return nil
}

func getDirFiles(dir string) (listFiles []string, err error) {
	listFiles = make([]string, 0)

	var scan = func(path string, fi os.FileInfo, err error) (e error) {

		if err != nil {
			return err
		}
		if fi.IsDir() {
			if strings.HasPrefix(fi.Name(), ".") && fi.Name() != "." && fi.Name() != ".." {
				return filepath.SkipDir
			}
		} else {
			if strings.HasPrefix(fi.Name(), ".") {
				return nil
			}

			listFiles = append(listFiles, path)

		}
		return nil
	}

	if err = filepath.Walk(dir, scan); err != nil {
		return
	}

	return

}

func GetUploadName(localname, inipath, destdir string, isdir bool) string {

	rel, _ := filepath.Rel(inipath, localname)
	if rel == localname {
		return destdir
	}
	return filepath.Join(destdir, rel)
}
