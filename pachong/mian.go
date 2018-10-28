package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	//reStr = `\d+@qq.com`
	reLink = `href="(https?://[\s\S]+?)"`
	//https://www.baidu.com?tn=SE_hldp03480_530ir7bs
	url1 = `https://bcy.net/search/home?k=尼禄`
	//bcy首页json获取url
	bcyJsUrl = `https://bcy.net/item/rec/getItemRec?since=0&grid_type=timeline`
	//https://i0.hdslb.com/bfs/live/8974a6ca62f20449cd330f588e8d663b2eefcbee.jpg@400w_250h.jpg
	//<img alt="国服刘备鬼谷子 在线无敌带粉顺便教学" src="https://i0.hdslb.com/bfs/live/8974a6ca62f20449cd330f588e8d663b2eefcbee.jpg@400w_250h.jpg">
	//reImg = `<img[\s\S]+?>`
	//reImg = `<img[\s\S]*?(src=")([^<>]+?)"[\s\S]+?>` //B站首页图片/tl640

	// BCY URL https://img-bcy-qn.pstatp.com/user/516637/item/web/179i3/a5aae8105e3c11e89b200d4d5bf1446c.jpg/tl640
	reImg = `<img[\s\S]*?(src=")(https://[^<>]+?)/w650"[\s\S]+?>` //B站首页图片/tl640

	//目标样例 "since":"6605067423102009607"
	reImgDetails = `"since":"([\s\S]*?)",`
	chanImgUrl   chan string
)

func HandleError(err error, why string) {
	fmt.Print(why, err)
}

func getPageStr(url string) (pageStr string) {
	//绑定 cookie

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	jar, _ := cookiejar.New(nil)
	//	jar.SetCookies(req.URL, []*http.Cookie{
	//		&http.Cookie{Name: "__tea_sdk__ssid", Path: "/", Domain: ".bcy.net", Value: "ea563902-f713-4005-83a5-aa0bb87769a7", HttpOnly: false},
	//		&http.Cookie{Name: "__tea_sdk__user_unique_id", Path: "/", Domain: ".bcy.net", Value: "6616233150227957261", HttpOnly: false},
	//		&http.Cookie{Name: "_csrf_token", Path: "/", Domain: ".bcy.net", Value: "5073c0f14d4d791cab7a6f9f284240c2", HttpOnly: false},
	//		&http.Cookie{Name: "_ga", Path: "/", Domain: ".bcy.net", Value: "GA1.2.15070018.1540461802", HttpOnly: false},
	//		&http.Cookie{Name: "_gat_gtag_UA_121535331_1", Path: "/", Domain: ".bcy.net", Value: "1", HttpOnly: false},
	//		&http.Cookie{Name: "_gid", Path: "/", Domain: ".bcy.net", Value: "GA1.2.78355028.1540461802", HttpOnly: false},
	//		&http.Cookie{Name: "ccid", Path: "/", Domain: ".bcy.net", Value: "891562a25ca4730bc57c571e11400e3a", HttpOnly: false},
	//		&http.Cookie{Name: "Hm_lpvt_330d168f9714e3aa16c5661e62c00232", Path: "/", Domain: ".bcy.net", Value: "1540537889", HttpOnly: false},
	//		&http.Cookie{Name: "Hm_lvt_330d168f9714e3aa16c5661e62c00232", Path: "/", Domain: ".bcy.net", Value: "1540461802", HttpOnly: false},
	//		&http.Cookie{Name: "lang_set", Path: "/", Domain: ".bcy.net", Value: "zh", HttpOnly: false},
	//		&http.Cookie{Name: "mobile_set", Path: "/", Domain: ".bcy.net", Value: "no", HttpOnly: false},
	//		&http.Cookie{Name: "passport_auth_status", Path: "/", Domain: ".bcy.net", Value: "e886861afef29dc6b795140c8ac67e12", HttpOnly: false},
	//		&http.Cookie{Name: "PHPSESSID", Path: "/", Domain: ".bcy.net", Value: "891e08afdb764cb1a1a53a9c5fbe47ad", HttpOnly: false},
	//		&http.Cookie{Name: "sessionid", Path: "/", Domain: ".bcy.net", Value: "ba6334a1219dca0b7967780112112716", HttpOnly: false},
	//		&http.Cookie{Name: "sid_guard", Path: "/", Domain: ".bcy.net", Value: "ba6334a1219dca0b7967780112112716%7C1540537884%7C5184000%7CTue%2C+25-Dec-2018+07%3A11%3A24+GMT", HttpOnly: false},
	//		&http.Cookie{Name: "sid_tt", Path: "/", Domain: ".bcy.net", Value: "ba6334a1219dca0b7967780112112716", HttpOnly: false},
	//		&http.Cookie{Name: "tt_webid", Path: "/", Domain: ".bcy.net", Value: "6616233150227957261", HttpOnly: false},
	//		&http.Cookie{Name: "uid_tt", Path: "/", Domain: ".bcy.net", Value: "a965ef73225e7063a2ee56f02233412c", HttpOnly: false},
	//	})
	client.Jar = jar

	//获取url的html文本字符串
	//resp, err := http.Get(url)
	resp, err := client.Do(req)

	HandleError(err, "http.Get Url")
	defer resp.Body.Close()

	//response.body返回的是一个io.ReadCloser 将它转换成bytes
	pageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")

	//bytes转成字符串
	pageStr = string(pageBytes)
	//fmt.Println(pageStr)
	return pageStr
}

