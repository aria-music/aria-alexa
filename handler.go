package main

import (
	"context"
	"log"

	"github.com/aria-music/aria-alexa/alexa"
)

type requestHandler func(ctx context.Context, request *alexa.Request) *alexa.Response

var requestHandlers = map[string]requestHandler{
	"IntentRequest":    handleIntentRequest,
	"LaunchRequest":    handleLaunchRequest,
	"PlaybackStarted":  getLogHandler("PlaybackStarted"),
	"PlaybackFinished": getLogHandler("PlaybackFinished"),
	"PlaybackStopped":  getLogHandler("PlaybackStopped"),
}

func getLogHandler(msg string) requestHandler {
	return func(_ context.Context, _ *alexa.Request) *alexa.Response {
		log.Println(msg)
		return nil
	}
}

func handleRequest(ctx context.Context, request *alexa.Request) *alexa.Response {
	if h, ok := requestHandlers[request.Type()]; ok {
		log.Printf("handling %s", request.Type())
		return h(ctx, request)
	}

	log.Printf("no handlers were found: %s", request.Type())
	return nil
}

// handleIntent
func handleIntentRequest(ctx context.Context, request *alexa.Request) *alexa.Response {
	return handleIntent(ctx, request)
}

func handleLaunchRequest(ctx context.Context, request *alexa.Request) *alexa.Response {
	return alexa.NewAudioPlayerPlayResponse(audioEndpoint)
}
