package file

import (
	"path/filepath"
	"strings"

	"gopkg.in/h2non/filetype.v1"
)

type MediaType struct {
	Mime        string
	Description string
}

// maps media types to descriptions most people would recognize.
var mediaTypeDescriptions = map[string]string{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "Spreadsheet",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "Word",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "Presentation",
	"Composite Document File V2 Document, No summary info":                      "Composite Document File",
	"application/vnd.ms-powerpoint.presentation.macroEnabled.12":                "MS-PowerPoint",
	"text/xml":                                 "XML",
	"image/jpeg":                               "JPEG",
	"application/postscript":                   "Postscript",
	"image/png":                                "PNG",
	"application/json":                         "JSON",
	"image/vnd.ms-modi":                        "MS-Document Imaging",
	"application/vnd.ms-xpsdocument":           "MS-Postscript",
	"image/vnd.radiance":                       "Radiance",
	"application/vnd.sealedmedia.softseal.pdf": "Softseal PDF",
	"application/vnd.hp-PCL":                   "PCL",
	"application/xslt+xml":                     "XSLT",
	"image/gif":                                "GIF",
	"application/matlab":                       "Matlab",
	"application/pdf":                          "PDF",
	"application/xml":                          "XML",
	"application/vnd.ms-excel":                 "MS-Excel",
	"image/bmp":                                "BMP",
	"image/x-ms-bmp":                           "BMP",
	"image/tiff":                               "TIFF",
	"image/vnd.adobe.photoshop":                "Photoshop",
	"application/pkcs7-signature":              "PKCS",
	"image/vnd.dwg":                            "DWG",
	"application/octet-stream":                 "Binary",
	"application/rtf":                          "RTF",
	"text/plain":                               "Text",
	"application/vnd.ms-powerpoint":            "MS-PowerPoint",
	"application/x-troff-man":                  "TROFF",
	"video/x-ms-wmv":                           "WMV Video",
	"application/vnd.chemdraw+xml":             "ChemDraw",
	"text/html":                                "HTML",
	"video/mpeg":                               "MPEG Video",
	"text/csv":                                 "CSV",
	"application/zip":                          "ZIP",
	"application/msword":                       "MS-Word",
	"unknown":                                  "Unknown",
}

func init() {
	filetype.AddType("txt", "text/plain")
	filetype.AddType("m", "application/matlab")
}

func GetMediaTypeByExtension(path string) MediaType {
	var mt MediaType
	ext := strings.ToLower(filepath.Ext(path)[1:])
	mime := filetype.GetType(ext)
	switch mime.Extension {
	case "unknown":
		mt.Description = mediaTypeDescriptions["unknown"]
		mt.Mime = "unknown"
	default:
		mt.Mime = mime.MIME.Value
		desc, ok := mediaTypeDescriptions[mt.Mime]
		if !ok {
			mt.Description = "Unknown"
		} else {
			mt.Description = desc
		}
	}

	return mt
}
