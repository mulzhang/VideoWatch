package video

import (
	"context"
	"fmt"
	"gopkg.in/vansante/go-ffprobe.v2"
	"time"
)

func GetMediaDurationByUrl(url string) float64 {
	ctx, cancelFn := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFn()
	data, err := ffprobe.ProbeURL(ctx, url)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	if data == nil {
		return 0
	}

	return data.Format.Duration().Seconds()
}
