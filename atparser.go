package asiatorrents

import (
	"errors"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"github.com/RangelReale/filesharetop/lib"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ATSort int

const (
	ATSORT_SEEDERS  ATSort = 5
	ATSORT_LEECHERS ATSort = 6
	ATSORT_COMPLETE ATSort = 7
)

type ATSortBy int

const (
	ATSORTBY_ASCENDING  ATSortBy = 1
	ATSORTBY_DESCENDING ATSortBy = 2
)

type ATParser struct {
	List   map[string]*fstoplib.Item
	config *Config
	logger *log.Logger
}

func NewATParser(config *Config, l *log.Logger) *ATParser {
	return &ATParser{
		List:   make(map[string]*fstoplib.Item),
		config: config,
		logger: l,
	}
}

func (p *ATParser) Parse(sort ATSort, sortby ATSortBy, pages int) error {

	if pages < 1 {
		return errors.New("Pages must be at least 1")
	}

	posct := int32(0)
	for pg := 1; pg <= pages; pg++ {
		var doc *gq.Document
		var e error

		// download the page
		u, e := url.Parse(fmt.Sprintf("http://www.asiatorrents.me/index.php?page=torrents&active=0&discount=0&order=%d&by=%d&pages=%d", sort, sortby, pg))
		if e != nil {
			return e
		}

		cookies, _ := cookiejar.New(nil)
		cookies.SetCookies(u, []*http.Cookie{
			&http.Cookie{Name: "PHPSESSID", Value: p.config.PHPSESSID, Path: "/", Domain: "www.asiatorrents.me"},
			&http.Cookie{Name: "lastseen", Value: p.config.Lastseen, Path: "/", Domain: "www.asiatorrents.me"},
			&http.Cookie{Name: "atvpnshown", Value: "YES", Path: "/", Domain: "www.asiatorrents.me"},
			&http.Cookie{Name: "pass", Value: p.config.Pass, Path: "/", Domain: "www.asiatorrents.me"},
			&http.Cookie{Name: "uid", Value: p.config.Uid, Path: "/", Domain: "www.asiatorrents.me"},
		})

		client := &http.Client{
			Jar: cookies,
		}

		req, e := http.NewRequest("GET", u.String(), nil)
		if e != nil {
			return e
		}

		resp, e := client.Do(req)
		if e != nil {
			return e
		}

		// parse the page
		if doc, e = gq.NewDocumentFromResponse(resp); e != nil {
			return e
		}

		// find the ordered column link using the requested order
		valid := doc.Find(fmt.Sprintf("div.b-content table table.lista tr td > a[href^=\"/index.php?page=torrents&active=0&discount=0&order=%d\"]", sort)).First()
		if valid.Length() == 0 {
			return errors.New("Doc not valid")
		}

		// On the ordered column an up or down arrow is added, check if it is present
		if !strings.ContainsAny(valid.Parent().Text(), "\u2191\u2193") {
			return errors.New("Doc not valid 2")
		}

		// Iterate on each record
		doc.Find("div.b-content table table.lista tr").Each(func(i int, s *gq.Selection) {
			var se error

			link := s.Find("td.lista > a[href^=\"index.php?page=torrent-detail\"]:not([href*=\"#comments\"])").First()
			if link.Length() == 0 {
				//p.logger.Println("ERROR: Link not found")
				return
			}

			href, hvalid := link.Attr("href")
			if !hvalid || href == "" {
				p.logger.Println("ERROR: Link not found")
				return
			}

			hu, se := url.Parse(href)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			hu.Scheme = "http"
			hu.Host = "www.asiatorrents.me"

			lid := hu.Query().Get("id")
			if lid == "" {
				p.logger.Println("ERROR: Link not found")
				return
			}

			category := s.Find("td > a[href^=\"index.php?page=torrents&category\"]").First()
			if category.Length() == 0 {
				p.logger.Println("ERROR: Category not found")
				return
			}
			cathref, catvalid := category.Attr("href")
			if !catvalid || cathref == "" {
				p.logger.Println("ERROR: Cat link not found")
				return
			}

			cu, se := url.Parse(cathref)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			catid := cu.Query().Get("category")

			seeder := s.Find("td > a[href^=\"index.php?page=peers\"]").First()
			if seeder.Length() == 0 {
				p.logger.Println("ERROR: Seeder not found")
				return
			}
			leecher := s.Find("td > a[href^=\"index.php?page=peers\"]").Eq(1)
			if leecher.Length() == 0 {
				p.logger.Println("ERROR: Leecher not found")
				return
			}
			complete := s.Find("td > a[href^=\"index.php?page=torrent_history\"]").First()
			/*
				if complete.Length() == 0 {
					p.logger.Printf("ERROR: Complete not found - %s - %s", link.Text(), hu.String())
					return
				}
			*/

			comments := s.Find("td > a[href^=\"index.php?page=torrent-details\"][href*=\"#comments\"]").First()
			/*
				if comments.Length() == 0 {
					p.logger.Printf("ERROR: Comments not found - %s - %s", link.Text(), hu.String())
					return
				}
			*/

			adddate := seeder.Parent().Prev().Prev().Prev()

			nseeder, se := strconv.ParseInt(seeder.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			nleecher, se := strconv.ParseInt(leecher.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			ncomplete := int64(0)
			if complete.Length() > 0 {
				ncomplete, se = strconv.ParseInt(complete.Text(), 10, 32)
				if se != nil {
					p.logger.Printf("ERROR: %s", se)
					return
				}
			}
			ncomments := int64(0)
			if comments.Length() > 0 && comments.Text() != "---" {
				ncomments, se = strconv.ParseInt(comments.Text(), 10, 32)
				if se != nil {
					p.logger.Printf("ERROR: %s", se)
					return
				}
			}
			nadddate, se := time.Parse("02/01/2006", adddate.Text())
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}

			//fmt.Printf("%s: %s\n", link.Text(), hu.Query().Get("id"))
			item, ok := p.List[lid]
			if !ok {
				item = fstoplib.NewItem()
				item.Id = lid
				item.Title = link.Text()
				item.Link = hu.String()
				item.Count = 0
				item.Category = catid
				item.AddDate = nadddate.Format("2006-01-02")
				item.Seeders = int32(nseeder)
				item.Leechers = int32(nleecher)
				item.Complete = int32(ncomplete)
				item.Comments = int32(ncomments)
				p.List[lid] = item
			}
			item.Count++
			posct++
			if sort == ATSORT_SEEDERS {
				item.SeedersPos = posct
			} else if sort == ATSORT_LEECHERS {
				item.LeechersPos = posct
			} else if sort == ATSORT_COMPLETE {
				item.CompletePos = posct
			}
		})
	}

	return nil
}
