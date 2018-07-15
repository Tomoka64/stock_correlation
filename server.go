package stock_correlation

import (
	"errors"
	"net/http"

	"google.golang.org/appengine"
)

type targetData struct {
	target string
	data   []float64
}

func Index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func Result(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	cdata := make(chan targetData)

	datas := []targetData{
		{
			target: r.FormValue("target1"),
		},
		{
			target: r.FormValue("target2"),
		},
	}

	for _, v := range datas {
		go getData(ctx, v, cdata)
	}

	data1 := <-cdata

	data2 := <-cdata

	if str, err := errorHandling(data1, data2); err != nil {
		tpl.ExecuteTemplate(w, "error.html", str)
		return
	}

	new1, new2, err := Adjust(data1.data, data2.data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	leng := len(new1)
	covariance, sd1, sd2, corre := Correlation(new1, new2)
	word := Related(corre)

	ret := &Data{
		Data1:                    data1.target,
		Data2:                    data2.target,
		Covariance:               covariance,
		StandardDeviationOfData1: sd1,
		StandardDeviationOfData2: sd2,
		Correlation:              corre,
		Explanation:              word,
		Compared:                 leng,
	}

	tpl.ExecuteTemplate(w, "result.html", ret)
	post(w, r, ret)
}

func errorHandling(data1, data2 targetData) (string, error) {
	var ret string
	if len(data1.data) == 0 && len(data2.data) == 0 {
		ret = data1.target + " and " + data2.target
		return ret, errors.New("length of the data is 0")
	}
	if len(data1.data) == 0 {
		return data1.target, errors.New("length of the data is 0")
	}
	if len(data2.data) == 0 {
		return data2.target, errors.New("length of the data is 0")
	}
	return "", nil
}
