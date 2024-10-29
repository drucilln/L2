package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	startURL  = flag.String("url", "", "Начальный URL")
	outputDir = flag.String("output", "output", "Выходная директория")
	maxDepth  = flag.Int("depth", 1, "Максимальная глубина рекурсии")
)

type Crawler struct {
	startUrl  string
	outputDir string
	maxDepth  int
	visited   map[string]struct{}
	visitedMu sync.Mutex
	urlQ      chan URLTask
	wg        sync.WaitGroup
	client    *http.Client
	taskCount int32
	closeOnce sync.Once
}

type URLTask struct {
	url   *url.URL
	depth int
}

func NewCrawler(startURL, outputDir string, maxDepth int) *Crawler {
	return &Crawler{
		startUrl:  startURL,
		outputDir: outputDir,
		maxDepth:  maxDepth,
		visited:   make(map[string]struct{}),
		urlQ:      make(chan URLTask),
		client:    &http.Client{},
		closeOnce: sync.Once{},
	}
}

func (c *Crawler) Run() {
	parsedURL, err := url.Parse(c.startUrl)
	if err != nil {
		log.Fatal(err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for task := range c.urlQ {
			c.processURL(task)
			if atomic.AddInt32(&c.taskCount, -1) == 0 {
				c.closeOnce.Do(func() {
					close(c.urlQ)
				})
			}
		}
	}()

	c.enqueueTask(URLTask{parsedURL, 0})

	c.wg.Wait()
}

func (c *Crawler) enqueueTask(task URLTask) {
	atomic.AddInt32(&c.taskCount, 1)
	c.urlQ <- task
}

func (c *Crawler) processURL(task URLTask) {
	if task.depth > c.maxDepth {
		return
	}

	urlStr := task.url.String()
	c.visitedMu.Lock()
	if _, ok := c.visited[urlStr]; ok {
		c.visitedMu.Unlock()
		return
	}

	c.visited[urlStr] = struct{}{}
	c.visitedMu.Unlock()

	fmt.Printf("Скачивание: %s\n", urlStr)

	resp, err := c.client.Get(urlStr)
	if err != nil {
		log.Printf("Ошибка при загрузке %s: %v", urlStr, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Статус ответа %d для %s", resp.StatusCode, urlStr)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		c.processHTML(task, resp.Body)
	} else {
		c.saveResource(task.url, resp.Body)
	}

	//c.enqueueTask(URLTask{task.url, task.depth + 1})
	//atomic.AddInt32(c.taskCount, -1)
	//if atomic.
}

func (c *Crawler) processHTML(task URLTask, body io.Reader) {
	data, err := io.ReadAll(body)
	if err != nil {
		log.Printf("Ошибка чтения тела ответа: %v", err)
		return
	}

	//filePath := filepath.Join(c.outputDir, task.url.Path)
	filePath := c.getFilePath(task.url)
	c.saveToFile(filePath, data)

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		log.Printf("Ошибка парсинга HTML: %v", err)
		return
	}

	links := c.extractLinks(doc, task.url)
	for _, link := range links {
		absUrl := task.url.ResolveReference(link)
		if !c.isSameDomain(absUrl) {
			continue
		}
		c.addTask(absUrl, task.depth+1)
	}

	resources := c.extractResources(doc, task.url)
	for _, res := range resources {
		absUrl := task.url.ResolveReference(res)
		c.addTask(absUrl, task.depth+1)
	}
}

func (c *Crawler) getFilePath(u *url.URL) string {
	path := u.Path
	if strings.HasSuffix(path, "/") || path == "" {
		path += "index.html"
	}
	filePath := filepath.Join(c.outputDir, u.Hostname(), path)
	return filePath
}

func (c *Crawler) saveResource(url *url.URL, body io.Reader) {
	filePath := c.getFilePath(url)
	data, err := io.ReadAll(body)
	if err != nil {
		log.Printf("Ошибка чтения ресурса %s: %v", url.String(), err)
		return
	}
	c.saveToFile(filePath, data)
}

func (c *Crawler) saveToFile(filePath string, data []byte) {
	err := os.MkdirAll(filepath.Dir(filePath), 0777)
	if err != nil {
		log.Printf("Ошибка создания директорий для %s: %v", filePath, err)
		return
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Ошибка записи файла %s: %v", filePath, err)
		return
	}
}

func (c *Crawler) extractLinks(n *html.Node, base *url.URL) []*url.URL {
	var links []*url.URL
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href := strings.TrimSpace(a.Val)
					if href != "" && !strings.HasPrefix(href, "javascript:") {
						parsedURL, err := url.Parse(href)
						if err == nil {
							links = append(links, parsedURL)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return links
}

func (c *Crawler) isSameDomain(u *url.URL) bool {
	startHost, err := url.Parse(c.startUrl)
	if err != nil {
		return false
	}
	return u.Hostname() == startHost.Hostname()
}

func (c *Crawler) addTask(absUrl *url.URL, depth int) {
	urlStr := absUrl.String()

	c.visitedMu.Lock()
	if _, ok := c.visited[urlStr]; ok {
		c.visitedMu.Unlock()
		return
	}

	c.visited[urlStr] = struct{}{}
	c.visitedMu.Unlock()

	//c.wg.Add(1)
	//go func() {
	//	defer c.wg.Done()
	//	c.urlQ <- URLTask{absUrl, depth}
	//}()
	c.enqueueTask(URLTask{absUrl, depth})
}

func (c *Crawler) extractResources(n *html.Node, base *url.URL) []*url.URL {
	var resources []*url.URL
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrKey string
			switch n.Data {
			case "img", "script":
				attrKey = "src"
			case "link":
				attrKey = "href"
			default:
				attrKey = ""
			}
			if attrKey != "" {
				for _, a := range n.Attr {
					if a.Key == attrKey {
						href := strings.TrimSpace(a.Val)
						if href != "" {
							parsedURL, err := url.Parse(href)
							if err == nil {
								resources = append(resources, parsedURL)
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return resources
}

func main() {
	flag.Parse()

	if startURL == nil || *startURL == "" {
		log.Fatal("Пожалуйста, укажите начальный URL с помощью флага -url")
	}

	err := os.MkdirAll(*outputDir, 0777)
	if err != nil {
		log.Fatalf("Не удалось создать директорию: %v", err)
	}

	crawler := NewCrawler(*startURL, *outputDir, *maxDepth)
	crawler.Run()
}
