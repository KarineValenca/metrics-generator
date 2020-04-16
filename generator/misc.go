package generator

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
)

func generateSample(median float64, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + median
}

func getRandomElem(items []string) string {
	return items[rand.Intn(len(items))]
}

func getRandomElemNormal(items []string) string {
	return items[randomRangeNormal(0, len(items))] //randomRangeNormal guarantees len(items) will be unreachable
}

func generateItems(prefix string, qtty int) []string {
	result := []string{}
	for i := 1; i <= qtty; i++ {
		result = append(result, fmt.Sprintf("%s%04d", prefix, i))
	}
	return result
}

func generateVersion(prefix string, qtty int) []string {
	result := []string{}
	for i := 1; i <= qtty; i++ {
		result = append(result, fmt.Sprintf("%s%d", prefix, i))
	}
	return result
}

func randomRangeNormal(min int, max int) int {
	median := float64((max - min) / 2)
	v := generateSample(median+0.5, median/3)
	mi := float64(min)
	ma := float64(max) - 0.1 // // max will be unreachable
	return min + int(math.Max(math.Min(float64(v), ma), mi))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func randomInt(seed int64, max int) int {
	s1 := rand.NewSource(seed)
	return rand.New(s1).Intn(max)
}

func getSampleRequestTime(uri string) (requestTime float64) {
	median := (hash(uri)%23 + 1) * 100                                        // 100 - 2300 indexed by uri
	requestTime = generateSample(float64(median), float64(median)/5) / 1000.0 // sample normal values around the defined median
	requestTime = math.Max(requestTime, 0.01)                                 // minimum request time defined as 0.01
	return
}

func getStatusWithErrorAccident(accident float64) (string, bool) {
	sr := rand.Intn(100)
	abs := int(accident) * 100
	if sr <= abs {
		return "5xx", true
	} else if sr > abs && sr < 1-abs {
		return "4xx", true
	} else {
		return "2xx", false
	}
}
