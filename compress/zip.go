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
//// files 文件数组，可以是不同dir下的文件或者文件夹
//// relativePath 压缩包内要去除的相对路径
//// dest 目标压缩文件的路径
//func Compress(dest string, relatedPath string, files []string) error {
//	fileList := isc.ListToMapFrom[string, *os.File](files).Map(func(item string) *os.File {
//		f, _ := os.Open(item)
//		return f
//	})
//
//	d, _ := os.Create(dest)
//	defer func(d *os.File) { _ = d.Close() }(d)
//	w := zip.NewWriter(d)
//	defer func(w *zip.Writer) { _ = w.Close() }(w)
//	for _, file := range fileList {
//		err := compress(file, relatedPath, "", w)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//// Decompress 解压
//// zipFile 压缩文件路径
//// dest 目标解压路径
//func Decompress(zipFile string, dest string) error {
//	reader, err := zip.OpenReader(zipFile)
//	if err != nil {
//		return err
//	}
//	defer func(reader *zip.ReadCloser) { _ = reader.Close() }(reader)
//	for _, file := range reader.File {
//		if rc, err := file.Open(); err != nil {
//			return err
//		} else {
//			filename := filepath.Join(dest, file.Name)
//			if err = os.MkdirAll(getDir(filename), 0755); err != nil {
//				return err
//			}
//			if w, err := os.Create(filename); err != nil {
//				return err
//			} else {
//				if _, err = io.Copy(w, rc); err != nil {
//					return err
//				}
//				_ = w.Close()
//			}
//			_ = rc.Close()
//		}
//	}
//	return nil
//}
//
//func compress(file *os.File, relatedPath string, prefix string, zw *zip.Writer) error {
//	defer func(file *os.File) { _ = file.Close() }(file)
//	if info, err := file.Stat(); err != nil {
//		return err
//	} else {
//		if info.IsDir() {
//			prefix = prefix + "/" + info.Name()
//			if fileInfos, err := file.Readdir(-1); err != nil {
//				return err
//			} else {
//				for _, fi := range fileInfos {
//					if f, err := os.Open(file.Name() + "/" + fi.Name()); err != nil {
//						return err
//					} else {
//						if err = compress(f, relatedPath, prefix, zw); err != nil {
//							return err
//						}
//					}
//				}
//			}
//		} else {
//			if header, err := zip.FileInfoHeader(info); err != nil {
//				return err
//			} else {
//				header.Name = prefix + "/" + header.Name
//				if strings.HasPrefix(header.Name, relatedPath) {
//					header.Name = strings.Replace(header.Name, relatedPath, "", 1)
//				}
//				if strings.HasPrefix(header.Name, "/") {
//					header.Name = strings.Replace(header.Name, "/", "", 1)
//				}
//				if writer, err := zw.CreateHeader(header); err != nil {
//					return err
//				} else {
//					if _, err = io.Copy(writer, file); err != nil {
//						return err
//					}
//				}
//			}
//		}
//	}
//	return nil
//}
//
//func getDir(path string) string {
//	return subString(path, 0, strings.LastIndex(path, "/"))
//}
//
//func subString(str string, start, end int) string {
//	rs := []rune(str)
//	length := len(rs)
//	if start < 0 || start > length {
//		panic("start is wrong")
//	}
//	if end < start || end > length {
//		panic("end is wrong")
//	}
//	return string(rs[start:end])
//}

func Zip(zipPath string, paths []string) error {
	// create zip file
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// new zip writer
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

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
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(headerWriter, f)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip decompresses a zip file to specified directory.
// Note that the destination directory don't need to specify the trailing path separator.
func Unzip(zipPath, dstDir string) error {
	// open zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
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
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
