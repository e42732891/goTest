package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	//reStr = `\d+@qq.com`
	reLink = `href="(https?://[\s\S]+?)"`
	//https://www.baidu.com?tn=SE_hldp03480_530ir7bs
	//url1 =: `https://bcy.net/`
	//https://i0.hdslb.com/bfs/live/8974a6ca62f20449cd330f588e8d663b2eefcbee.jpg@400w_250h.jpg
	//<img alt="国服刘备鬼谷子 在线无敌带粉顺便教学" src="https://i0.hdslb.com/bfs/live/8974a6ca62f20449cd330f588e8d663b2eefcbee.jpg@400w_250h.jpg">
	//reImg = `<img[\s\S]+?>`
	//reImg = `<img[\s\S]*?(src=")([^<>]+?)"[\s\S]+?>` //B站首页图片/tl640

	// BCY URL https://img-bcy-qn.pstatp.com/user/516637/item/web/179i3/a5aae8105e3c11e89b200d4d5bf1446c.jpg/tl640
	reImg      = `<img[\s\S]*?(src=")(https://[^<>]+?)"[\s\S]+?>` //B站首页图片/tl640
	chanImgUrl chan string
)

func HandleError(err error, why string) {
	fmt.Print(why, err)
}

func getPageStr(url string) (pageStr string) {

	http.Cookie{}
	//获取url的html文本字符串
	resp, err := http.Get(url)
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

func SpiderPicImg(url string) (urls []string) {

	pageStr := getPageStr(url)
	//fmt.Println(pageStr)

	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Printf("共找到%d条结果:\n", len(results))

	for _, result := range results {

		//fmt.Println(result[0])

		//fmt.Println(result[2])
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
func main() {
	urls := SpiderPicImg("https://bcy.net/")

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

}
