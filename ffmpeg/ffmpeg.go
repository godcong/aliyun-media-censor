package ffmpeg

import (
	"log"
	"os"
	"os/exec"
)

// Mpeg ...
type Mpeg struct {
	OutPath         string
	MessageCallback func(map[string]interface{}) error
}

// MpegType ...
type MpegType string

// OutPath ...
const (
	OutPath MpegType = "outpath"
)

// New ...
func New(args map[MpegType]string) *Mpeg {
	return &Mpeg{
		OutPath: args[OutPath],
	}
}

// Run ...
func (m *Mpeg) Run() {

}

// Run ...
func Run(args ...string) (string, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffmpeg", args...)
	log.Println(cmd.Args)
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return string(stdout), err
	}

	return string(stdout), nil
}

// TranToMp4 ...
func TranToMp4(src string, out string) (string, error) {
	//ffmpeg -i input.mkv -acodec libfaac -vcodec libx264 out.mp4
	return Run("-i", src, "-y", "-vcodec", "libx264", "-acodec", "aac", out)

}

// CopyToMp4 ...
func CopyToMp4(src string, out string) (string, error) {
	//cmd:ffmpeg -i input.mkv -acodec copy -vcodec copy out.mp4
	return Run("-i", src, "-y", "-acodec", "copy", "-vcodec", "copy", out)
}

// TransToTS 废弃
func TransToTS(src string, out string) (string, error) {
	//cmd:ffmpeg -i INPUT.mp4 -codec copy -bsf:v h264_mp4toannexb OUTPUT.ts
	return Run("-i", src, "-y", "-codec", "copy", "-bsf:v", "h264_mp4toannexb", out)
}

// SplitToTS 废弃
func SplitToTS(src string, out string) (string, error) {
	//cmd:ffmpeg -i ./tmp/ELTbmjn2IZY6EtLFCibQPL4pIyfMvN8jQS67ntPlFaFo3NUkM3PpCFBgMivKk67W_out.mp4 -f segment -segment_time 10 -segment_format mpegts -segment_list ./split/list_file.m3u8 -c copy -bsf:v h264_mp4toannexb ./split/output_file-%d.ts
	return Run("-i", src,
		"-c", "copy", "-bsf:v", "h264_mp4toannexb",
		"-f", "segment", "-segment_time", "10",
		"-segment_format", "mpegts", "-segment_list", out+".m3u8",
		out+"-%03d.ts")

}

// Split ...
func Split(src string, out string, media, m3u8 string) (string, error) {
	//ffmpeg -i input.mp4 -c copy -bsf:v h264_mp4toannexb -hls_time 10 -hls_key_info_file key_info playlist.m3u8
	return Run("-i", src,
		"-y", "-c:v", "libx264", "-c:a", "aac",
		//"-c:v", "copy",
		"-bsf:v", "h264_mp4toannexb",
		//"-f", "segment", "-segment_time", "10",
		//"-hls_flags", "delete_segments",
		"-f", "hls", "-hls_time", "10",
		"-hls_playlist_type", "vod",
		//"-segment_format", "mpegts",
		"-hls_segment_filename", out+"/"+media+"-%03d.ts",
		out+"/"+m3u8)
}

// QuickSplit ...
func QuickSplit(src string, out string, media, m3u8 string) (string, error) {
	//ffmpeg -i input.mp4 -c copy -bsf:v h264_mp4toannexb -hls_time 10 -hls_key_info_file key_info playlist.m3u8
	return Run("-i", src,
		"-y", "-c:v", "copy", "-c:a", "copy",
		"-bsf:v", "h264_mp4toannexb",
		//"-f", "segment", "-segment_time", "10",
		//"-hls_flags", "delete_segments",
		"-f", "hls", "-hls_time", "10",
		"-hls_playlist_type", "vod",
		//"-segment_format", "mpegts",
		"-hls_segment_filename", out+"/"+media+"-%03d.ts",
		out+"/"+m3u8)
}

// SplitWithKey ...
func SplitWithKey(src string, out string, key string, media, m3u8 string) (string, error) {
	//ffmpeg -i input.mp4 -c copy -bsf:v h264_mp4toannexb -hls_time 10 -hls_key_info_file key_info playlist.m3u8
	return Run("-i", src,
		"-y", "-c:v", "libx264", "-c:a", "aac",
		//"-c:v", "copy",
		"-bsf:v", "h264_mp4toannexb",
		//"-f", "segment", "-segment_time", "10",
		//"-hls_flags", "delete_segments",
		"-f", "hls", "-hls_time", "10",
		"-hls_playlist_type", "vod",
		//"-segment_format", "mpegts",
		"-hls_segment_filename", out+"/"+media+"-%03d.ts",
		"-hls_key_info_file", out+"/"+key,
		out+"/"+m3u8)
}

// QuickSplitWithKey ...
func QuickSplitWithKey(src string, out string, key string, media, m3u8 string) (string, error) {
	//ffmpeg -i input.mp4 -c copy -bsf:v h264_mp4toannexb -hls_time 10 -hls_key_info_file key_info playlist.m3u8
	return Run("-i", src,
		"-y", "-c:v", "copy", "-c:a", "copy",
		"-bsf:v", "h264_mp4toannexb",
		//"-f", "segment", "-segment_time", "10",
		//"-hls_flags", "delete_segments",
		"-f", "hls", "-hls_time", "10",
		"-hls_playlist_type", "vod",
		//"-segment_format", "mpegts",
		"-hls_segment_filename", out+"/"+media+"-%03d.ts",
		"-hls_key_info_file", out+"/"+key,
		out+"/"+m3u8)
}

// SplitWithKey1 废弃
func SplitWithKey1(src string, out string, keyPath string) (string, error) {
	//ffmpeg -i input.mp4 -c copy -bsf:v h264_mp4toannexb -hls_time 10 -hls_key_info_file key_info playlist.m3u8
	return Run("-i", src,
		"-c", "copy", "-bsf:v", "h264_mp4toannexb",
		"-hls_time", "10", "-hls_key_info_file", keyPath,
		out+".m3u8")
}

// TransferSplitWith 切分:320*240,帧率:24
func TransferSplitWith(path string, src string, size string, time string) (string, error) {
	os.MkdirAll("transferred/"+src, os.ModePerm)
	return Run("-i", path+src,
		"-y", "-c:v", "libx264", "-s", size,
		"-r", "24",
		"-bsf:v", "h264_mp4toannexb",
		"-f", "hls", "-hls_time", time,
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", "transferred/"+src+"/output-%03d.ts",
		"transferred/"+src+"/output.m3u8")
}

// TransferSplit ...
func TransferSplit(path string, src string) (string, error) {
	return TransferSplitWith(path, src, "320*240", "1800")
}

// TransferJPG ...
func TransferJPG(src string, dest string) (string, error) {
	_ = os.MkdirAll(dest, os.ModePerm)
	return Run("-i", src,
		"-f", "image2",
		"-vf", "fps=fps=1/15",
		dest+"/censor-%03d.jpg",
	)
}
