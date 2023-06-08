package sse

import (
	"github.com/go-redis/redis/v8"
)

//var (
//	rdb    *redis.Client
//	pubSub *redis.PubSub
//)

func InitSSE() (*redis.Client, *redis.PubSub) {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize PubSub
	pubSub := rdb.Subscribe(rdb.Context(), "broadcast")
	return rdb, pubSub
	//// Handle SSE requests
	//http.HandleFunc("/events", Handler)
	//
	//// Start the HTTP server
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

//func HandleSSE(app *config.Application) func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		// Set response headers for SSE
//		w.Header().Set("Content-Type", "text/event-stream")
//		w.Header().Set("Cache-Control", "no-cache")
//		w.Header().Set("Connection", "keep-alive")
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//
//		// Create a new PubSub channel
//		//channel := make(chan *redis.Message)
//		//channel := pubSub.Channel()
//
//		// Add the channel to the PubSub pubSub
//		// We get the user from the request, we then subscribe him to his specific Redis channel
//		user := context_config.ContextGetAuthenticatedUser(r)
//		app.Logger.Infof("%s joined", user.Username)
//		userChannelName := fmt.Sprintf("user#%d", user.ID)
//		err := app.PubSub.Subscribe(app.Rdb.Context(), userChannelName)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		defer app.PubSub.Close()
//
//		channel := app.PubSub.Channel()
//
//		for msg := range channel {
//			app.Logger.Infof(msg.Channel, msg.Payload)
//			fmt.Fprintf(w, "data: %s\n\n", strings.TrimSuffix(msg.Payload, "\n"))
//			w.(http.Flusher).Flush()
//		}
//
//	}
//}
