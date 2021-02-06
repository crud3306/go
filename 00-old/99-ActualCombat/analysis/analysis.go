package main 

import (
	"flag"
	"time"
	"os"
	"bufio"
	"io"
	"strings"
	"net/url"
	"crypto/md5"
	"encodeing/hex"
	"strconv"
	// "net"

	// "crypto/dsa"
	// "golang.org/x/net/html/atom"
	"github.com/sirupsen/logrus"
	"github.com/mgutz/str"
	// "github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/pool"
)

const HANDEL_DIG = "/dig?"
const HANDEL_MOVIE = "/movie/"
const HANDEL_LIST = "/list/"
const HANDEL_HTML = ".html"

// 接收的命令行参数
type cmdParams struct {
	logFilePath string
	routineNum int
}

// 日志数据
type digData struct {
	time string
	url  string
	refer string
	ua string
}

// 每个资源的信息
type urlData struct {
	data digData
	uid  string
	unode urlNode
}

// 资源结构
type urlNode struct {
	unType string // 详情页\列表页\首页
	unRid int // resouce id 资源id
	unUrl string // 当前该页面的url
	unTime string // 当前访问该页面的时间
}

// 存储结构
type storageBlock struct {
	counterType string
	storageModel string
	unode		 urlNode
}

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	// 获取参数
	logFilePath := flag.String("logFilePath", "/data/tmp/log/dig.log")
	routineNum  := flag.Int("routineNum", 5, "consumer num by goroutine")
	l := flag.String("l", "/data/tmp/log/dig_runtime.log", "this programe runtime log target file path")
	flag.Parse()

	// flag 取的到值，使用时记得前面带*
	params := cmdParams{*logFilePath, *routineNum}

	// 打日志
	logFd, err := os.OpenFile(*l, os.O_CREATE|os.WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infof("exec start")
	log.Infof("params: logFilePath=%s, routineNum=%d", params.logFilePath, params.routineNum)

	// reids pool
	redisPool, err := pool.New("tcp", "localhost:6379", 2*params.routineNum)
	if err != nil {
		log.Fatalln("redis pool created failed.")
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3*time.Second)
			}
		}()
	}

	// 初始化一些channel，用于数据传递
	var logChannel = make(chan string, 3*params.routineNum)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)

	// 日志消费者
	go readFileLineByLine(params, logChannel)

	// 创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	// 创建pv uv 统计器
	go pvCounter(pvCounter, storageChannel)
	go uvCounter(uvCounter, storageChannel, redisPool)

	// 创建存储器
	go dataStorege(storageChannel, redisPool)

	time.Sleep(1000*time.Sleep)
}

func dataStorege(storageChannel chan storageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.counterType+"_"

		// 天、小时、分钟
		// 层级：顶级/列表/页面
		// 存储结构 redis sortedset
		setKeys := []string{
			prefix+"day_"+getTime(block.unode.unTime, "day"),
			prefix+"hour_"+getTime(block.unode.unTime, "hour"),
			prefix+"min_"+getTime(block.unode.unTime, "min"),
			prefix+block.unode.unType+"_day_"+getTime(block.unode.unTime, "day"),
			prefix+block.unode.unType+"_hour_"+getTime(block.unode.unTime, "hour"),
			prefix+block.unode.unType+"_min_"+getTime(block.unode.unTime, "min"),
		}

		rowId := block.unnode.unRid

		for _, key := range setKeys{
			ret, err := redisPool.Cmd(block.storageModel, key, 1, rowId).Int()
			if ret <= 0 || err != nil {
				log.Errorln("dataStorege redis storate error", block.storageModel, key, rowId)
			}
		}
	}
}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {
	for data := range pvChannel {
		sItem := storageBlock{
			"pv"
			"ZINCREBY",
			data.unode
		}
		storageChannel <- sItem
	}
}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock, redisPool *pool.Pool) {
	for data := range pvChannel {
		// HyperLoglog redis
		hyperLogLogKey := "uv_hpll_"+getTime(data.data.time, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, data.uid, "EX", 86400).Int()
		if err != nil {
			log.Warningln("uvcounter check redis hyperloglog failed")
		}
		if ret != 1 {
			continue
		}

		sItem := storageBlock{
			"uv"
			"ZINCREBY",
			data.unode
		}
		storageChannel <- sItem
	}
}

