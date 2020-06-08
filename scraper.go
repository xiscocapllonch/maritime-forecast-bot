package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"time"
)

type Result struct {
	Warning  Warning  `xml:"aviso"`
	Forecast Forecast `xml:"prediccion"`
	Trend    Trend    `xml:"tendencia"`
}

type Warning struct {
	End  Date   `xml:"fin"`
	Text string `xml:"texto"`
}

type Forecast struct {
	Start Date `xml:"inicio"`
	End   Date `xml:"fin"`
	Zones []struct {
		Name     string `xml:"nombre,attr"`
		SubZones []struct {
			Name string `xml:"nombre,attr"`
			Text string `xml:"texto"`
		} `xml:"subzona"`
	} `xml:"zona"`
}

type Trend struct {
	Start Date   `xml:"inicio"`
	End   Date   `xml:"fin"`
	Text  string `xml:"texto"`
}

type Date string

func getXML(url string) (Result, error) {
	res, err := http.Get(url)
	if err != nil {
		return Result{}, err
	}


	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Result{}, errors.New("status code error: " + string(res.StatusCode) + " / " + res.Status)

	}

	var result Result

	decoder := xml.NewDecoder(res.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&result)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (d Date) formatDate() string {
	EsWeeksDays := []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	EsMonths := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

	date, err := time.Parse("2006-01-02T15:04:05", string(d))

	if err != nil {
		log.Println("formatDate error: ", err)
		return string(d)
	}

	day := EsWeeksDays[date.Weekday()]
	month := EsMonths[int(date.Month())-1]

	return fmt.Sprintf(
		"%v, %v %v %v a las %v hora oficial",
		day,
		date.Day(),
		month,
		date.Year(),
		date.Format("15:04"),
	)
}

func (w Warning) formatText() string {
	return fmt.Sprintf(
		"<u><b>Avisos válidos hasta el %v</b></u>\n\n%v",
		w.End.formatDate(),
		w.Text,
	)
}

func (f Forecast) formatText() string {
	forecastText := fmt.Sprintf(
		"<u><b>Predicción</b></u>\n\nFecha de inicio: %v\nFecha de fin: %v\n\n\n",
		f.Start.formatDate(),
		f.End.formatDate(),
	)

	for _, zone := range f.Zones {
		forecastText += fmt.Sprintf("<b>%v</b>\n\n", zone.Name)

		if len(zone.SubZones) == 1 {
			forecastText += fmt.Sprintf(
				"%v\n\n\n",
				zone.SubZones[0].Text,
			)
		} else {
			for _, subZone := range zone.SubZones {
				forecastText += fmt.Sprintf(
					"•\t<b>%v</b>: %v\n\n",
					subZone.Name,
					subZone.Text,
				)
			}
		}
	}

	return forecastText
}

func (t Trend) formatText() string {
	return fmt.Sprintf(
		"<u><b>Tendencia</b></u>\n\nFecha de inicio: %v\nFecha de fin: %v\n\n%v\n\n\n",
		t.Start.formatDate(),
		t.End.formatDate(),
		t.Text,
	)
}

func (r Result) formatText() string {
	return fmt.Sprintf(
		"%v\n\n\n\n%v",
		r.Warning.formatText(),
		r.Forecast.formatText(),
	)
}
