package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func post(w http.ResponseWriter, r *http.Request, d *Data) {
	ctx := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(ctx, "Result", nil)
	if _, err := datastore.Put(ctx, key, d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ranking(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	datas := make([]*Data, 0, 40)
	q := datastore.NewQuery("Result").Order("-correlation").Limit(40)
	for it := q.Run(ctx); ; {
		var post Data
		_, err := it.Next(&post)
		if err == datastore.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		datas = append(datas, &post)
	}
	ret := WithoutDepulicate(datas)
	if err := tpl.ExecuteTemplate(w, "ranking.html", ret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func WithoutDepulicate(datas []*Data) []*Data {
	result := []*Data{}

	for i := 0; i < len(datas); i++ {
		exists := false
		for v := 0; v < i; v++ {
			if datas[v].Correlation == datas[i].Correlation {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, datas[i])
		}
	}
	return result
}
