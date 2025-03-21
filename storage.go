package reqparser

import (
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"io"
	"os"
)

var defaultPathFile = "storage/"

var allowedFileExt = []string{
	".png", ".jpg", ".jpeg", // image type
	".ogv", ".jpm", ".mp4", ".webm", ".mpg", // video type
	".mpeg", ".mpe", ".mpv", ".ogg", ".qt", ".3gp", ".flv", ".swf", // video type
	".avi", ".mov", ".wmv", ".yuv", ".rm", ".rmvb", ".xlsx", ".zip", // misc type
	".7z", ".docx", ".pptx", ".csv", ".gz", ".pdf", // document type
}

func SaveFileToStorage(c FiberInterface, fieldName string, required bool) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if required {
			return "", NewCommonError(fieldName, "NOT_BLANK")
		}
	}

	if file == nil {
		return "", nil
	}

	theFile, err := file.Open()
	if err != nil {
		return "", err
	}

	mime, err := mimetype.DetectReader(theFile)
	if err != nil {
		return "", err
	}

	if !slicesContain(allowedFileExt, mime.Extension()) {
		return "", NewCommonError(fieldName, "EXTENSION_NOT_ALLOWED")
	}

	fileName := uuid.NewString() + mime.Extension()
	filePath := defaultPathFile + fileName

	storageFile, err := os.Open(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		} else {
			storageFile, err = os.Create(filePath)
			if err != nil {
				return "", err
			}
		}
	}

	_, err = io.Copy(storageFile, theFile)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	return fileName, nil
}

func SaveFileToStorageWithExt(c FiberInterface, fieldName string, required bool, allowedExt []string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if required {
			return "", NewCommonError(fieldName, "NOT_BLANK")
		}
	}

	if file == nil && !required {
		file = nil
		return "", nil
	}

	theFile, err := file.Open()
	if err != nil {
		return "", err
	}

	extension, err := CheckFileExtension(theFile, fieldName, allowedExt)
	if err != nil {
		return "", err
	}

	fileName := uuid.NewString() + extension
	filePath := "storage/files/" + fileName

	storageFile, err := os.Open(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		} else {
			storageFile, err = os.Create(filePath)
			if err != nil {
				return "", err
			}
		}
	}

	_, err = io.Copy(storageFile, theFile)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	return fileName, nil
}

func CheckFileExtension(theFile io.ReadSeeker, fieldName string, allowedExtension []string) (string, error) {
	mimetype, err := mimetype.DetectReader(theFile)
	if err != nil {
		return "", err
	}

	extension := mimetype.Extension()
	if !slicesContain(allowedExtension, extension) {
		return "", NewCommonError(fieldName, "EXTENSION_NOT_ALLOWED")
	}

	_, err = theFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return extension, nil
}

func ModifyDefaultPathFile(path string) {
	defaultPathFile = path
}

func ModifyAllowedFileExt(extensions []string) {
	allowedFileExt = extensions
}

func slicesContain[T comparable](s []T, contain T) bool {
	for _, v := range s {
		if v == contain {
			return true
		}
	}

	return false
}
