package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/patnaikankit/Forward-Proxy/utils"
)

var windowFocus = "requests" // "requests" or "blocked"
var selectedRequestIndex int
var selectedBlockedIndex int
var showDetailsPane bool
var detailsPaneRequest string
var showDetailsScreen bool
var detailsScreenRequest *utils.LoggedRequest

func ConsoleRunner() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	windowFocus = "requests"
	selectedRequestIndex = 0
	selectedBlockedIndex = 0
	showDetailsPane = false
	detailsPaneRequest = ""
	showDetailsScreen = false
	detailsScreenRequest = nil

	drawUI()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				drawUI()
			case <-quit:
				return
			}
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if showDetailsScreen {
				if ev.Key == termbox.KeyEnter || ev.Key == termbox.KeyEsc {
					showDetailsScreen = false
					drawUI()
				}
				continue
			}
			if ev.Key == termbox.KeyCtrlC || ev.Ch == 'q' {
				close(quit)
				return
			} else if ev.Ch == 's' || ev.Ch == 'S' {
				if windowFocus == "requests" {
					windowFocus = "blocked"
				} else {
					windowFocus = "requests"
				}
				drawUI()
			} else if ev.Key == termbox.KeyArrowUp {
				if windowFocus == "requests" {
					if selectedRequestIndex > 0 {
						selectedRequestIndex--
						drawUI()
					}
				} else {
					if selectedBlockedIndex > 0 {
						selectedBlockedIndex--
						drawUI()
					}
				}
			} else if ev.Key == termbox.KeyArrowDown {
				if windowFocus == "requests" {
					requests := utils.GetLogs()
					if selectedRequestIndex < len(requests)-1 {
						selectedRequestIndex++
						drawUI()
					}
				} else {
					blockedURLs := utils.GetBlockedURLs()
					if selectedBlockedIndex < len(blockedURLs)-1 {
						selectedBlockedIndex++
						drawUI()
					}
				}
			} else if (ev.Ch == 'b' || ev.Ch == 'B') && windowFocus == "requests" {
				requests := utils.GetLogs()
				if selectedRequestIndex >= 0 && selectedRequestIndex < len(requests) {
					url := extractURL(requests[selectedRequestIndex].LogLine)
					if url != "" {
						utils.BlockURL(url)
						drawUI()
					}
				}
			} else if (ev.Ch == 'u' || ev.Ch == 'U') && windowFocus == "blocked" {
				blockedURLs := utils.GetBlockedURLs()
				if selectedBlockedIndex >= 0 && selectedBlockedIndex < len(blockedURLs) {
					url := blockedURLs[selectedBlockedIndex]
					if url != "" {
						utils.UnblockURL(url)
						if selectedBlockedIndex >= len(blockedURLs)-1 && selectedBlockedIndex > 0 {
							selectedBlockedIndex--
						}
						drawUI()
					}
				}
			} else if ev.Key == termbox.KeyEnter && windowFocus == "requests" {
				requests := utils.GetLogs()
				if selectedRequestIndex >= 0 && selectedRequestIndex < len(requests) {
					showDetailsScreen = true
					detailsScreenRequest = &requests[selectedRequestIndex]
					drawUI()
				}
			}
		case termbox.EventResize:
			drawUI()
		}
	}
}

