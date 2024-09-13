/*
计算两个GPS坐标这件的距离
*/
package distance

import "math"

const (
	earthRadius = 6378137              //赤道半径
	radian      = 0.017453292519943295 //math.Pi / 180.0
)

func rad(value float64) float64 {
	return value * radian
}

func Distance(longitude, latitude, otherLongitude, otherLatitude float64) float64 {
	radLatitude := rad(latitude)
	radOtherLatitude := rad(otherLatitude)

	a := radLatitude - radOtherLatitude
	b := rad(longitude) - rad(otherLongitude)

	return earthRadius * 2 * math.Asin(
		math.Sqrt(
			math.Pow(
				math.Sin(a/2), 2)+
				math.Cos(radLatitude)*
					math.Cos(radOtherLatitude)*
					math.Pow(math.Sin(b/2), 2)))
}
