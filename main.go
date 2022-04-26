package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/meetcw/blc/api"
)

func run(root string, cookie string, lux string) error {
	bilibili, err := api.NewBilibili(root, cookie, lux)
	if err != nil {
		return err
	}
	videos, err := bilibili.FetchAllMedias()
	if err != nil {
		return err
	}
	log.Printf("获取到 %d 个视频\n", len(videos))
	activeVideos := []api.Media{}
	inactiveVideos := []api.Media{}
	for _, video := range videos {
		// log.Println(video.Folder + "/" + video.Title)
		if video.Active {
			activeVideos = append(activeVideos, video)
		} else {
			inactiveVideos = append(inactiveVideos, video)
		}
	}

	for index, video := range activeVideos {
		log.Printf("[%d/%d (%d)]开始下载 `%s`\n", index+1, len(activeVideos), len(videos), video.Title)
		err = bilibili.Download(video)
		if err != nil {
			log.Printf("\n下载失败 `%s`\n %s\n\n", video.Title, err)
			continue
		}
		log.Printf("\n下载成功 `%s`\n\n", video.Title)

		log.Printf("保存元数据 `%s`\n", video.Title)
		err = bilibili.SaveMetaData(video)
		if err != nil {
			log.Printf("\n保存失败 `%s`\n %s\n\n", video.Title, err)
			continue
		}
		log.Printf("保存成功 `%s`\n", video.Title)
	}
	return nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cookie := flag.String("cookie", "cookie", "Cookie 文件路径")
	lux := flag.String("lux", "lux", "Lux 可执行文件路径")
	interval := flag.Int("interval", 3600, "扫描间隔，单位秒")
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("获取当前目录失败：", err)
	}
	root := flag.String("root", cwd, "Directory for download")
	flag.Parse()
	for {
		log.Println("开始执行...")
		err = run(*root, *cookie, *lux)
		if err != nil {
			log.Println("执行失败：", err)
		} else {
			log.Println("执行成功")
		}

		log.Printf("下次执行时间：%s", time.Now().Add(time.Duration(*interval)*time.Second))
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
