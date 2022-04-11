package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path"
	"strconv"

	cookiemonster "github.com/MercuryEngineering/CookieMonster"
	utils "github.com/meetcw/blc/utils"
)

type Bilibili struct {
	client *http.Client
	id     int64
	root   string
	cookie string
	lux    string
}

func NewBilibili(root string, cookie string, lux string) (*Bilibili, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}

	cookies, err := cookiemonster.ParseFile(cookie)
	if err != nil {
		return nil, err
	}
	var id int64
	for _, item := range cookies {
		if item.Name == "DedeUserID" {
			id, err = strconv.ParseInt(item.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if id == 0 {
		return nil, &ApiError{Message: "Invalid id"}
	}
	return &Bilibili{client, id, root, cookie, lux}, nil
}

func (this *Bilibili) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	cookies, err := cookiemonster.ParseFile(this.cookie)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	return request, nil
}

func (this *Bilibili) GetString(url string, headers map[string]string) (string, error) {
	bytes, err := this.GetBytes(url, headers)
	return string(bytes), err
}

func (this *Bilibili) GetBytes(url string, headers map[string]string) ([]byte, error) {
	request, err := this.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (this *Bilibili) FetchFavorites() (*FavoritesResponse, error) {

	url := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/folder/created/list-all?up_mid=%d", this.id)

	request, err := this.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data FavoritesResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (this *Bilibili) FetchFavoriteMediasByPage(favoriteID int, page int) (*MediasResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/resource/list?media_id=%d&pn=%d&ps=20&keyword=&order=mtime&type=0&tid=0&platform=web&jsonp=jsonp", favoriteID, page)

	request, err := this.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data MediasResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (this *Bilibili) FetchFavoriteMedias(favoriteID int) (*[]Media, error) {
	var videos []Media
	page := 1
	for {
		response, err := this.FetchFavoriteMediasByPage(favoriteID, page)
		if err != nil {
			return nil, err
		}
		if response.Code != 0 {
			return nil, &ApiError{Message: response.Message}
		}
		for _, media := range response.Data.Medias {
			videos = append(videos, Media{
				Identity:     media.BVID,
				UpperName:    media.Upper.Name,
				UpperID:      media.Upper.MemberID,
				Introduction: media.Introduction,
				Cover:        media.Cover,
				Folder:       response.Data.Info.Title,
				Title:        media.Title,
				Active:       media.Attr == 0,
			})
		}
		if response.Data.HasMore {
			page++
		} else {
			break
		}
	}
	return &videos, nil
}

func (this *Bilibili) FetchCollections(page int) (*CollectionsResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/v3/fav/folder/collected/list?pn=%d&ps=20&up_mid=%d&platform=web&jsonp=jsonp", page, this.id)

	request, err := this.newRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data CollectionsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (this *Bilibili) FetchCollectionMediasByPage(collectedID int, page int) (*MediasResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/space/fav/season/list?season_id=%d&pn=%d&ps=20&jsonp=jsonp", collectedID, page)

	request, err := this.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data MediasResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (this *Bilibili) FetchCollectionMedias(collectedID int) (*[]Media, error) {
	var videos []Media
	page := 1
	for {
		response, err := this.FetchCollectionMediasByPage(collectedID, page)
		if err != nil {
			return nil, err
		}
		if response.Code != 0 {
			return nil, &ApiError{Message: response.Message}
		}
		for _, media := range response.Data.Medias {
			videos = append(videos, Media{
				Identity:     media.BVID,
				UpperName:    media.Upper.Name,
				UpperID:      media.Upper.MemberID,
				Introduction: media.Introduction,
				Cover:        media.Cover,
				Folder:       response.Data.Info.Title,
				Title:        media.Title,
				Active:       media.Attr == 0,
			})
		}

		if response.Data.HasMore {
			page++
			// TODO: 目前 B 站有 bug，调用一次就可以获取所有数据
			break
		} else {
			break
		}
	}
	return &videos, nil
}

func (this *Bilibili) FetchAllMedias() ([]Media, error) {
	response, err := this.FetchFavorites()
	if err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, &ApiError{Message: response.Message}
	}
	var videos []Media
	for _, favorite := range response.Data.List {
		list, err := this.FetchFavoriteMedias(favorite.ID)
		if err != nil {
			log.Printf("获取收藏夹 `%s` 中的视频失败：%s", favorite.Title, err)
			continue
		}
		videos = append(videos, *list...)
	}
	page := 1
	for {
		response, err := this.FetchCollections(page)
		if err != nil {
			return nil, err
		}
		if response.Code != 0 {
			return nil, &ApiError{Message: response.Message}
		}
		for _, collection := range response.Data.List {
			if collection.Type == 11 {
				list, err := this.FetchFavoriteMedias(collection.ID)
				if err != nil {
					log.Printf("获取收藏夹 `%s` 中的视频失败：%s", collection.Title, err)
					continue
				}
				videos = append(videos, *list...)
			} else {
				list, err := this.FetchCollectionMedias(collection.ID)
				if err != nil {
					log.Printf("获取收藏的合集 `%s` 中的视频失败：%s", collection.Title, err)
					continue
				}
				videos = append(videos, *list...)
			}

		}

		if response.Data.HasMore {
			page++
		} else {
			break
		}
	}
	return videos, nil
}
func (this *Bilibili) FetchMediaPages(id string) (*PageResponse, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/player/pagelist?bvid=%s&jsonp=jsonp", id)
	request, err := this.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := this.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data PageResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (this *Bilibili) isPlaylist(id string) (bool, error) {
	response, err := this.FetchMediaPages(id)
	if err != nil {
		return false, err
	}
	if response.Code != 0 {
		return false, &ApiError{Message: response.Message}
	}
	return len(response.Data) > 1, nil
}

func (this *Bilibili) MakeDownloadDirectory(media Media) string {
	folder := utils.Pathnamify(media.Folder)

	mediaFolder := utils.Pathnamify(media.Title)

	directory := path.Join(this.root, folder, mediaFolder)
	os.MkdirAll(directory, 0733)
	return directory
}

func (this *Bilibili) SaveMetaData(video Media) error {

	directory := this.MakeDownloadDirectory(video)
	os.MkdirAll(directory, 0733)
	err := os.WriteFile(path.Join(directory, utils.Filenamify(video.Title, "md")), []byte(fmt.Sprintf("# [%s](https://www.bilibili.com/video/%s) \n\nUP：[%s](https://space.bilibili.com/%d) \n\n简介：%s \n\n封面：[cover](%s)", video.Title, video.Identity, video.UpperName, video.UpperID, video.Introduction, video.Cover)), 0733)
	return err
}

func (this *Bilibili) Download(video Media) error {

	directory := this.MakeDownloadDirectory(video)
	isPlaylist, err := this.isPlaylist(video.Identity)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if isPlaylist {
		log.Printf("正在下载合集 `%s`\n", video.Title)
		cmd = exec.Command(this.lux, "-C", "-p", "--cookie", this.cookie, "-o", directory, video.Identity)
	} else {
		log.Printf("正在下载视频 `%s`\n", video.Title)
		cmd = exec.Command(this.lux, "-C", "--cookie", this.cookie, "-o", directory, video.Identity)
	}
	log.Println(cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return err
}
