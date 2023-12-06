package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
}

func Scrape(term string) {
	var baseUrl string = "https://www.saramin.co.kr/zf_user/search/recruit?searchword=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := getPages(baseUrl)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseUrl, c)
		// // extractedJobs... : [] + [] + [] => []
		// jobs = append(jobs, extractedJobs...)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)

	}

	writeJobs(jobs)
}

func getPage(i int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url + "&recruitPage=" + strconv.Itoa(i+1) // https://www.saramin.co.kr/zf_user/search/recruit?searchword=python&recruitPage=1

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCard := doc.Find(".item_recruit")
	searchCard.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCard.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://www.saramin.co.kr/zf_user/jobs/relay/view?isMypage=no&rec_idx=" + job.id, job.title, job.location}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)

	}

}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	location := ""

	id, _ := card.Attr("value")
	title, _ := card.Find(".job_tit > a").Attr("title")

	card.Find(".job_condition>span>a").Each(func(i int, locations_part *goquery.Selection) {
		location += CleanString(locations_part.Text())
	})

	c <- extractedJob{id: id, title: title, location: location}

}

// CleanString: cleanr string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	fmt.Println(pages)
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}

}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: ", res.StatusCode)
	}
}
