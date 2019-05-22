package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/ullaakut/disgo/style"

	"github.com/PuerkitoBio/goquery"
	"github.com/ullaakut/disgo"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const dictionaryURL = "https://community.geniusvision.net/platform/cprndr/manulist"

var rtspURLsFound sync.Map

func main() {
	if err := updateDictionary(); err != nil {
		log.Fatalf(err.Error())
	}
}

func updateDictionary() error {
	disgo.SetTerminalOptions(disgo.WithColors(true), disgo.WithDebug(true))

	disgo.StartStep("Fetching dictionary list")
	resp, err := http.Get(dictionaryURL)
	if err != nil {
		return disgo.FailStepf("unable to download dictionaries: %v", err)
	}
	defer resp.Body.Close()

	disgo.StartStep("Parsing dictionary list")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return disgo.FailStepf("unable to read from dictionary list: %v", err)
	}

	var vendorURLs []string
	doc.Find("td.simpletable a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			return
		}

		if url != "javascript:void(0)" {
			vendorURLs = append(vendorURLs, url)
		}
	})

	disgo.StartStep("Loading current cameradar dictionary")
	currentDictionary, err := ioutil.ReadFile("dictionaries/routes")
	if err != nil {
		return disgo.FailStepf("unable to read current dictionary: %v", err)
	}

	dictionaryEntries := bytes.Split(currentDictionary, []byte("\n"))

	for _, rtspURL := range dictionaryEntries {
		rtspURLsFound.Store(string(rtspURL), struct{}{})
	}

	disgo.Debugf("Current dictionary has %d entries\n", len(dictionaryEntries))
	disgo.EndStep()

	p := mpb.New(mpb.WithWidth(64))
	name := fmt.Sprintf("Fetching default routes from %d constructors:", len(vendorURLs))
	bar := p.AddBar(int64(len(vendorURLs)),
		// set custom bar style, default one is "[=>-]"
		mpb.BarStyle("╢▌▌░╟"),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(name, decor.WC{W: len(name), C: decor.DidentRight}),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	for _, url := range vendorURLs {
		go loadRoutes(url, bar)
	}

	p.Wait()

	disgo.StartStep("Converting found routes into proper data model")

	var rtspURLs []string
	rtspURLsFound.Range(func(rtspURL, _ interface{}) bool {
		disgo.Infoln("Adding URL", rtspURL.(string))
		rtspURLs = append(rtspURLs, rtspURL.(string))
		return true
	})

	sort.Slice(rtspURLs, func(a, b int) bool {
		return rtspURLs[a] < rtspURLs[b]
	})

	disgo.EndStep()

	if len(dictionaryEntries) < len(rtspURLs) {
		disgo.Infof("%s Saving them in cameradar default dictionary.\n", style.Success("Found ", len(rtspURLs)-len(dictionaryEntries), " new entries!"))
		saveRoutes(rtspURLs)
	} else {
		disgo.Infoln(style.Success("No new entry found, dictionary up-to-date! :)"))
	}

	return nil
}

func loadRoutes(url string, bar *mpb.Bar) {
	defer bar.IncrBy(1)

	var (
		failureCounter int
		resp           *http.Response
		err            error
	)
	for failureCounter < 5 {
		resp, err = http.Get(url)
		if err != nil {
			failureCounter++
		} else {
			break
		}
	}

	if failureCounter == 5 {
		disgo.Errorln("Request failed 5 times in a row, giving up on this vendor")
		return
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		disgo.Errorf("unable to read from dictionary list for URL %q: %v\n", url, err)
		return
	}

	doc.Find("tr.simpletable td.simpletable:nth-child(4) a").Each(func(i int, s *goquery.Selection) {
		rtspURL := s.Text()

		if strings.HasPrefix(rtspURL, "(") && strings.HasSuffix(rtspURL, ")") {
			return
		}

		if strings.HasPrefix(rtspURL, "[") && strings.HasSuffix(rtspURL, "]") {
			return
		}

		if strings.HasPrefix(rtspURL, "http://") {
			return
		}

		// Skip the port and only get the route.
		if strings.HasPrefix(rtspURL, "rtsp://ip-addr:") {
			routeAndPort := strings.TrimSpace(strings.TrimPrefix(rtspURL, "rtsp://ip-addr:"))
			route := strings.TrimLeft(routeAndPort, "0123456789/")
			rtspURLsFound.Store(route, struct{}{})
			return
		}

		switch rtspURL {
		case "", "[Details]", "rtsp://ip-addr/", "http://ip-addr/":
			return
		default:
			route := strings.TrimSpace(strings.TrimPrefix(rtspURL, "rtsp://ip-addr/"))
			rtspURLsFound.Store(route, struct{}{})
		}
	})
}

func saveRoutes(rtspURLs []string) {
	contents := strings.Join(rtspURLs, "\n")

	disgo.StartStep("Writing new dictionary file")
	err := ioutil.WriteFile("dictionaries/routes", []byte(contents), 0644)
	if err != nil {
		disgo.FailStepf("unable to write dictionnary: %v", err)
	}
}
