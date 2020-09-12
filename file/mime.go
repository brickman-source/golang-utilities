package file

import "strings"

func DetectMimeType(fileName string) string {
	if strings.HasSuffix(fileName, ".jpg") {
		return `image/jpeg`
	} else if strings.HasSuffix(fileName, ".png") {
		return `image/png`
	} else if strings.HasSuffix(fileName, ".gif") {
		return `image/gif`
	} else if strings.HasSuffix(fileName, ".webp") {
		return `image/webp`
	} else if strings.HasSuffix(fileName, ".cr2") {
		return `image/x-canon-cr2`
	} else if strings.HasSuffix(fileName, ".tif") {
		return `image/tiff`
	} else if strings.HasSuffix(fileName, ".bmp") {
		return `image/bmp`
	} else if strings.HasSuffix(fileName, ".jxr") {
		return `image/vnd.ms-photo`
	} else if strings.HasSuffix(fileName, ".psd") {
		return `image/vnd.adobe.photoshop`
	} else if strings.HasSuffix(fileName, ".ico") {
		return `image/x-icon`
	} else if strings.HasSuffix(fileName, ".mp4") {
		return `video/mp4`
	} else if strings.HasSuffix(fileName, ".m4v") {
		return `video/x-m4v`
	} else if strings.HasSuffix(fileName, ".mkv") {
		return `video/x-matroska`
	} else if strings.HasSuffix(fileName, ".webm") {
		return `video/webm`
	} else if strings.HasSuffix(fileName, ".mov") {
		return `video/quicktime`
	} else if strings.HasSuffix(fileName, ".avi") {
		return `video/x-msvideo`
	} else if strings.HasSuffix(fileName, ".wmv") {
		return `video/x-ms-wmv`
	} else if strings.HasSuffix(fileName, ".mpg") {
		return `video/mpeg`
	} else if strings.HasSuffix(fileName, ".flv") {
		return `video/x-flv`
	} else if strings.HasSuffix(fileName, ".mid") {
		return `audio/midi`
	} else if strings.HasSuffix(fileName, ".mp3") {
		return `audio/mpeg`
	} else if strings.HasSuffix(fileName, ".m4a") {
		return `audio/m4a`
	} else if strings.HasSuffix(fileName, ".ogg") {
		return `audio/ogg`
	} else if strings.HasSuffix(fileName, ".flac") {
		return `audio/x-flac`
	} else if strings.HasSuffix(fileName, ".wav") {
		return `audio/x-wav`
	} else if strings.HasSuffix(fileName, ".amr") {
		return `audio/amr`
	} else if strings.HasSuffix(fileName, ".epub") {
		return `application/epub+zip`
	} else if strings.HasSuffix(fileName, ".zip") {
		return `application/zip`
	} else if strings.HasSuffix(fileName, ".tar") {
		return `application/x-tar`
	} else if strings.HasSuffix(fileName, ".rar") {
		return `application/x-rar-compressed`
	} else if strings.HasSuffix(fileName, ".gz") {
		return `application/gzip`
	} else if strings.HasSuffix(fileName, ".bz2") {
		return `application/x-bzip2`
	} else if strings.HasSuffix(fileName, ".7z") {
		return `application/x-7z-compressed`
	} else if strings.HasSuffix(fileName, ".xz") {
		return `application/x-xz`
	} else if strings.HasSuffix(fileName, ".pdf") {
		return `application/pdf`
	} else if strings.HasSuffix(fileName, ".exe") {
		return `application/x-msdownload`
	} else if strings.HasSuffix(fileName, ".swf") {
		return `application/x-shockwave-flash`
	} else if strings.HasSuffix(fileName, ".rtf") {
		return `application/rtf`
	} else if strings.HasSuffix(fileName, ".eot") {
		return `application/octet-stream`
	} else if strings.HasSuffix(fileName, ".ps") {
		return `application/postscript`
	} else if strings.HasSuffix(fileName, ".sqlite") {
		return `application/x-sqlite3`
	} else if strings.HasSuffix(fileName, ".nes") {
		return `application/x-nintendo-nes-rom`
	} else if strings.HasSuffix(fileName, ".crx") {
		return `application/x-google-chrome-extension`
	} else if strings.HasSuffix(fileName, ".cab") {
		return `application/vnd.ms-cab-compressed`
	} else if strings.HasSuffix(fileName, ".deb") {
		return `application/x-deb`
	} else if strings.HasSuffix(fileName, ".ar") {
		return `application/x-unix-archive`
	} else if strings.HasSuffix(fileName, ".Z") {
		return `application/x-compress`
	} else if strings.HasSuffix(fileName, ".lz") {
		return `application/x-lzip`
	} else if strings.HasSuffix(fileName, ".rpm") {
		return `application/x-rpm`
	} else if strings.HasSuffix(fileName, ".elf") {
		return `application/x-executable`
	} else if strings.HasSuffix(fileName, ".doc") {
		return `application/msword`
	} else if strings.HasSuffix(fileName, ".docx") {
		return `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
	} else if strings.HasSuffix(fileName, ".xls") {
		return `application/vnd.ms-excel`
	} else if strings.HasSuffix(fileName, ".xlsx") {
		return `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
	} else if strings.HasSuffix(fileName, ".ppt") {
		return `application/vnd.ms-powerpoint`
	} else if strings.HasSuffix(fileName, ".pptx") {
		return `application/vnd.openxmlformats-officedocument.presentationml.presentation`
	} else if strings.HasSuffix(fileName, ".woff") {
		return `application/font-woff`
	} else if strings.HasSuffix(fileName, ".woff2") {
		return `application/font-woff`
	} else if strings.HasSuffix(fileName, ".ttf") {
		return `application/font-sfnt`
	} else if strings.HasSuffix(fileName, ".otf") {
		return `application/font-sfnt`
	} else {
		return "application/octet-stream"
	}
}
