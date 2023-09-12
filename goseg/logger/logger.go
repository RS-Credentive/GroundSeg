package logger

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	logPath     string
	logFile     *os.File
	multiWriter io.Writer
	Logger      *slog.Logger
)

type MuMultiWriter struct {
	Writers []io.Writer
	Mu      sync.Mutex
}

func init() {
	fmt.Println("                                       !G#:\n                                   " +
		" .7G@@@^\n          .                       :J#@@@@P.\n     .75GB#BG57.                ~5&@@" +
		"@&Y^  \n    ?&@@@@@@@@@&J             !G@@@@B?. .^ \n   Y@@@@@@@@@@@@@J         :?B@@@@G!  :" +
		"Y&&:\n   @@@@@@@@@@@@@@B       ^Y&@@@&5^  ~P&@@@:\n   Y@@@@@@@@@@@@@J     !P&@@@#J.  7B@@@@G" +
		"! \n    ?#@@@@@@@@@&?   .7B@@@@G7  :J#@@@&5~   \n     .!YGBBBGY7.  :J#@@@&P~  ^5&@@@#J:  !?." +
		"\n                ^5&@@@#J:  !G@@@@G7. .?B@@:\n             .7G@@@@G7. :J#@@@&P~  ^Y#@@@&:\n" +
		"            .P&&&&G!   ~B&&&#5^   ~#&&&&&P. \n\nＮａｔｉｖｅ Ｐｌａｎｅｔ")
	fmt.Println(" ▄▄ • ▄▄▄        ▄• ▄▌ ▐ ▄ ·▄▄▄▄  .▄▄ · ▄▄▄ . ▄▄ • 𝐯𝟐!\n▐█ ▀ ▪▀▄ █·▪     █▪██▌•" +
		"█▌▐███▪ ██ ▐█ ▀. ▀▄.▀·▐█ ▀ ▪\n▄█ ▀█▄▐▀▀▄  ▄█▀▄ █▌▐█▌▐█▐▐▌▐█· ▐█▌▄▀▀▀█▄▐▀▀▪▄▄█ ▀█▄ 🪐\n▐█▄▪▐█▐" +
		"█•█▌▐█▌.▐▌▐█▄█▌██▐█▌██. ██ ▐█▄▪▐█▐█▄▄▌▐█▄▪▐█\n·▀▀▀▀ .▀  ▀ ▀█▄▀▪ ▀▀▀ ▀▀ █▪▀▀▀▀▀•  ▀▀▀▀  ▀▀▀" +
		" ·▀▀▀▀ (~)")
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "/opt/nativeplanet/groundseg/"
	}
	logPath = basePath + "logs/"
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create log directory: %v", err))
		fmt.Println("\n\n.・。.・゜✭・.・✫・゜・。..・。.・゜✭・.・✫・゜・。.")
		fmt.Println("Please run GroundSeg as root!  \n    /) /)\n   ( . . )" +
			"\n   (  >< )\n Love, Native Planet")
		fmt.Println(".・。.・゜✭・.・✫・゜・。..・。.・゜✭・.・✫・゜・。.\n\n")
		panic("")
	}
	logFile, err := os.OpenFile(SysLogfile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}

	multiWriter = muMultiWriter(os.Stdout, logFile)
	Logger = slog.New(slog.NewJSONHandler(multiWriter, nil))
}

func SysLogfile() string {
	currentTime := time.Now()
	return fmt.Sprintf("%s%d-%02d.log", logPath, currentTime.Year(), currentTime.Month())
}

func muMultiWriter(writers ...io.Writer) *MuMultiWriter {
	return &MuMultiWriter{
		Writers: writers,
	}
}

func (m *MuMultiWriter) Write(p []byte) (n int, err error) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	var firstError error
	for _, w := range m.Writers {
		n, err := w.Write(p)
		if err != nil && firstError == nil {
			firstError = err
		}
		if n != len(p) && firstError == nil {
			firstError = io.ErrShortWrite
		}
	}
	return len(p), firstError
}

func TailLogs(filename string, n int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	return lines, scanner.Err()
}
