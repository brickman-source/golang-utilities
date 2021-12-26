/*
 * Copyright (c) 2019. Brickman.
 */

package amap

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/json"
	"math"
	"time"

	"github.com/brickman-source/golang-utilities/converter"
	"github.com/brickman-source/golang-utilities/errors"
	httpUtil "github.com/brickman-source/golang-utilities/http" //
)

/*
下方策略仅返回一条路径规划结果

0，速度优先，不考虑当时路况，此路线不一定距离最短
1，费用优先，不走收费路段，且耗时最少的路线
2，距离优先，不考虑路况，仅走距离最短的路线，但是可能存在穿越小路/小区的情况
3，速度优先，不走快速路，例如京通快速路（因为策略迭代，建议使用13）
4，躲避拥堵，但是可能会存在绕路的情况，耗时可能较长
5，多策略（同时使用速度优先、费用优先、距离优先三个策略计算路径）。
	其中必须说明，就算使用三个策略算路，会根据路况不固定的返回一~三条路径规划信息。
6，速度优先，不走高速，但是不排除走其余收费路段
7，费用优先，不走高速且避免所有收费路段
8，躲避拥堵和收费，可能存在走高速的情况，并且考虑路况不走拥堵路线，但有可能存在绕路和时间较长
9，躲避拥堵和收费，不走高速
*/
func DrivingDistance(long1, lat1, long2, lat2 float64, strategy int, key string) (distance float64, duration time.Duration, tollsFee float64, err error) {
	directionURL := fmt.Sprintf(
		"http://restapi.amap.com/v3/direction/driving?origin=%0.6f,%0.6f&destination=%0.6f,%0.6f&extensions=base&strategy=%d&ferry=%d&nosteps=%d&key=%s",
		long1, lat1,
		long2, lat2,
		strategy,
		0, //0：使用渡轮(默认) 1：不使用渡轮
		1, //0：steps字段内容正常返回；1：steps字段内容为空；
		key,
	)

	directionResp := &DirectionResponse{}

	//amapApiStart := time.Now()

	err = httpUtil.GetJSON(directionURL, directionResp)

	//callAmapElapsed := time.Since(amapApiStart)

	if err != nil {
		return -1, time.Duration(0), 0, err
	}

	if directionResp.Status != "1" {
		return -1, time.Duration(0), 0, errors.New(directionResp.Status)
	}
	if len(directionResp.Route.Paths) < 1 {
		return -1, time.Duration(0), 0, errors.New("no route")
	}

	shortestDistance := math.MaxFloat64
	routeDuration := 0.0
	amapRouteIsValid := false
	amapTollsFee := 0.0

	for _, routePath := range directionResp.Route.Paths {
		routePathDistance := converter.AsFloat64(routePath.Distance, 0)
		if routePathDistance > 0 && routePathDistance < shortestDistance {
			amapTollsFee = converter.AsFloat64(routePath.Tolls, 0)
			routeDuration = converter.AsFloat64(routePath.Duration, 0)
			shortestDistance = routePathDistance
		}
		amapRouteIsValid = true
	}

	if !amapRouteIsValid {
		return -1, time.Duration(0), 0, errors.New("cannot calc fee")
	}

	return shortestDistance, time.Duration(routeDuration), amapTollsFee, nil
}

// 计算距离
// 返回中，distance单位是米
func Distance(long1, lat1, long2, lat2 float64, key string) (distance float64, duration time.Duration, err error) {
	directionURL := fmt.Sprintf(
		`http://restapi.amap.com/v3/distance?origins=%0.6f,%0.6f&destination=%0.6f,%0.6f&output=json&key=%s&type=1`,
		long1, lat1,
		long2, lat2,
		key,
	)

	distResp := &DistanceResponse{}

	err = httpUtil.GetJSON(directionURL, distResp)
	if err != nil {
		return -1, time.Duration(0), err
	}

	println(string(json.ShouldMarshal(distResp)))
	if distResp.Status != "1" {
		return -1, time.Duration(0), errors.New(distResp.Status)
	}

	if len(distResp.Results) < 1 {
		return -1, time.Duration(0), errors.New(distResp.Info)
	}

	return converter.AsFloat64(distResp.Results[0].Distance, -1),
		converter.AsDuration(distResp.Results[0].Duration, 0),
		nil
}