// Box/line drawing helpers
func drawHLine(y, x1, x2 int) {
	for x := x1; x <= x2; x++ {
		termbox.SetCell(x, y, '-', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func drawVLine(x, y1, y2 int) {
	for y := y1; y <= y2; y++ {
		termbox.SetCell(x, y, '|', termbox.ColorWhite, termbox.ColorBlack)
	}
}

func drawBox(x1, y1, x2, y2 int) {
	drawHLine(y1, x1, x2)
	drawHLine(y2, x1, x2)
	drawVLine(x1, y1, y2)
	drawVLine(x2, y1, y2)
	// Corners
	termbox.SetCell(x1, y1, '+', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(x2, y1, '+', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(x1, y2, '+', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(x2, y2, '+', termbox.ColorWhite, termbox.ColorBlack)
}

func drawUI() {
	if showDetailsScreen && detailsScreenRequest != nil {
		drawFullScreenDetails(detailsScreenRequest)
		termbox.Flush()
		return
	}
	requests := utils.GetLogs()
	blockedURLs := utils.GetBlockedURLs()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	// Layout: Top half for requests, bottom half split for keymap and blocked URLs
	midY := height / 2
	if midY < 5 {
		midY = 5
	}
	bottomY := height - 1
	leftBoxX2 := width/2 - 1
	if leftBoxX2 < 20 {
		leftBoxX2 = 20
	}

	// Draw requests box (top half)
	drawBox(0, 0, width-1, midY)
	header := fmt.Sprintf("Management Console - %d Total Packets", len(requests))
	writeStr((width-len(header))/2, 0, header, termbox.ColorYellow, termbox.ColorBlack)

	maxRequests := midY - 2
	startIdx := 0
	if len(requests) > maxRequests {
		startIdx = len(requests) - maxRequests
	}
	if selectedRequestIndex < 0 {
		selectedRequestIndex = 0
	}
	if selectedRequestIndex > len(requests)-1 {
		selectedRequestIndex = len(requests) - 1
	}
	for i := startIdx; i < len(requests) && i-startIdx < maxRequests; i++ {
		line := requests[i].LogLine
		if len(line) > width-4 {
			line = line[:width-7] + "..."
		}
		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if windowFocus == "requests" && i == selectedRequestIndex {
			fg, bg = termbox.ColorBlack, termbox.ColorWhite // Highlight selected
		}
		writeStr(2, i-startIdx+1, line, fg, bg)
	}

	// Draw keymap box (bottom left)
	drawBox(0, midY+1, leftBoxX2, bottomY)
	keymapY := midY + 2
	writeStr(2, keymapY, "Keymap:", termbox.ColorCyan, termbox.ColorBlack)
	writeStr(2, keymapY+1, "[Q] - Quit", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+2, "[S] - Switch window - Requests/Blocked pane", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+3, "[R] - Refresh requests list", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+4, "[B] - Block URL (of selected packet)", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+5, "[U] - Unblock URL (of selected blocked)", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+6, "[Up/Down Arrow] - Select packet or URL", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(2, keymapY+7, "[Enter] - Inspect HTTP requests", termbox.ColorWhite, termbox.ColorBlack)

	// Draw blocked URLs box (bottom right)
	drawBox(leftBoxX2+1, midY+1, width-1, bottomY)
	writeStr(leftBoxX2+3, midY+2, "Blocked URLs", termbox.ColorRed, termbox.ColorBlack)
	maxBlocked := bottomY - (midY + 3)
	for i := 0; i < len(blockedURLs) && i < maxBlocked; i++ {
		url := blockedURLs[i]
		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if windowFocus == "blocked" && i == selectedBlockedIndex {
			fg, bg = termbox.ColorBlack, termbox.ColorWhite
		}
		writeStr(leftBoxX2+3, midY+3+i, url, fg, bg)
	}

	// Draw status line at the very bottom (optional)
	status := ""
	if windowFocus == "requests" && selectedRequestIndex >= 0 && selectedRequestIndex < len(requests) {
		status = requests[selectedRequestIndex].LogLine
	} else if windowFocus == "blocked" && selectedBlockedIndex >= 0 && selectedBlockedIndex < len(blockedURLs) {
		status = blockedURLs[selectedBlockedIndex]
	}
	if status != "" && height > 1 {
		writeStr(2, height-1, "Selected: "+status, termbox.ColorYellow, termbox.ColorBlack)
	}

	termbox.Flush()

	if showDetailsPane {
		drawRequestDetailsPane(detailsPaneRequest, width, height)
		termbox.Flush()
	}
}

// Draw a centered details pane with the full request info and cache/raw data if available
func drawRequestDetailsPane(logLine string, width, height int) {
	// Extract URL from log line
	url := extractURL(logLine)
	cache := ""
	if url != "" {
		cached := utils.GetCache(url)
		if cached != nil {
			cache = string(cached)
		}
	}
	// Prepare details text
	details := logLine
	if cache != "" {
		details += "\n\n--- Cached/Raw Data ---\n" + cache
	}
	lines := splitLines(details, width-8)
	boxWidth := width - 8
	if boxWidth > 100 {
		boxWidth = 100
	}
	boxHeight := len(lines) + 4
	if boxHeight > height-4 {
		boxHeight = height - 4
	}
	x1 := (width - boxWidth) / 2
	y1 := (height - boxHeight) / 2
	x2 := x1 + boxWidth
	y2 := y1 + boxHeight
	drawBox(x1, y1, x2, y2)
	for i, line := range lines {
		writeStr(x1+2, y1+2+i, line, termbox.ColorYellow, termbox.ColorBlack)
	}
	writeStr(x1+2, y2-1, "[Esc/Enter] to close", termbox.ColorWhite, termbox.ColorBlack)
}

// Helper to split long text into lines of max width
func splitLines(s string, maxWidth int) []string {
	var lines []string
	curr := ""
	for _, r := range s {
		if r == '\n' || len(curr) >= maxWidth {
			lines = append(lines, curr)
			if r == '\n' {
				curr = ""
			} else {
				curr = string(r)
			}
		} else {
			curr += string(r)
		}
	}
	if curr != "" {
		lines = append(lines, curr)
	}
	return lines
}

// Draw a centered details pane with the full request info
func drawDetailsPane(line string, width, height int) {
	boxWidth := len(line) + 4
	if boxWidth > width-4 {
		boxWidth = width - 4
	}
	boxHeight := 5
	x1 := (width - boxWidth) / 2
	y1 := (height - boxHeight) / 2
	x2 := x1 + boxWidth
	y2 := y1 + boxHeight
	drawBox(x1, y1, x2, y2)
	writeStr(x1+2, y1+2, line, termbox.ColorYellow, termbox.ColorBlack)
	writeStr(x1+2, y2-1, "[Esc/Enter] to close", termbox.ColorWhite, termbox.ColorBlack)
}

// Extract URL from a log line (new format: [time] PROTOCOL --> METHOD URL HTTP/1.1)
func extractURL(line string) string {
	// Example: [21:26:10] HTTPS --> CONNECT push.services.mozilla.com HTTP/1.1
	// We want to extract 'push.services.mozilla.com' (the URL after METHOD)
	parts := line
	// Find the first occurrence of ">" (end of -->), then split by spaces
	arrowIdx := -1
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == '>' && parts[i-1] == '-' && parts[i-2] == '-' && parts[i-3] == ' ' {
			arrowIdx = i + 1
			break
		}
	}
	if arrowIdx == -1 || arrowIdx >= len(parts) {
		return ""
	}
	// Skip spaces after -->
	for arrowIdx < len(parts) && parts[arrowIdx] == ' ' {
		arrowIdx++
	}
	// Now, the next word is METHOD, then URL, then maybe HTTP/1.1
	fields := []string{}
	curr := ""
	for i := arrowIdx; i < len(parts); i++ {
		if parts[i] == ' ' {
			if curr != "" {
				fields = append(fields, curr)
				curr = ""
			}
		} else {
			curr += string(parts[i])
		}
	}
	if curr != "" {
		fields = append(fields, curr)
	}
	if len(fields) < 2 {
		return ""
	}
	return fields[1] // METHOD fields[0], URL fields[1]
}

func drawRequests(requests []utils.LoggedRequest, width, height int) {
	// Header
	header := fmt.Sprintf("Management Console - %d Total Packets", len(requests))
	writeStr(0, 0, header, termbox.ColorYellow, termbox.ColorBlack)

	// Calculate how many requests we can show
	maxRequests := height - 10 // Leave space for keymap
	if maxRequests < 1 {
		maxRequests = 1
	}

	// Show recent requests (limit to screen size)
	startIdx := 0
	if len(requests) > maxRequests {
		startIdx = len(requests) - maxRequests
	}

	for i := startIdx; i < len(requests) && i-startIdx < maxRequests; i++ {
		line := requests[i].LogLine
		// Truncate line if it's too long for terminal width
		if len(line) > width-1 {
			line = line[:width-4] + "..."
		}
		writeStr(0, i-startIdx+2, line, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func drawBlocked(urls []string, startY int) {
	writeStr(50, startY, "Blocked URLs", termbox.ColorRed, termbox.ColorBlack)
	for i, url := range urls {
		writeStr(50, startY+i+1, url, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func drawKeymap(height int) {
	// Position keymap at bottom of screen
	base := height - 8
	if base < 0 {
		base = 0
	}

	writeStr(0, base, "Keymap:", termbox.ColorCyan, termbox.ColorBlack)
	writeStr(0, base+1, "[Q] - Quit", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+2, "[S] - Switch window", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+3, "[R] - Refresh requests list", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+4, "[B] - Block URL (of selected packet)", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+5, "[U] - Unblock URL", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+6, "[Ent/Esc] - View detailed packet info", termbox.ColorWhite, termbox.ColorBlack)
	writeStr(0, base+7, "[Up/Down Arrow] - Select packet or URL", termbox.ColorWhite, termbox.ColorBlack)
}

func writeStr(x, y int, msg string, fg, bg termbox.Attribute) {
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func drawFullScreenDetails(req *utils.LoggedRequest) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, _ := termbox.Size()
	y := 1
	writeStr((width-18)/2, 0, "Request Details", termbox.ColorCyan, termbox.ColorBlack)
	writeStr(2, y, "General Information", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	y++
	writeStr(2, y, fmt.Sprintf("Timestamp - %s", req.Timestamp.Format("15:04:05")), termbox.ColorWhite, termbox.ColorBlack)
	y++
	writeStr(2, y, fmt.Sprintf("Method - %s", req.Request.Method), termbox.ColorWhite, termbox.ColorBlack)
	y++
	writeStr(2, y, fmt.Sprintf("URL - %s", req.Request.URL), termbox.ColorWhite, termbox.ColorBlack)
	y++
	writeStr(2, y, fmt.Sprintf("Version - %s", req.Request.Version), termbox.ColorWhite, termbox.ColorBlack)
	y += 2
	writeStr(2, y, "HTTP Headers", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	y++
	// Parse headers from Raw
	headers := parseHeadersFromRaw(req.Request.Raw)
	for _, h := range headers {
		writeStr(2, y, h, termbox.ColorWhite, termbox.ColorBlack)
		y++
	}
	y++
	writeStr(2, y, "Raw Request Data", termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	y++
	raw := string(req.Request.Raw)
	for _, line := range splitLines(raw, width-4) {
		writeStr(2, y, line, termbox.ColorWhite, termbox.ColorBlack)
		y++
	}
	writeStr(2, y+1, "[Esc/Enter] to return", termbox.ColorYellow, termbox.ColorBlack)
}

// Helper to parse headers from raw HTTP request bytes
func parseHeadersFromRaw(raw []byte) []string {
	lines := strings.Split(string(raw), "\r\n")
	var headers []string
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}
		headers = append(headers, line)
	}
	return headers
}
