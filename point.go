package model

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type Point struct {
	Point *wkb.Point
}

func (p Point) LatLng() *LatLng {
	if p.Point == nil {
		return nil
	}

	coords := p.Point.Coords()

	return &LatLng{
		Lat: coords.Y(),
		Lng: coords.X(),
	}
}

func (p Point) GormDataType() string {
	return "geometry"
}

// Set point as Latitude and Longitude coordinates.
func (p *Point) SetCoordinates(lat, lng float64) {
	p.Point = &wkb.Point{
		Point: geom.NewPoint(geom.XY).MustSetCoords([]float64{lng, lat}).SetSRID(4326),
	}
}

func (p Point) Value() (driver.Value, error) {
	if p.Point == nil {
		return nil, nil
	}

	value, err := p.Point.Value()
	if err != nil {
		return nil, err
	}

	buf, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("did not convert value: expected []byte, but was %T", value)
	}

	mysqlEncoding := make([]byte, 4)
	binary.LittleEndian.PutUint32(mysqlEncoding, 4326)
	mysqlEncoding = append(mysqlEncoding, buf...)

	return mysqlEncoding, err
}

func (p *Point) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	mysqlEncoding, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", src)
	}

	var srid uint32 = binary.LittleEndian.Uint32(mysqlEncoding[0:4])

	var pnt wkb.Point

	if err := pnt.Scan(mysqlEncoding[4:]); err != nil {
		return err
	}

	pnt.SetSRID(int(srid))

	*p = Point{
		Point: &pnt,
	}

	return nil
}

type LatLng struct {
	Lat float64
	Lng float64
}
