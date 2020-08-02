package alexa

import "log"

type Response struct {
	Version  string         `json:"version"`
	Response ResponseObject `json:"response"`
}

type ResponseObject struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	ShouldEndSession *bool         `json:"shouldEndSession,omitempty"`
	Directives       []interface{} `json:"directives,omitempty"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type AudioPlayerPlayDirective struct {
	Type         string `json:"type"`
	PlayBehavior string `json:"playBehavior"`
	AudioItem    struct {
		Stream struct {
			URL                  string `json:"url"`
			Token                string `json:"token"`
			OffsetInMilliseconds int    `json:"offsetInMilliSeconds"`
		} `json:"stream"`
	} `json:"audioItem"`
}

var EmptyResponse = newEmptyResponse()
var ServerErrorResponse = newServerErrorResponse()

func newEmptyResponse() *Response {
	res := newBaseResponse()
	end := true
	res.Response.ShouldEndSession = &end
	return res
}

func newServerErrorResponse() *Response {
	end := true
	return NewTextOutputSpeechResponse("サーバーに接続できません", &end)
}

func NewTextOutputSpeechResponse(text string, end *bool) *Response {
	if text == "" {
		log.Printf("Empty text is not allowed.")
		return nil
	}

	res := newBaseResponse()
	res.Response.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
	res.Response.ShouldEndSession = end
	return res
}

func NewAudioPlayerPlayResponse(audioEndpoint string) *Response {
	end := true
	res := NewTextOutputSpeechResponse("起動しました。 音声コントロールを利用できます。", &end)
	audio := &AudioPlayerPlayDirective{
		Type:         "AudioPlayer.Play",
		PlayBehavior: "REPLACE_ALL",
	}
	audio.AudioItem.Stream.URL = audioEndpoint
	audio.AudioItem.Stream.OffsetInMilliseconds = 0
	audio.AudioItem.Stream.Token = "aria"
	res.Response.Directives = append(res.Response.Directives, &audio)

	return res
}

func newBaseResponse() *Response {
	return &Response{
		Version: "1.0",
	}
}
