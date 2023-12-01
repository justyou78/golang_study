package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
}

var baseUrl string = "https://www.saramin.co.kr/zf_user/search/recruit?searchword=python"

func main() {
	var jobs []extractedJob

	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		// extractedJobs... : [] + [] + [] => []
		jobs = append(jobs, extractedJobs...)
	}

	fmt.Println(jobs)
}

func getPage(i int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseUrl + "&recruitPage=" + strconv.Itoa(i+1) // https://www.saramin.co.kr/zf_user/search/recruit?searchword=python&recruitPage=1

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCard := doc.Find(".item_recruit")
	searchCard.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	location := ""

	id, _ := card.Attr("value")
	title, _ := card.Find(".job_tit > a").Attr("title")

	card.Find(".job_condition>span>a").Each(func(i int, locations_part *goquery.Selection) {
		location += cleanString(locations_part.Text())
	})

	job := extractedJob{id: id, title: title, location: location}

	return job

}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseUrl)
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
