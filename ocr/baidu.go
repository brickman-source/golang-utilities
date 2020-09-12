package ocr

//
//import (
//	"encoding/base64"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"os"
//)
//
//func baiduRec(baseUrl, path string) ([]byte, error) {
//	filebytes, _ := ioutil.ReadFile(path)
//	sourcestring := base64.StdEncoding.EncodeToString(filebytes)
//	token := getToken()
//	urlstr := fmt.Sprintf("%s?access_token=%s", baseUrl, token)
//	//todo options参数抽出来
//	params := url.Values{
//		"image": {sourcestring},
//	}
//	res, err := http.PostForm(urlstr, params)
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//	//fmt.Println(string(body))
//	return body, err
//}
//
//func GeneralBasic(path string) ([]byte, error) {
//	//通用文字识别
//	general_basic := "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
//	return baiduRec(general_basic, path)
//}
//
//func AccurateBasic(path string) ([]byte, error) {
//	//通用文字识别（高精度版）
//	accurate_basic := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"
//	return baiduRec(accurate_basic, path)
//}
//
//func General(path string) ([]byte, error) {
//	//通用文字识别（含位置信息版）
//	general := "https://aip.baidubce.com/rest/2.0/ocr/v1/general"
//	return baiduRec(general, path)
//}
//
//func Accurate(path string) ([]byte, error) {
//	//通用文字识别（含位置高精度版）
//	accurate := "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate"
//	return baiduRec(accurate, path)
//}
//
//func GeneralEnhanced(path string) ([]byte, error) {
//	//通用文字识别（含生僻字版）
//	general_enhanced := "https://aip.baidubce.com/rest/2.0/ocr/v1/general_enhanced"
//	return baiduRec(general_enhanced, path)
//}
//
//func Webimage(path string) ([]byte, error) {
//	//网络图片文字识别
//	webimage := "https://aip.baidubce.com/rest/2.0/ocr/v1/webimage"
//	return baiduRec(webimage, path)
//}
//
//////////////////
//
//
//const TOKEN_URL_FORMAT = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"
//
///**
// * 判断文件是否存在  存在返回 true 不存在返回false
// */
//func checkFileIsExist(filename string) bool {
//	var exist = true
//	if _, err := os.Stat(filename); os.IsNotExist(err) {
//		exist = false
//	}
//	return exist
//}
//
////获取配置信息
//func getConfig(path string) (config map[string]string, err error) {
//	if checkFileIsExist(path) {
//		fmt.Println("已存在配置文件")
//		//存在配置文件就读取
//		f, err1 := os.Open(path)
//		defer f.Close()
//		if err1 != nil {
//			return nil, err1
//		}
//		br := bufio.NewReader(f)
//		config = make(map[string]string)
//		for {
//			str, _, eof := br.ReadLine()
//			if eof == io.EOF {
//				break
//			} else {
//				config_line := strings.Split(string(str), ":")
//				config[config_line[0]] = config_line[1]
//			}
//		}
//		return config, nil
//	} else {
//		fmt.Println("不存在配置文件")
//	}
//	return
//}
//
//func WriteMaptoFile(configs map[string]string, filePath string) error {
//	f, err := os.Create(filePath)
//	if err != nil {
//		fmt.Printf("create map file error: %v\n", err)
//		return err
//	}
//	defer f.Close()
//
//	w := bufio.NewWriter(f)
//	for k, v := range configs {
//		lineStr := fmt.Sprintf("%s:%s", k, v)
//		fmt.Fprintln(w, lineStr)
//	}
//	return w.Flush()
//}
//
//func getToken() string {
//	//检查配置文件中是否存在token
//	configFile := "config.txt"
//	configs, err := getConfig(configFile)
//	checkError(err)
//	if expires, ok := configs["expires"]; ok {
//		lastTime, _ := strconv.Atoi(configs["time"])
//		expireTime, _ := strconv.Atoi(expires)
//		if int64(lastTime)+int64(expireTime) > time.Now().Unix() {
//			//还没过期
//			fmt.Println("还没有过期")
//			return configs["token"]
//		}
//	} else {
//		fmt.Println("配置文件token已经过期")
//	}
//
//	//否则重新请求
//	tokenUrl := fmt.Sprintf(TOKEN_URL_FORMAT, CLIENT_ID, CLIENT_SECRET)
//	resp, err := http.Get(tokenUrl)
//	checkError(err)
//	data, err := ioutil.ReadAll(resp.Body)
//	checkError(err)
//	var tokens interface{}
//	json.Unmarshal(data, &tokens)
//	result := tokens.(map[string]interface{})
//	if errStr, ok := result["error"].(string); ok {
//		log.Fatal(errors.New(result["error_description"].(string) + errStr))
//	} else {
//		//写入配置文件
//		token := result["access_token"].(string)
//		configs := make(map[string]string)
//		configs["expires"] = strconv.Itoa(int(result["expires_in"].(float64)))
//		configs["time"] = strconv.Itoa(int(time.Now().Unix()))
//		configs["token"] = token
//		WriteMaptoFile(configs, configFile)
//		return token
//	}
//	return ""
//}
