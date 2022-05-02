package files

import (
	"archive/zip"
	"bytes"
	"net/http"
)

// CheckOffice 2007及以后版本(OfficeOpenXML格式)
func CheckOffice(b []byte) string {
	if http.DetectContentType(b) != "application/zip" {
		return ""
	}
	archive, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return ""
	}
	for _, f := range archive.File {
		switch f.Name {
		case "word/document.xml":
			return "docx"
		case "xl/workbook.xml":
			return "xlsx"
		case "ppt/presentation.xml":
			return "pptx"
		}
	}
	return ""
}
