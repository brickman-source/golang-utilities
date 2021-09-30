package http

import (
	"github.com/brickman-source/golang-utilities/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
)

func GetRequestBody(r *http.Request) []byte {
	if r.Method == "POST" || r.Method == "PUT" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			println("ioutil.ReadAll(r.Body) err:", err.Error())
			return nil
		}
		return data
	} else {
		println("method:", r.Method)
	}
	return nil
}

func GetRequestBodyAsMap(r *http.Request) map[string]interface{} {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	ret := make(map[string]interface{})
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}
	return ret
}

func AddURLQuery(reqURL string, key, value string) string {
	x, err := url.Parse(reqURL)
	if err != nil {
		return reqURL
	}
	q := x.Query()
	q.Set(key, value)
	x.RawQuery = q.Encode()
	return x.String()
}

func ContentTypeByUrl(httpUrl string) string {
	ext := ""
	u, err := url.Parse(httpUrl)
	if err == nil {
		ext = filepath.Ext(path.Base(u.Path))
	}
	switch ext {
	case ".jpg":
		return `image/jpeg`
	case ".png":
		return `image/png`
	case ".gif":
		return `image/gif`
	case ".webp":
		return `image/webp`
	case ".cr2":
		return `image/x-canon-cr2`
	case ".tif":
		return `image/tiff`
	case ".bmp":
		return `image/bmp`
	case ".jxr":
		return `image/vnd.ms-photo`
	case ".psd":
		return `image/vnd.adobe.photoshop`
	case ".ico":
		return `image/x-icon`
	case ".mp4":
		return `video/mp4`
	case ".m4v":
		return `video/x-m4v`
	case ".mkv":
		return `video/x-matroska`
	case ".webm":
		return `video/webm`
	case ".mov":
		return `video/quicktime`
	case ".avi":
		return `video/x-msvideo`
	case ".wmv":
		return `video/x-ms-wmv`
	case ".mpg":
		return `video/mpeg`
	case ".flv":
		return `video/x-flv`
	case ".mid":
		return `audio/midi`
	case ".mp3":
		return `audio/mpeg`
	case ".m4a":
		return `audio/m4a`
	case ".ogg":
		return `audio/ogg`
	case ".flac":
		return `audio/x-flac`
	case ".wav":
		return `audio/x-wav`
	case ".amr":
		return `audio/amr`
	case ".epub":
		return `application/epub+zip`
	case ".zip":
		return `application/zip`
	case ".tar":
		return `application/x-tar`
	case ".rar":
		return `application/x-rar-compressed`
	case ".gz":
		return `application/gzip`
	case ".bz2":
		return `application/x-bzip2`
	case ".7z":
		return `application/x-7z-compressed`
	case ".xz":
		return `application/x-xz`
	case ".pdf":
		return `application/pdf`
	case ".exe":
		return `application/x-msdownload`
	case ".swf":
		return `application/x-shockwave-flash`
	case ".rtf":
		return `application/rtf`
	case ".eot":
		return `application/octet-stream`
	case ".ps":
		return `application/postscript`
	case ".sqlite":
		return `application/x-sqlite3`
	case ".nes":
		return `application/x-nintendo-nes-rom`
	case ".crx":
		return `application/x-google-chrome-extension`
	case ".cab":
		return `application/vnd.ms-cab-compressed`
	case ".deb":
		return `application/x-deb`
	case ".ar":
		return `application/x-unix-archive`
	case ".Z":
		return `application/x-compress`
	case ".lz":
		return `application/x-lzip`
	case ".rpm":
		return `application/x-rpm`
	case ".elf":
		return `application/x-executable`
	case ".doc":
		return `application/msword`
	case ".docx":
		return `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
	case ".xls":
		return `application/vnd.ms-excel`
	case ".xlsx":
		return `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
	case ".ppt":
		return `application/vnd.ms-powerpoint`
	case ".pptx":
		return `application/vnd.openxmlformats-officedocument.presentationml.presentation`
	case ".woff":
		return `application/font-woff`
	case ".woff2":
		return `application/font-woff`
	case ".ttf":
		return `application/font-sfnt`
	case ".otf":
		return `application/font-sfnt`
	default:
		return "application/octet-stream"
	}
}
