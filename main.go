package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type Data struct {
	Data1       string
	Data2       string
	Correlation float64
}

var tpl *template.Template

func init() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/result", Result)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func buildURL(key1 string) string {
	return fmt.Sprintf("https://kabutan.jp/stock/kabuka?code=%s",
		url.QueryEscape(key1))
}

func Result(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	target1 := r.FormValue("target1")
	target2 := r.FormValue("target2")
	data1, err := getData(ctx, target1)
	if err != nil {
		log.Fatalln(err)
	}

	data2, err := getData(ctx, target2)
	if err != nil {
		log.Fatalln(err)
	}

	data1json, _ := json.Marshal(data1)
	data2json, _ := json.Marshal(data2)

	c := correlation(data1, data2)
	ret := &Data{
		Data1:       string(data1json),
		Data2:       string(data2json),
		Correlation: c,
	}
	tpl.ExecuteTemplate(w, "result.gohtml", ret)
}

func getData(ctx context.Context, target string) ([]float64, error) {
	url := buildURL(target)
	client := urlfetch.Client(ctx)

	result, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	reader := csv.NewReader(result.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	const (
		dateColumn  = 0
		closeColumn = 1
	)
	var data []float64
	for _, row := range records {
		val, _ := strconv.ParseFloat(row[closeColumn], 64)
		data = append(data, val)
	}
	return relativize(data), nil
}
