package ffmpeg

import "testing"

func TestTransferJPG(t *testing.T) {
	s, e := TransferJPG("download/[Thz.la]ipx-091.mp4", "transfer/out1")
	t.Log(s, e)
}
