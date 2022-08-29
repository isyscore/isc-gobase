package compress

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//// Compress 压缩文件

func Zip(zipPath string, paths []string) error {
	// create zip file
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer func(archive *os.File) {
		_ = archive.Close()
	}(archive)

	// new zip writer
	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	// traverse the file or directory
	for _, srcPath := range paths {
		// remove the trailing path separator if path is a directory
		srcPath = strings.TrimSuffix(srcPath, string(os.PathSeparator))

		// visit all the files or directories in the tree
		err = filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// create a local file header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// set compression
			header.Method = zip.Deflate

			// set relative path of a file as the header name
			header.Name, err = filepath.Rel(filepath.Dir(srcPath), path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				header.Name += string(os.PathSeparator)
			}

			// create writer for the file header and save content of the file
			headerWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			err = copyFile(path, headerWriter)
			if err != nil {
				return err
			}
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(fileSrcPath string, dstWriter io.Writer) error {
	f, err := os.Open(fileSrcPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, err = io.Copy(dstWriter, f)
	return err
}

// Unzip decompresses a zip file to specified directory.
// Note that the destination directory don't need to specify the trailing path separator.
func Unzip(zipPath, dstDir string) error {
	// open zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(w)

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
