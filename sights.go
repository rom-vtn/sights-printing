package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"

	_ "github.com/gmlewis/go-fonts-a/fonts/aaarghnormal"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	trainmapdb "github.com/rom-vtn/trainmap-db"
)

func getSights(fetchUrl string) ([]trainmapdb.RealTrainSight, error) {
	httpResp, err := http.Get(fetchUrl)
	if err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	type Response struct {
		Success       bool
		Passing_times []trainmapdb.RealTrainSight
		Error         string
	}
	var response Response
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New(response.Error)
	}
	return response.Passing_times, nil
}

func renderSight(sight trainmapdb.RealTrainSight, renderFont *opentype.Font) (image.Image, error) {
	regularFace, err := opentype.NewFace(renderFont, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	bigFace, err := opentype.NewFace(renderFont, &opentype.FaceOptions{
		Size:    50,
		DPI:     75,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: 576,
			Y: 5000,
		},
	})

	const SIZE = 30
	const LINE_HEIGHT = SIZE * 1.1
	const COL_WIDTH = SIZE

	addLabel(img, bigFace, 0, 2*LINE_HEIGHT, fmt.Sprintf("%s %s", sight.Timestamp.Format("02/01 15:04"), sight.TrainSight.Feed.DisplayName))
	addLabel(img, bigFace, COL_WIDTH*4, 4*LINE_HEIGHT, fmt.Sprintf("%s %s", sight.TrainSight.RouteName, sight.TrainSight.Trip.TripShortName))

	getStopTimeRender := func(st trainmapdb.StopTime, prefix string) string {
		return fmt.Sprintf("%s %s (%s)", st.ArrivalTime.Format("15:04"), st.Stop.StopName, prefix)
	}

	addLabel(img, regularFace, 0, 5*LINE_HEIGHT, getStopTimeRender(sight.TrainSight.FirstSt, "departure"))
	addLabel(img, regularFace, 0, 6*LINE_HEIGHT, getStopTimeRender(sight.TrainSight.StBefore, "previous stop"))
	addLabel(img, regularFace, 0, 7*LINE_HEIGHT, getStopTimeRender(sight.TrainSight.StAfter, "next stop"))
	addLabel(img, regularFace, 0, 8*LINE_HEIGHT, getStopTimeRender(sight.TrainSight.LastSt, "last stop"))

	return img.SubImage(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: 576,
			Y: LINE_HEIGHT * (9),
		},
	}), nil
}

func addLabel(img *image.RGBA, face font.Face, x, y int, label string) error {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)

	return nil
}
