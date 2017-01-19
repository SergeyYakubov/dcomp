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
