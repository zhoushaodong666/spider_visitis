package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

//增加的访问量
var num int = 0

func main() {

	var waitgroup sync.WaitGroup

	//需要刷访问量的文章列表
	article_list := []string{
		"https://blog.csdn.net/weixin_36162966/article/details/90634035",
		"https://blog.csdn.net/weixin_36162966/article/details/86649197",
		"https://blog.csdn.net/weixin_36162966/article/details/86645765",
		"https://blog.csdn.net/weixin_36162966/article/details/86648219",
		"https://blog.csdn.net/weixin_36162966/article/details/90605065",
		"https://blog.csdn.net/weixin_36162966/article/details/90752923",
		"https://blog.csdn.net/weixin_36162966/article/details/91432201",
		"https://blog.csdn.net/weixin_36162966/article/details/91383006",
		"https://blog.csdn.net/weixin_36162966/article/details/88866796",
		"https://blog.csdn.net/weixin_36162966/article/details/91463689",
	}

	//快代理的页数
	//开始页数
	startPage := 251
	//截止页数
	endPage := 300

	for i := startPage; i <= endPage; i++ {
		time.Sleep(2 * time.Second)
		fmt.Println("开始执行第" + strconv.Itoa(i) + "页")
		waitgroup.Add(1)
		allProxy := GetProxy(i)
		go func(i int) {

			if allProxy != nil {
				for _, oneProxy := range allProxy {
					fmt.Print(oneProxy, "----")
					rand.Seed(time.Now().UnixNano())
					client := NewHttpClient(oneProxy)
					ranNum := rand.Intn(len(article_list))

					err := HttpGET(client, article_list[ranNum])
					if err == nil {
						alist := append(article_list[:ranNum], article_list[ranNum+1:]...)
						for _, v := range alist {
							HttpGET(client, v)
						}
					}
				}

				fmt.Println("第" + strconv.Itoa(i) + "页执行完毕")

			}

			waitgroup.Done()
		}(i)
	}

	waitgroup.Wait()

}

/**
抓取快代理的IP https://www.kuaidaili.com/free
抓取一页代理的IP 因为快代理有反爬机智 太快就返回503状态码  没有ip的数据
*/
func GetProxy(pageIndex int) []string {
	allProxy := []string{}
	url := "https://www.kuaidaili.com/free/inha/" + strconv.Itoa(pageIndex)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode == 503 {
		fmt.Println("爬取快代理第" + strconv.Itoa(pageIndex) + "页失败！")
		return nil
	}
	//关闭body
	defer resp.Body.Close()
	//从body中读取数据
	data, _ := ioutil.ReadAll(resp.Body)
	//匹配IP
	reg_ip := regexp.MustCompile(`<td data-title="IP">(.*)</td>`)
	allIp := reg_ip.FindAllStringSubmatch(string(data), -1)

	//匹配PROT
	reg_prot := regexp.MustCompile(`<td data-title="PORT">(.*?)</td>`)
	allPort := reg_prot.FindAllStringSubmatch(string(data), -1)

	for key, ip := range allIp {
		groupStr := "http://" + ip[1] + ":" + allPort[key][1] + "/"
		allProxy = append(allProxy, groupStr)
	}
	return allProxy
}

