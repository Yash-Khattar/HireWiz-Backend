package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
)

func ValidatePDFFile(file *multipart.FileHeader) error {
	// Check file extension
	if filepath.Ext(file.Filename) != ".pdf" {
		return errors.New("only PDF files are allowed")
	}

	// Check file size (e.g., max 5MB)
	if file.Size > 5*1024*1024 {
		return errors.New("file size exceeds 5MB limit")
	}

	return nil
} 