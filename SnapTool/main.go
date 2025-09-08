package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	fps        = 25
	concurrent = 5 // 同时处理的并发数
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: SnapTool <video_file> [output_directory]")
		fmt.Println("Example: SnapTool video.mp4 ./screenshots")
		os.Exit(1)
	}

	videoFile := os.Args[1]
	outputDir := "./screenshots"
	if len(os.Args) > 2 {
		outputDir = os.Args[2]
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// 检查视频文件是否存在
	if _, err := os.Stat(videoFile); os.IsNotExist(err) {
		log.Fatal("Video file does not exist:", videoFile)
	}

	// 检查ffmpeg和ffprobe是否可用
	if !checkFFmpeg() {
		log.Fatal("FFmpeg or ffprobe not found. Please install FFmpeg and make sure it's in your PATH.")
	}

	fmt.Printf("Starting screenshot process for %s at %d fps\n", videoFile, fps)
	fmt.Printf("Output directory: %s\n", outputDir)

	// 截图过程
	if err := takeScreenshots(videoFile, outputDir); err != nil {
		log.Fatal("Error taking screenshots:", err)
	}

	fmt.Println("Screenshot process completed successfully!")
}

func checkFFmpeg() bool {
	_, err1 := exec.LookPath("ffmpeg")
	_, err2 := exec.LookPath("ffprobe")
	return err1 == nil && err2 == nil
}

func takeScreenshots(videoFile, outputDir string) error {
	// 获取视频时长
	duration, err := getVideoDuration(videoFile)
	if err != nil {
		return fmt.Errorf("failed to get video duration: %v", err)
	}

	fmt.Printf("Video duration: %.2f seconds\n", duration)

	// 计算帧间隔时间(以秒为单位)
	frameInterval := 1.0 / float64(fps)
	totalFrames := int(duration * float64(fps))

	fmt.Printf("Total frames to capture: %d\n", totalFrames)

	// 使用goroutine和channel来并发处理截图
	jobs := make(chan job, totalFrames)
	results := make(chan result, totalFrames)

	// 启动worker goroutines
	var wg sync.WaitGroup
	for w := 0; w < concurrent; w++ {
		wg.Add(1)
		go worker(videoFile, outputDir, jobs, results, &wg)
	}

	// 发送任务到jobs channel
	go func() {
		defer close(jobs)
		for i := 0; i < totalFrames; i++ {
			timestamp := float64(i) * frameInterval
			filename := fmt.Sprintf("frame_%06d.jpg", i)
			jobs <- job{i, timestamp, filename}
		}
	}()

	// 关闭results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果并显示进度
	completed := 0
	for res := range results {
		if res.err != nil {
			return fmt.Errorf("failed to capture frame %d: %v", res.frameIndex, res.err)
		}
		completed++
		
		// 显示进度
		if completed%fps == 0 || completed == totalFrames {
			progress := float64(completed) / float64(totalFrames) * 100
			fmt.Printf("Progress: %.1f%% (%d/%d frames)\n", progress, completed, totalFrames)
		}
	}

	return nil
}

type job struct {
	frameIndex int
	timestamp  float64
	filename   string
}

type result struct {
	frameIndex int
	err        error
}

func worker(videoFile, outputDir string, jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		outputPath := filepath.Join(outputDir, j.filename)
		err := captureFrame(videoFile, j.timestamp, outputPath)
		results <- result{j.frameIndex, err}
	}
}

func getVideoDuration(videoFile string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-show_entries", "format=duration", "-of", "csv=p=0", videoFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get video duration: %v", err)
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %v", err)
	}

	return duration, nil
}

func captureFrame(videoFile string, timestamp float64, outputPath string) error {
	// 使用ffmpeg截图指定时间点的帧
	timestampStr := fmt.Sprintf("%.6f", timestamp)
	
	// 使用更快的截图方法: 精确定位到时间戳并只输出一帧
	cmd := exec.Command("ffmpeg", "-ss", timestampStr, "-i", videoFile, 
		"-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2", // 确保尺寸为偶数
		"-vframes", "1", "-q:v", "2", "-y", outputPath)
	
	// 设置CPU亲和性以提高性能
	if runtime.GOOS != "windows" {
		cmd.Env = append(os.Environ(), "OMP_NUM_THREADS=1")
	}
	
	// 忽略stdout输出，但保留stderr以便调试
	cmd.Stdout = nil
	
	return cmd.Run()
}