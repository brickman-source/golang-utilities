/*
 * Copyright (c) 2019. Brickman Source.
 */

package amap

import jsoniter "github.com/json-iterator/go"

type DirectionResponse struct {
	Status   string `json:"status" xml:"status"`
	Info     string `json:"info" xml:"info"`
	Infocode string `json:"infocode" xml:"infocode"`
	Count    string `json:"count" xml:"count"`
	Route    struct {
		Origin      string `json:"origin" xml:"origin"`
		Destination string `json:"destination" xml:"destination"`
		Paths       []struct {
			Duration     string `json:"duration" xml:"duration"`
			Strategy     string `json:"strategy" xml:"strategy"`
			Tolls        string `json:"tolls" xml:"tolls"`
			TollDistance string `json:"toll_distance" xml:"toll_distance"`
			Steps        []struct {
				Tolls           string          `json:"tolls" xml:"tolls"`
				Duration        string          `json:"duration" xml:"duration"`
				Action          jsoniter.RawMessage `json:"action" xml:"action"`
				Instruction     string          `json:"instruction" xml:"instruction"`
				Orientation     string          `json:"orientation" xml:"orientation"`
				Road            string          `json:"road,omitempty" xml:"road,omitempty"`
				Distance        string          `json:"distance" xml:"distance"`
				TollDistance    string          `json:"toll_distance" xml:"toll_distance"`
				TollRoad        jsoniter.RawMessage `json:"toll_road" xml:"toll_road"`
				Polyline        string          `json:"polyline" xml:"polyline"`
				AssistantAction jsoniter.RawMessage `json:"assistant_action" xml:"assistant_action"`
			} `json:"steps" xml:"steps"`
			Restriction   string `json:"restriction" xml:"restriction"`
			TrafficLights string `json:"traffic_lights" xml:"traffic_lights"`
			Distance      string `json:"distance" xml:"distance"`
		} `json:"paths" xml:"paths"`
	} `json:"route" xml:"route"`
}

type DistanceResponse struct {
	Status string `json:"status" xml:"status"`
	Info string `json:"info" xml:"info"`
	Infocode string `json:"infocode" xml:"infocode"`
	Results []struct {
		Distance string `json:"distance" xml:"distance"`
		Duration string `json:"duration" xml:"duration"`
		Info string `json:"info" xml:"info"`
		Code string `json:"code" xml:"code"`
		OriginID string `json:"origin_id" xml:"origin_id"`
		DestID string `json:"dest_id" xml:"dest_id"`
	} `json:"results" xml:"results"`
}
