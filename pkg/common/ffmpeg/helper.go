package ffmpeg

import (
	"log/slog"
	"os/exec"
)

// GetCover 通过ffmpeg获取封面
func GetCover(filename string, outPath string) {
	cmd := exec.Command("./utils/ffmpeg.exe", "-i", filename, "-ss", "1", "-f", "image2", "-y", outPath)
	err := cmd.Run()
	if err != nil {
		slog.Error("get cover by ffmpeg failed: ", err)
	}
}
