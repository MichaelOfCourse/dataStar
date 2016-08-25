package dataStar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MichaelOfCourse/file"
	"github.com/geops/gtfsparser"
)

func getLine(a string, datas jsonData) *busLine {
	for _, b := range datas.BusLines {
		if b.Name == a {
			return b
		}
	}
	return nil
}

func directionExists(dir int, line busLine) bool {
	return line.Path[dir] != nil
}

func getStop(all []*stop, id string) *stop {
	for _, v := range all {
		if v.ID == id {
			return v
		}
	}
	return nil
}

type stop struct {
	ID   string
	Code string
	Lat  float32
	Lon  float32
	Name string
}

type linePath struct {
	Dir   string
	Stops []*stop
}

type busLine struct {
	Name string
	Path [2]*linePath
}

type jsonData struct {
	Date     time.Time
	BusLines []*busLine
}

func order(feed *gtfsparser.Feed, jsonData *jsonData) {
	var newStop *stop
	var allStops []*stop
	var line *busLine

	for _, v := range feed.Stops {
		newStop = new(stop)
		newStop.ID = v.Id
		newStop.Code = v.Code
		newStop.Name = v.Name
		newStop.Lat = v.Lat
		newStop.Lon = v.Lon
		allStops = append(allStops, newStop)
	}

	for _, v := range feed.Trips {
		line = getLine(v.Route.Short_name, *jsonData)
		if line == nil {
			line = new(busLine)
			line.Name = v.Route.Short_name
			jsonData.BusLines = append(jsonData.BusLines, line)
		} else if directionExists(v.Direction_id, *line) {
			continue
		}
		line.Path[v.Direction_id] = new(linePath)
		line.Path[v.Direction_id].Dir = v.Headsign
		for _, j := range v.StopTimes {
			line.Path[v.Direction_id].Stops = append(line.Path[v.Direction_id].Stops, getStop(allStops, j.Stop.Id))
		}
	}
}

// Create provide a file containing json infos
func Create() {
	feed := gtfsparser.NewFeed()
	jsonData := new(jsonData)

	if e := feed.Parse("starGtfs.zip"); e != nil {
		fmt.Println(e)
	}
	jsonData.Date = time.Now()
	order(feed, jsonData)
	b, err := json.Marshal(*jsonData)
	if err != nil {
		fmt.Println(err)
	}
	f, err := file.Create("jsonData.json")
	if err != nil {
		fmt.Println(err)
	}
	if nb, err := f.Write(b); err != nil || nb != len(b) {
		fmt.Println(err)
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
}