func SpiderPicImg(pageStr string) (urls []string) {

	//fmt.Println(pageStr)
	//正则处理
	re := regexp.MustCompile(reImg)

	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果:\n", len(results))

	for _, result := range results {

		//fmt.Println(result[0])

		fmt.Println(result[2])
		//fmt.Println()
		urls = append(urls, result[2])

	}
	return
}

func SpiderLink() {
	urlStr := "http://tieba.baidu.com/p/2544042204"
	pageStr := getPageStr(urlStr)
	//fmt.Println(pageStr)

	re := regexp.MustCompile(reLink)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果:\n", len(results))

	for _, result := range results {
		fmt.Println(result[1])
	}

}

func DownLoadFile(url string, filename string) (ok bool) {

	resp, err := http.Get(url)
	HandleError(err, "http.Get(url)")
	defer resp.Body.Close()

	fBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll(resp.Body)")

	err = ioutil.WriteFile(filename, fBytes, 0644)
	HandleError(err, "ioutil.WriteFile(filename, fBytes,0644).var")

	if err != nil {
		return false
	} else {
		return true
	}
}

func GetFilenameFromUrl(url string, dirPath string) (filename string) {
	lastIndex := strings.LastIndex(url, "/")
	filename = url[lastIndex+1:]
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + filename
	filename = dirPath + filename
	//fmt.Println(fileName)
	return
}

func getWorkJson(url string) (jsonBodyStr string) {
	resp, err := http.Get(url)
	HandleError(err, "http.Get Url")
	defer resp.Body.Close()
	pageBytes, err := ioutil.ReadAll(resp.Body)
	jsonBodyStr = string(pageBytes)
	return jsonBodyStr
}

func jsonSwitchStringArray(jsonBodyStr string) (urls []string) {

	//正则处理
	re := regexp.MustCompile(reImgDetails)

	results := re.FindAllStringSubmatch(jsonBodyStr, -1)
	fmt.Printf("共找到%d条结果:\n", len(results))

	for _, result := range results {

		//		fmt.Println(result[0])
		//		fmt.Println(result[1])
		//		fmt.Println()
		urls = append(urls, result[1])

	}
	return
}

func getDetailsUrlStr(headUrl string, workIdStrs []string) (urls []string) {

	for _, result := range workIdStrs {
		//fmt.Println(result)
		url := headUrl + result
		urls = append(urls, url)
		fmt.Println("这是获取到的内容url:" + headUrl + url)
	}
	return urls
}

func main() {
	//通过bcy的js url获取首页动态生成的url
	jsonStr := getWorkJson(bcyJsUrl)

	//取出json部分中需要的作品id
	workIdStrs := jsonSwitchStringArray(jsonStr)

	//合并bcy的作品url头和id 合成可用的url组
	urls := getDetailsUrlStr(`https://bcy.net/item/detail/`, workIdStrs)

	//测试--urls的首个作品url 获取此url的页面
	pageStr := getPageStr(urls[0])
	fmt.Println("正在查询ID:" + urls[0])
	//把页面内容处理成图片url []string
	urls2 := SpiderPicImg(pageStr)
	fmt.Println(urls2)
	for _, imgUrl := range urls2 {
		filename := GetFilenameFromUrl(imgUrl, `F:\goGZ\goTest\pachong\img\`)
		ok := DownLoadFile(imgUrl, filename)
		if ok {
			fmt.Println("下载成功")
		} else {
			fmt.Println("下载失败")
		}
		fmt.Println()
	}

	//		for _, url := range urls2 {
	//			pageStr := getPageStr(url)

	//			//SpiderPicImg(pageStr)
	//			urlss := SpiderPicImg(pageStr)

	//			for _, result := range urlss {

	//				result = "" + result
	//				fmt.Println(result)

	//				//这里获取URL的名称
	//				filename := GetFilenameFromUrl(result, `F:\goGZ\goTest\pachong\img\`)

	//				ok := DownLoadFile(result, filename)
	//				if ok {
	//					fmt.Println("下载成功")
	//				} else {
	//					fmt.Println("下载失败")
	//				}
	//				fmt.Println()

	//			}
	//		}

	/* //这个片段是 下载单url的测试
	pageStr := getPageStr(url1)
	SpiderPicImg(pageStr)
	urls := SpiderPicImg(pageStr)

	for _, result := range urls {

		result = "" + result
		fmt.Println(result)

		//这里获取URL的名称
		filename := GetFilenameFromUrl(result, `F:\goGZ\goTest\pachong\img\`)

		ok := DownLoadFile(result, filename)
		if ok {
			fmt.Println("下载成功")
		} else {
			fmt.Println("下载失败")
		}
		fmt.Println()

	}
	*/
}
