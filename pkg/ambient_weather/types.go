// Implementation of iface.WeatherCloud that accesses ambientweather.net's API,
// their api is documented at https://ambientweather.docs.apiary.io/
package ambient_weather

// iso8601Format date format for parsing ISO8601 dates
const iso8601Format = "2006-01-02T15:04:05Z0700"

const externalTempKey = "tempf"
const externalTempFeelsKey = "feelsLike"
const internalTempKey = "tempinf"
const internalTempFeelsKey = "feelsLikein"
const dateKey = "dateutc"

// Cloud represents an Ambient Weather cloud backend
type Cloud struct {
	deviceID       string
	apiKey         string
	applicationKey string
}

// NewCloud instantiates an Ambient Weather Cloud
func NewCloud(deviceID string, apiKey string, applicationKey string) *Cloud {
	return &Cloud{
		deviceID:       deviceID,
		apiKey:         apiKey,
		applicationKey: applicationKey,
	}
}
