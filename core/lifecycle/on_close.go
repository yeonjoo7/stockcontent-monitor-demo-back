package lifecycle

import "stockcontent-monitor-demo-back/core/app"

func ProvidesOnClose() app.OnClose {
	return func(err error) {

	}
}
