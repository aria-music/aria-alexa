package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aria-music/aria-alexa/alexa"
	"github.com/aria-music/aria-alexa/aria"
)

type intentHandler func(ctx context.Context, intent *alexa.Intent) *alexa.Response

var intentHandlers = map[string]intentHandler{
	"AMAZON.PauseIntent":      handleAmazonPauseIntent,
	"AMAZON.ResumeIntent":     handleAmazonResumeIntent,
	"AMAZON.NextIntent":       handleAmazonNextIntent,
	"AMAZON.ShuffleOnIntent":  handleAmazonShuffleIntent,
	"AMAZON.ShuffleOffIntent": handleAmazonShuffleIntent,
	"AMAZON.RepeatIntent":     handleAmazonRepeatIntent,
}

// parseIntent parses alexa.Intent from *alexa.Request.
// returns *alexa.Intent: this won't be null as declared in alexa.IntentRequest
func parseIntent(r *alexa.Request) (intent *alexa.Intent) {
	ir := new(alexa.IntentRequest)
	if err := json.Unmarshal(r.Request, ir); err != nil {
		log.Printf("failed to unmarshal IntentRequest: %v (body: %s)", err, r.Request)
		return
	}
	return &ir.Intent
}

func handleIntent(ctx context.Context, request *alexa.Request) *alexa.Response {
	intent := parseIntent(request)
	if intent == nil {
		log.Printf("failed in parseIntent.")
		return nil
	}

	if h, ok := intentHandlers[intent.Name]; ok {
		log.Printf("handling %s", intent.Name)
		return h(ctx, intent)
	}

	log.Printf("no intent handlers were found: %s", intent.Name)
	return nil
}

func handleAmazonPauseIntent(ctx context.Context, intent *alexa.Intent) *alexa.Response {
	if err := sendOP(ctx, "pause"); err != nil {
		log.Printf("failed to pause: restapi: %v", err)
		return alexa.ServerErrorResponse
	}
	return alexa.EmptyResponse
}

func handleAmazonResumeIntent(ctx context.Context, intent *alexa.Intent) *alexa.Response {
	if err := sendOP(ctx, "resume"); err != nil {
		log.Printf("failed to resume: restapi: %v", err)
		return alexa.ServerErrorResponse
	}
	return alexa.EmptyResponse
}

func handleAmazonNextIntent(ctx context.Context, intent *alexa.Intent) *alexa.Response {
	if err := sendOP(ctx, "skip"); err != nil {
		log.Printf("failed to skip: restapi: %v", err)
		return alexa.ServerErrorResponse
	}
	end := true
	return alexa.NewTextOutputSpeechResponse("スキップしました", &end)
}

func handleAmazonShuffleIntent(ctx context.Context, intent *alexa.Intent) *alexa.Response {
	if err := sendOP(ctx, "shuffle"); err != nil {
		log.Printf("failed to shuffle: restapi: %v", err)
		return alexa.ServerErrorResponse
	}
	end := true
	return alexa.NewTextOutputSpeechResponse("シャッフルしました", &end)
}

func handleAmazonRepeatIntent(ctx context.Context, intent *alexa.Intent) *alexa.Response {
	st, err := sendRequest(ctx, "state", nil)
	if err != nil {
		log.Printf("failed to repeat: failed to get state: %v", err)
		return alexa.ServerErrorResponse
	}

	sdata, ok := st.Data.(*aria.StateData)
	if !ok || sdata.Entry.URI == "" {
		log.Printf("failed to repeat: invalid response for op state (%v)", sdata)
		return alexa.ServerErrorResponse
	}

	if _, err := sendRequest(ctx, "repeat", &aria.RepeatRequestData{
		URI: sdata.Entry.URI,
	}); err != nil {
		log.Printf("failed to repeat: restapi: %v", err)
		return alexa.ServerErrorResponse
	}
	end := true
	return alexa.NewTextOutputSpeechResponse("リピートしました", &end)
}
