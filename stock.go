package stock_correlation

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func Related(result float64) string {
	ab := math.Abs(result)
	if ab >= 0 && ab <= 0.2 {
		return "ほとんど相関なし"
	}
	if ab >= 0.2 && ab <= 0.4 {
		return "弱い相関あり"
	}
	if ab >= 0.4 && ab <= 0.7 {
		return "やや相関あり"
	}
	if ab >= 0.7 && ab <= 1.0 {
		return "かなり強い相関がある"
	}
	return "相関説明不可能な値"
}

func Adjust(data1, data2 []float64) ([]float64, []float64, error) {
	var ret []float64
	if len(data1) > len(data2) {
		ret = data1[:len(data2)]
		return ret, data2, nil
	}
	if len(data2) > len(data1) {
		ret = data2[:len(data1)]
		return data1, ret, nil
	}
	return []float64{}, []float64{}, errors.New("could not make length of slice equal")
}

func getData(ctx context.Context, td targetData, result chan targetData) {
	var data []string
	var counter = 1
	for counter < 20 {
		url := BuildURL(td.target, counter)
		client := urlfetch.Client(ctx)
		resp, err := client.Get(url)
		if err != nil {
			log.Infof(ctx, "could not get the url:", err)
		}
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			log.Infof(ctx, "error happened during goquery:", err)
		}
		doc.Find("table.stock_kabuka1 > tbody > tr").Each(func(i int, s *goquery.Selection) {
			t := s.Find("td").Next().Next().Next().Next().Text()
			data = append(data, t)

			td.data = CloseRow(data)
		})
		counter++
	}
	result <- td
}

func BuildURL(t string, s int) string {
	a := strconv.Itoa(s)
	return fmt.Sprintf("https://kabutan.jp/stock/kabuka?code=%s&ashi=day&page=%s",
		url.QueryEscape(t), url.QueryEscape(a))
}

func CloseRow(s []string) []float64 {
	var ret []float64
	for _, da := range s {
		for i := 0; i < len(da); i++ {
			if da[i] == 43 || da[i] == 45 {
				new := strings.Replace(da[:i], ",", "", -1)
				num, _ := strconv.Atoi(new)
				ret = append(ret, float64(num))
				break
			}
		}
	}
	return ret
}