//生成一个client
func NewHttpClient(proxyAddr string) *http.Client {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		return nil
	}
	//设置头信息
	header := make(http.Header)
	header["User-Agent"] = []string{GetRandUA()}
	header["Connection"] = []string{"keep-alive"}
	header["User-Agent"] = []string{GetRandUA()}
	header["Accept-Language"] = []string{"zh-CN,zh;q=0.9"}
	header["Upgrade-Insecure-Requests"] = []string{"1"}
	header["Cookie"] = []string{`BAIDU_SSP_lcr=https://www.baidu.com/link?url=DEn9RlgKH6p1K4Avgft4DNNIJf97laz_XgfDb2qoRQ4YPk_g5m59pM07fk3mf7DNSKLZbl9D7M77GEsjwB3Mrsw-vo3sJAgTt7n6IWXp4Ee&wd=&eqid=d840b3a2000032a1000000065d075309; ARK_ID=JSc22425b0dd9e005733148b426f793200c224; smidV2=20180920150921cda58f67c278544839ea082b55f9852a004b6d1e8f48bd990; gr_user_id=81b6631c-e813-46fc-b011-39e394d1c549; _dg_id.40e39cb6d36d5282.c482=60ef5b4272863541%7C%7C%7C1539053046%7C%7C%7C2%7C%7C%7C1539053046%7C%7C%7C1539053046%7C%7C%7C%7C%7C%7Ca77cf9609c0dcd23%7C%7C%7Chttps%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3Dq2M0iPoPtFfxChjWneDR16Kru4JR0Ja__sFk7p6BMv2XmgY-T1QvCpwTVlD3KCicjGauH-XsjKhkfrlWJ7iRydxhv58x78g9LmNe6oVLKNS%26wd%3D%26eqid%3Dad0614dc00006cf7000000065bbc159b%7C%7C%7Chttps%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3Dq2M0iPoPtFfxChjWneDR16Kru4JR0Ja__sFk7p6BMv2XmgY-T1QvCpwTVlD3KCicjGauH-XsjKhkfrlWJ7iRydxhv58x78g9LmNe6oVLKNS%26wd%3D%26eqid%3Dad0614dc00006cf7000000065bbc159b%7C%7C%7C1%7C%7C%7Cundefined; pt_7cd998c4=uid=aYnnP2QYpIzF3eRR3XIF-w&nid=1&vid=p-AlrZJ8WZzyG7Y4UZ9Kzw&vn=1&pvn=1&sact=1539053092760&to_flag=0&pl=nqfQkfCgcIDcfW87CmcPNA*pt*1539053046306; UN=weixin_36162966; __yadk_uid=cNc4IAfl78TBtlwTjuBB729vT9126luc; uuid_tt_dd=10_28867322960-1540794791804-345450; _ga=GA1.2.973882576.1541646120; ADHOC_MEMBERSHIP_CLIENT_ID1.0=9c021c12-8952-7bfd-ecdd-c78f781689bc; Hm_ct_6bcd52f51e9b3dce32bec4a3997715ac=1788*1*PC_VC!5744*1*weixin_36162966!6525*1*10_28867322960-1540794791804-345450; CNZZDATA1259587897=885232125-1553757005-https%253A%252F%252Fwww.baidu.com%252F%7C1553757005; UM_distinctid=16b06d85730945-0503ae4a93ef05-7a1437-1fa400-16b06d85731936; UserName=weixin_36162966; UserInfo=eeb980df3925434da47da24ce27abea6; UserToken=eeb980df3925434da47da24ce27abea6; UserNick=%E9%A3%8E%E9%9B%AA%E4%B9%8B%E6%97%85; AU=6D0; BT=1560309064855; p_uid=U000000; dc_session_id=10_1560413480933.885304; firstDie=1; dc_tos=pt8iyv; Hm_lvt_6bcd52f51e9b3dce32bec4a3997715ac=1560755050,1560760564,1560760583,1560761096; Hm_lpvt_6bcd52f51e9b3dce32bec4a3997715ac=1560761096`}

	netTransport := &http.Transport{
		Proxy:              http.ProxyURL(proxy),
		ProxyConnectHeader: header,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Duration(5*time.Second))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   -1,                              //每个host最大空闲连接
		ResponseHeaderTimeout: time.Duration(10 * time.Second), //数据收发5秒超时
	}

	return &http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: netTransport,
	}
}

//发起http请求
func HttpGET(client *http.Client, url string) error {
	fmt.Println(url)
	rsp, err := client.Get(url)
	if err != nil {
		fmt.Println()
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK || err != nil {
		fmt.Errorf("HTTP GET Code=%v, URI=%v, err=%v", rsp.StatusCode, url, err)
		return err
	}

	num += 1
	fmt.Printf("访问 %s --- 状态码[%s] 访问量:%d \n", url, rsp.StatusCode, num)
	return nil
}

//伪造一个User-Agent更不容易给反爬机制发现
//返回一个随机的User-Agent
func GetRandUA() string {
	//uaMap := make(map[string][]string)
	//休眠300ms 不然for执行过快 会出现一直取到一个特定元素的情况
	//间隔时间越长让随机数更加准确
	time.Sleep(200 * time.Millisecond)
	userAgentList := []string{
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0`,
		`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50`,
		`Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50`,
		`Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 2.0.50727; SLCC2; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; Tablet PC 2.0; .NET4.0E)`,
		`Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; InfoPath.3)`,
		`Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; GTB7.0)`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)`,
		`Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1)`,
		`Mozilla/5.0 (Windows; U; Windows NT 6.1; ) AppleWebKit/534.12 (KHTML, like Gecko) Maxthon/3.0 Safari/534.12`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)`,
		`Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.472.33 Safari/534.3 SE 2.X MetaSr 1.0`,
		`Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)`,
		`Mozilla/5.0 (Windows NT 6.1) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/13.0.782.41 Safari/535.1 QQBrowser/6.9.11079.201`,
		`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E) QQBrowser/6.9.11079.201`,
		`Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)`,
	}
	//str := []string{}
	rand.Seed(time.Now().UnixNano())
	ua := userAgentList[rand.Intn(len(userAgentList))]
	//str = append(str, ua)
	//uaMap["User-Agent"] = str
	return ua
}
