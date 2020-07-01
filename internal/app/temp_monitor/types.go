package temp_monitor

type SourceWeatherCloud string

const (
	AmbientWeather SourceWeatherCloud = "AmbientWeather"
)

type EnvConfig struct {
	SourceWeatherCloud string `required:"true" split_words:"true"`
	FunctionID         string `required:"true" split_words:"true"`
	SNSTopicARN        string `required:"true" split_words:"true"`
}
