package generator

import (
	"context"
	"strconv"
	"time"

	"github.com/abilioesteves/metrics-generator/generator/accidenttypes"
	"github.com/abilioesteves/metrics-generator/metrics"
	"github.com/sirupsen/logrus"
)

// Tabajara generates metrics
type Tabajara struct {
	Collector *metrics.Collector
	Entropy   Entropy
	Accidents map[string]Accident
}

// NewGeneratorTabajara instantiates a new
func NewGeneratorTabajara(collector *metrics.Collector, entropy Entropy) *Tabajara {
	return &Tabajara{
		Collector: collector,
		Entropy:   entropy,
		Accidents: make(map[string]Accident),
	}
}

// Init initializes the generation of the dummy metrics
func (gen *Tabajara) Init(ctx context.Context) {
	logrus.Infof("Initializing metrics generator...")
	go func() {
		c := time.Tick(10 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				logrus.Info("Generator Tabajara stopped!")
				return
			case <-c:
				gen.FillMetrics()
			}
		}
	}()
	logrus.Infof("Metrics generator initialized!")
}

// CreateAccident creates observation accidents to an specific resource
func (gen *Tabajara) CreateAccident(accident Accident) (err error) {
	key := GetAccidentKey(accident.ResourceName, accident.Type)
	gen.Accidents[key] = accident
	logrus.Infof("Accident '%v' -> %v' installed!", key, accident)
	return
}

// DeleteAccident deletes observation accidents to an specific resource
func (gen *Tabajara) DeleteAccident(accidentType, resourceName string) (err error) {
	key := GetAccidentKey(resourceName, accidentType)
	delete(gen.Accidents, key)
	logrus.Infof("Accident '%v' removed!", key)
	return
}

// DeleteAccidents deletes all accidents
func (gen *Tabajara) DeleteAccidents() (err error) {
	gen.Accidents = make(map[string]Accident)
	logrus.Infof("All accidents removed!")
	return
}

// SetEntropy increases the number of returned time-series by n
func (gen *Tabajara) SetEntropy(e Entropy) (err error) {
	gen.Entropy = e
	logrus.Infof("Entropy '%v' set!", e)
	return
}

// FillMetrics advances the state of the registered generator metrics with configurable random values
func (gen *Tabajara) FillMetrics() {
	methods := []string{"POST", "GET", "DELETE", "PUT"}

	uri := getRandomElemNormal(gen.getUris())
	name := getRandomElemNormal(gen.getServiceNames())
	version := getRandomElemNormal(gen.getVersions())
	calls := int(gen.getValueAccident(accidenttypes.Calls, accidenttypes.DefaultNumberOfCalls, uri))

	for i := 0; i < calls; i++ {
		method := methods[randomInt(int64(hash(uri)), len(methods))]
		status, isError := getStatusWithErrorAccident(gen.getValueAccident(accidenttypes.ErrorRate, accidenttypes.DefaultErrorRate, uri))

		gen.FillRequests(uri, method, status, version, isError)
		gen.FillResponses(uri, method, status, version, isError)
	}

	gen.FillDependencies(name)
}

// FillRequests fills the RequestSecondsHistogram metric
func (gen *Tabajara) FillRequests(uri, method, status, version string, isError bool) {
	gen.Collector.RequestSecondsHistogram.WithLabelValues(
		"http",
		status,
		method,
		uri,
		strconv.FormatBool(isError),
		version,
	).Observe(gen.getValueAccident(accidenttypes.Latency, getSampleRequestTime(uri), uri))
}

// FillResponses fills the ResponseBytesCounter metric
func (gen *Tabajara) FillResponses(uri, method, status, version string, isError bool) {
	gen.Collector.ResponseBytesCounter.WithLabelValues(
		"http",
		status,
		method,
		uri,
		strconv.FormatBool(isError),
		version,
	).Add(gen.getValueAccident(accidenttypes.Latency, getSampleRequestTime(uri), uri))

}

// FillDependencies fills the DependencyUp metric
func (gen *Tabajara) FillDependencies(name string) {
	h := hash(name)
	gen.Collector.DependencyUp.WithLabelValues(
		name,
	).Set(float64(h % 2))
}

func (gen *Tabajara) getUris() []string {
	return generateItems("/resource/test-", gen.Entropy.URICount)
}

func (gen *Tabajara) getServiceNames() []string {
	return generateItems("Fake Dependency ", 10)
}

func (gen *Tabajara) getValueAccident(accidentType string, defaultValue float64, resourceName string) float64 {
	key := GetAccidentKey(resourceName, accidentType)
	if accident, ok := gen.Accidents[key]; ok {
		return accident.Value
	}
	return defaultValue
}

func (getn *Tabajara) getVersions() []string {
	return generateVersion("0.0.", 4)
}
