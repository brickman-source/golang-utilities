package geojson

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/map/helper"
	"github.com/paulmach/go.geojson"
	"io/ioutil"
	"math"
	"path"
	"sync"
)

const (
	X_PI   = math.Pi * 3000.0 / 180.0
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

var (
	geoFeatures     map[int]*geojson.FeatureCollection
	geoFeaturesLock sync.RWMutex
)

func Preload(jsonDir string) {
	gfs := make(map[int]*geojson.FeatureCollection, 0)

	content, err := ioutil.ReadFile(path.Join(jsonDir, "100000_full.json"))

	if err != nil {
		panic(err)
	}
	fc1, err := geojson.UnmarshalFeatureCollection(content)
	if err != nil || fc1 == nil {
		panic(err)
	}
	gfs[100000] = fc1

	for _, feature := range fc1.Features {
		childAdCode := helper.GetValueAsInt(feature.Properties, "adcode", 0)
		level := helper.GetValueAsString(feature.Properties, "level", "")

		if childAdCode == 100000 || childAdCode == 0 || level == "" {
			continue
		}
		// 直辖市: 北京、天津、上海、重庆
		// 香港、澳门
		if childAdCode == 110000 || childAdCode == 120000 || childAdCode == 310000 || childAdCode == 500000 || childAdCode == 810000 || childAdCode == 820000 {
			continue
		}
		if feature.Geometry == nil {
			continue
		}

		childContent, err := ioutil.ReadFile(path.Join(jsonDir, fmt.Sprintf("%d_full.json", childAdCode)))
		if err != nil {
			panic(err)
		}
		childFc1, err := geojson.UnmarshalFeatureCollection(childContent)
		if err != nil || childFc1 == nil {
			panic(err)
		}
		gfs[childAdCode] = childFc1
	}

	geoFeaturesLock.Lock()
	defer geoFeaturesLock.Unlock()

	geoFeatures = gfs
}

func FindGeoFeature(lng, lat float64, adCode int) *geojson.Feature {
	var fc1 *geojson.FeatureCollection
	var ok bool

	if geoFeatures == nil {
		return nil
	}
	if adCode == 0 {
		adCode = 100000
	}
	geoFeaturesLock.RLock()
	if fc1, ok = geoFeatures[adCode]; !ok || fc1 == nil {
		geoFeaturesLock.RUnlock()
		return nil
	}
	geoFeaturesLock.RUnlock()

	for _, feature := range fc1.Features {
		childAdCode := helper.GetValueAsInt(feature.Properties, "adcode", 0)
		level := helper.GetValueAsString(feature.Properties, "level", "")

		if childAdCode == 100000 || childAdCode == 0 || level == "" {
			continue
		}
		if feature.Geometry == nil {
			continue
		}

		findInPolygon := func(polygon [][][]float64) *geojson.Feature {
			for _, points := range polygon {
				if containsPoint(points, []float64{lng, lat}) {
					if level == "province" {
						// 直辖市: 北京、天津、上海、重庆
						// 行政区: 香港、澳门
						if childAdCode == 110000 ||
							childAdCode == 120000 ||
							childAdCode == 310000 ||
							childAdCode == 500000 ||
							childAdCode == 810000 ||
							childAdCode == 820000 {
							return feature
						}
						return FindGeoFeature(lng, lat, childAdCode)
					} else {
						return feature
					}
				}
			}
			return nil
		}
		if feature.Geometry.Type == geojson.GeometryPolygon {
			f := findInPolygon(feature.Geometry.Polygon)
			if f != nil {
				return f
			}
		} else if feature.Geometry.Type == geojson.GeometryMultiPolygon {
			for _, polygon := range feature.Geometry.MultiPolygon {
				f := findInPolygon(polygon)
				if f != nil {
					return f
				}
			}

		}
	}
	return nil
}

func FindGeoFeatureByAddCode(adCode int) *geojson.Feature {
	var provinceFeature *geojson.Feature
	var cityFeature *geojson.Feature

	if geoFeatures == nil {
		return nil
	}
	if adCode == 0 {
		return nil
	}
	geoFeaturesLock.RLock()
	defer geoFeaturesLock.RUnlock()

	provinceAdCode := adCode / 1000 * 1000

	if fcProvinces, ok := geoFeatures[100000]; ok && fcProvinces != nil {
		for _, pf := range fcProvinces.Features {
			adCodeT := helper.GetValueAsInt(pf.Properties, "adcode", 0)
			if provinceAdCode == adCodeT {
				provinceFeature = pf
				break
			}
		}
	}
	if provinceFeature == nil {
		return nil
	}

	if adCode%1000 == 0 {
		return provinceFeature
	}

	cityAdCode := adCode / 100 * 100
	if fcCities, ok := geoFeatures[provinceAdCode]; ok && fcCities != nil {
		for _, cf := range fcCities.Features {
			adCodeT := helper.GetValueAsInt(cf.Properties, "adcode", 0)
			if cityAdCode == adCodeT {
				cityFeature = cf
				break
			}
		}
	}

	if cityFeature != nil {
		return cityFeature
	}

	return provinceFeature
}

func ConvertCoordinates(geoStr string, from, to string) string {
	fc, err := geojson.UnmarshalFeatureCollection([]byte(geoStr))
	if err != nil {
		return geoStr
	}
	for _, feature := range fc.Features {
		if feature.Geometry == nil {
			continue
		}
		if feature.Geometry.Type == "Polygon" {
			if len(feature.Geometry.Polygon) > 0 {
				polygon := make([][][]float64, len(feature.Geometry.Polygon))
				for pointsIndex, points := range feature.Geometry.Polygon {
					if len(points) > 0 {
						polygon[pointsIndex] = make([][]float64, len(points))
						for pointIndex, point := range points {
							lng, lat := convertGeoPoint(point[0], point[1], from, to)
							polygon[pointsIndex][pointIndex] = []float64{
								lng, lat,
							}
						}
					}
				}
				feature.Geometry.Polygon = polygon
			}
		} else if feature.Geometry.Type == "MultiPolygon" {
			if len(feature.Geometry.MultiPolygon) > 0 {
				multiPolygon := make([][][][]float64, len(feature.Geometry.MultiPolygon))
				for polygonIndex, polygon := range feature.Geometry.MultiPolygon {
					nPolygon := make([][][]float64, len(polygon))
					for pointsIndex, points := range polygon {
						if len(points) > 0 {
							nPolygon[pointsIndex] = make([][]float64, len(points))
							for pointIndex, point := range points {
								lng, lat := convertGeoPoint(point[0], point[1], from, to)
								nPolygon[pointsIndex][pointIndex] = []float64{
									lng, lat,
								}
							}
						}
					}
					multiPolygon[polygonIndex] = nPolygon
				}
				feature.Geometry.MultiPolygon = multiPolygon
			}
		}
	}
	converted, err := fc.MarshalJSON()
	if err != nil {
		return geoStr
	}
	return string(converted)
}
func containsPoint(polygon [][]float64, testp []float64) bool {
	minX := polygon[0][0]
	maxX := polygon[0][0]
	minY := polygon[0][1]
	maxY := polygon[0][1]

	for _, p := range polygon {
		minX = min(p[0], minX)
		maxX = max(p[0], maxX)
		minY = min(p[1], minY)
		maxY = max(p[1], maxY)
	}

	if testp[0] < minX || testp[0] > maxX || testp[1] < minY || testp[1] > maxY {
		return false
	}

	inside := false
	j := len(polygon) - 1
	for i := 0; i < len(polygon); i++ {
		if (polygon[i][1] > testp[1]) != (polygon[j][1] > testp[1]) && testp[0] < (polygon[j][0]-polygon[i][0])*(testp[1]-polygon[i][1])/(polygon[j][1]-polygon[i][1])+polygon[i][0] {
			inside = !inside
		}
		j = i
	}

	return inside
}

func min(n1, n2 float64) float64 {
	if n1 <= n2 {
		return n1
	}
	return n2
}
func max(n1, n2 float64) float64 {
	if n1 >= n2 {
		return n1
	}
	return n2
}
func convertGeoPoint(lng, lat float64, from, to string) (float64, float64) {
	if from == "WGS84" && to == "GCJ02" {
		mgLng, mgLat := delta(lng, lat)
		return mgLng, mgLat
	}
	if from == "GCJ02" && to == "WGS84" {
		mgLon, mgLat := delta(lng, lat)
		return lng*2 - mgLon, lat*2 - mgLat
	}
	return lng, lat
}
func delta(lon, lat float64) (float64, float64) {
	dlat, dlon := transform(lon-105.0, lat-35.0)
	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)

	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)

	mgLat := lat + dlat
	mgLon := lon + dlon

	return mgLon, mgLat
}
func transform(lon, lat float64) (x, y float64) {
	var lonlat = lon * lat
	var absX = math.Sqrt(math.Abs(lon))
	var lonPi, latPi = lon * math.Pi, lat * math.Pi
	var d = 20.0*math.Sin(6.0*lonPi) + 20.0*math.Sin(2.0*lonPi)
	x, y = d, d
	x += 20.0*math.Sin(latPi) + 40.0*math.Sin(latPi/3.0)
	y += 20.0*math.Sin(lonPi) + 40.0*math.Sin(lonPi/3.0)
	x += 160.0*math.Sin(latPi/12.0) + 320*math.Sin(latPi/30.0)
	y += 150.0*math.Sin(lonPi/12.0) + 300.0*math.Sin(lonPi/30.0)
	x *= 2.0 / 3.0
	y *= 2.0 / 3.0
	x += -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lonlat + 0.2*absX
	y += 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lonlat + 0.1*absX
	return
}
