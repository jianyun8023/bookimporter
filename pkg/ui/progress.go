package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// ProgressTracker 进度跟踪器
type ProgressTracker struct {
	total       int
	current     int
	startTime   time.Time
	lastMessage string
	width       int
}

// NewProgressTracker 创建新的进度跟踪器
func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		total:     total,
		current:   0,
		startTime: time.Now(),
		width:     40,
	}
}

// Increment 增加进度
func (p *ProgressTracker) Increment() {
	p.current++
}

// SetMessage 设置当前消息
func (p *ProgressTracker) SetMessage(msg string) {
	p.lastMessage = msg
}

// Render 渲染进度条
func (p *ProgressTracker) Render() string {
	if p.total == 0 {
		return ""
	}

	percentage := float64(p.current) / float64(p.total)
	filled := int(percentage * float64(p.width))

	// 构建进度条
	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	percent := int(percentage * 100)

	// 计算已用时间
	elapsed := time.Since(p.startTime)
	elapsedStr := formatDuration(elapsed)

	// 估算剩余时间
	var etaStr string
	if p.current > 0 {
		avgTime := elapsed / time.Duration(p.current)
		remaining := time.Duration(p.total-p.current) * avgTime
		etaStr = formatDuration(remaining)
	}

	progressText := fmt.Sprintf("[%s] %d%% (%d/%d)",
		StyleSuccess.Render(bar),
		percent,
		p.current,
		p.total,
	)

	timeInfo := StyleMuted.Render(fmt.Sprintf("已用 %s", elapsedStr))
	if etaStr != "" && p.current < p.total {
		timeInfo += StyleMuted.Render(fmt.Sprintf(" | 预计剩余 %s", etaStr))
	}

	result := progressText + "\n" + timeInfo

	if p.lastMessage != "" {
		result += "\n" + StyleInfo.Render(p.lastMessage)
	}

	return result
}

// RenderSimple 渲染简单版本的进度（单行）
func (p *ProgressTracker) RenderSimple() string {
	if p.total == 0 {
		return ""
	}

	percentage := float64(p.current) / float64(p.total)
	percent := int(percentage * 100)

	return StyleMuted.Render(fmt.Sprintf("进度: %d%% (%d/%d)", percent, p.current, p.total))
}

// formatDuration 格式化时间长度
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

// SimpleProgress 简单进度条（用于静态显示）
type SimpleProgress struct {
	prog progress.Model
}

// NewSimpleProgress 创建简单进度条
func NewSimpleProgress() *SimpleProgress {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return &SimpleProgress{
		prog: prog,
	}
}

// Render 渲染进度条
func (sp *SimpleProgress) Render(percent float64) string {
	return sp.prog.ViewAs(percent)
}

// SpinnerFrames Spinner 动画帧
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner 简单的 Spinner
type Spinner struct {
	frame   int
	message string
}

// NewSpinner 创建新的 Spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frame:   0,
		message: message,
	}
}

// Next 下一帧
func (s *Spinner) Next() {
	s.frame = (s.frame + 1) % len(SpinnerFrames)
}

// SetMessage 设置消息
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// Render 渲染 Spinner
func (s *Spinner) Render() string {
	spinnerStyle := lipgloss.NewStyle().Foreground(ColorPrimary)
	return spinnerStyle.Render(SpinnerFrames[s.frame]) + " " + s.message
}

