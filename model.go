package main

type Data struct {
	Data1                    string  `datastore:"data1"`
	Data2                    string  `datastore:"data2"`
	Covariance               float64 `datastore:"covariance"`
	StandardDeviationOfData1 float64 `datastore:"sd1"`
	StandardDeviationOfData2 float64 `datastore:"sd2"`
	Correlation              float64 `datastore:"correlation"`
	Explanation              string  `datastore:"explantion"`
	Compared                 int     `datastore:"compared"`
}
