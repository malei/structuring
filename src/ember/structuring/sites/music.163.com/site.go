package m1c

import (
	"ember/structuring/types"
	"net/http"
	"regexp"
	"strings"
	"io"
	"os"
	//"fmt"
)

func (p *Song) Run(appender types.Appender) (err error) {
	// TODO
	if len(p.url) > 100 {
		return
	}
	ret, err := p.Crawl()
	if ret == nil || err != nil {
		return  err
	}

	task := types.NewTaskInfo(p.url, "song", 0)
	for _, v := range ret {
		task.Url = "http://" + p.site.domain + "/" + v
		err = appender(task)
	}

	return err
}

func (p *Song) Crawl() (ret []string, err error) {
	body, err := p.site.FetchHtml(p.url)
	if err != nil {
		return nil, err
	}
	pv, err := p.site.ParseHtml(body)
	if pv == nil || err != nil {
		return nil, err
	}
	err = p.site.Write(p.url, pv)
	if err != nil {
		println(err.Error())
		return
	}
	return p.site.ExtractUrl(body)
}

type Song struct {
	site *Site
	url string
}

func (p *Site) NewTask(info types.TaskInfo) types.Task {
	switch info.Type {
	}
	return &Song{p, info.Url}
}

func (p *Site) Close()(err error) {
	p.data.file.Flush()
	return p.data.file.Close()
}

func (p *Site) FetchHtml(url string) (ret []byte, err error) {
	cookie, err := p.GetCookie()
	// TODO check err 
	p.html.cookie = cookie
	return p.html.fetch(url)
}

type SongInfo struct {
	Version string
	Url string
	SongName, Singer, Album, IssueDate string
	IssueCompany, Note, SongLyric string
}

func (p *Site) ParseHtml(body []byte) (ret []string, err error) {
	return p.html.parse(body)
}

func (p *Site) ExtractUrl(body []byte) (ret []string, err error) {
	return p.url.extract(body)
}

func (p *Site) Write(url string, ret []string) (err error) {
	str := p.version + "\t" + url
	for _, v := range ret {
		str = str + "\t" + v
	}
	str = str + "\n"
	return p.data.write([]byte(str), 0)
}

func (p *Site) Search(key string) (ret [][]string, err error) {
	scanner, err := p.data.file.OpenScanner()
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`[^\t\n]+`)
	var x [][]string
	for {
		buf, err := scanner.Scan()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		line := string(buf)
		word := reg.FindAllString(line, -1)
		if strings.Contains(word[2], key) {
			x = append(x, word)
		}
	}
	scanner.Close()
	return x, err
}

func (p *Site) Serialize() (ret []byte, err error) {
	return ret, err
}

func New(root string) (p *Site, err error) {
	domain := "music.163.com"
	path := root + "/" + domain
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return
	}
	p = &Site{domain, "01", NewUrl(), NewHtml(), NewCrawl(), NewData(path)}
	return
}

func (p *Site) GetCookie() (cookie string, err error) {
	for i := 0; i < 3; i++ {
		resp, err := http.Head(p.domain)
		if err != nil {
			continue
		}
		cookie = resp.Header.Get("Set-Cookie")
		break
	}
	return cookie, err
}

type Site struct {
	domain string
	version string
	url Url
	html Html
	crawl Crawl
	data Data
}
