package main

import (
    "flag"
    "fmt"
    "math"
    "strings"
    "time"
)

func wave(period, amp float64, offset int, sig chan string) {
    for {
        for x := 1.0; x >= 0.0; x -= 1.0 / period {
            l := int(math.Ceil((2 * math.Asin(x) / math.Pi) * amp))
            r := int(math.Floor((amp - float64(l)) * 2))
            e := int(amp)*2 - l - r
            if x == 1.0 {
                e--
            }
            sig <- fmt.Sprintf("%*s%*s%*s", l+offset, "*",
                r+offset, "*", e+offset, "")
        }
        for x := 1.0 / period; x <= 1.0; x += 1.0 / period {
            l := int(math.Ceil((2 * math.Asin(x) / math.Pi) * amp))
            r := int(math.Floor((amp - float64(l)) * 2))
            e := int(amp)*2 - l - r
            sig <- fmt.Sprintf("%*s%*s%*s", l+offset, "*",
                r+offset, "*", e+offset, "")
        }
    }
}

func main() {
    // parse arguments
    period := flag.Int("p", 20, "Period: length (lines) of wave")
    amp := flag.Int("a", 20, "Amplitude: width (chars) of wave")
    num := flag.Int("n", 2, "Number of waves")
    freq := flag.Int("f", 20, "Frequency of waves (Hz)")
    flag.Parse()
    sigs := make([]chan string, *num)
    // start up the waves
    for i, _ := range sigs {
        sigs[i] = make(chan string)
        go wave(float64(*period), float64(*amp), 0, sigs[i])
    }
    // sync up the waves at some frequency in Hz
    go func(frequency int, sigs []chan string) {
        for {
            line := make([]string, len(sigs))
            for i, s := range sigs {
                line[i] = <-s
            }
            // join up the parts and put a newline at the end
            fmt.Printf("%s\n", strings.Join(line, ""))
            time.Sleep(time.Second / time.Duration(frequency))
        }
    }(*freq, sigs)
    // press any key + enter to exit
    var exit string
    fmt.Scan(&exit)
}