// 从logChannel中，消费每条日志
func logConsumer(logChannel chan string, pvChannel chan urlData, uvChannel chan urlData) error {
	for logStr := range logChannel {
		// 切割日志每行的数据，匹配需要的数据
		data := cutLogFetchData(logStr)

		//uid
		// 正常情况下的日志，可以用cookie唯一值
		// 这里因是假数据，临时用md5(refer+ua)
		hasher := md5.New()
		hasher.Write([]byte(data.refer+data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))

		// 很多更复杂的解析工作都可以在这里完成
		// ...

		uData := urlData{ data, uid, formatUrl(data.url, data.time) }

		// log.Infoln(uData)

		pvChannel <- urlData
		uvChannel <- urlData
	}

	return nil
}

// 根据url与time，拼装urlNode
func formatUrl(url, t string) urlNode {
	// 从量大的着手，详情>列表>首页
	pos1 := str.IndexOf(url, HANDEL_MOVIE, 0)
	if pos1 != -1 {
		pos1 += len(HANDEL_MOVIE)
		pos2 := str.IndexOf(url, HANDEL_HTML, 0)
		idStr := str.Substr(url, pos1, pos2-pos1)
		id, _ := strconv.Atoi(idStr)
		return urlNode{"movie", id, url, t}
	} else {
		pos1 = str.IndexOf(url, HANDEL_LIST 0)
		if pos1 != -1 {
			pos1 += len(HANDEL_LIST)
			pos2 := str.IndexOf(url, HANDEL_HTML, 0)
			idStr := str.Substr(url, pos1, pos2-pos1)
			id, _ := strconv.Atoi(idStr)
			return urlNode{"list", id, url, t}

		} else {

			return urlNode{"home", 0, url, t}
		} // 如果页面url 有很多种，就不断在这里扩展
	}
}

// 逐行分析原始日志，提取digData
func cutLogFetchData(logStr String) digData {
	logStr = strings.TrimSpace(logStr)
	pos1 := str.IndexOf(logStr, HANDEL_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}

	pos1 += len(HANDEL_DIG)
	pos2 := str.IndexOf(logStr, " HTTP/", pos1)
	d := str.Substr(logStr, pos1, pos2-pos1)

	urlInfo, err := url.Parse("http://localhost/?"+d)
	if err != nil {
		return digData{}
	}

	data := urlInfo.Query()
	return digData{
		data.Get("time")
		data.Get("refer")
		data.Get("url")
		data.Get("ua")
	}
}

// 逐行读日志，并写入logChannel
func readFileLineByLine(params cmdParams, logChannel chan string) error {
	fd, err = os.Open(params.logFilePath)
	if err != nil {
		log.Warningf("readFileLineByLine can not open file: %s", params.logFilePath)
		return err
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString("\n")
		logChannel <- line
		count++

		// 如果每条都记log，就太多了，这里设为每读1000条log，才记一次
		if count%(1000*params.routineNum) {
			log.Infof("readFileLineByLine line: %d", count)
		}
		if err != nil {
			if err == io.EOF {
				time.Sleep(3*time.Second)
				log.Infof("readFileLineByLine wait, readline: %d", count)
			} else {
				log.Warningf("readFileLineByLine read log error on line:%d", count)
			}
		}
	}

	return nil
}

func getTime(logTime, timeType string) string {
	var item string
	switch timeType {
		case "day":
			item = "2006-01-02"
		case "hour":
			item = "2006-01-02 15"
		case "min":
			item = "2006-01-02 15:04"
	}
	t, _ := time.Parse(item, time.Now().Format(item))
	return strconv.FormatInt(t.Unix(), 10)
}



























