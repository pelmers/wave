package main

import (
    "flag"
    "fmt"
    "math"
    "strings"
    "time"
)

// Wave draws asin curve with given period and amplitude and sends out
// the string on given channel.
// The length of each string is amp*2. Spaces are padded to the front and end.
func Wave(period, amp float64, offset int, sig chan<- string) {
    step := 1.0/period
    for {
        for x := 1.0; x >= 0.0; x -= step {
            l := int(math.Ceil((2 * math.Asin(x) / math.Pi) * amp))
            r := int(math.Floor((amp - float64(l)) * 2))
            e := int(amp)*2 - l - r
            // fix rounding misalignment
            if x == 1.0 {
                e--
            }
            sig <- fmt.Sprintf("%*s%*s%*s", l+offset, "*",
                r+offset, "*", e+offset, "")
        }
        for x := step; x < 1.0; x += step {
            l := int(math.Ceil((2 * math.Asin(x) / math.Pi) * amp))
            r := int(math.Floor((amp - float64(l)) * 2))
            e := int(amp)*2 - l - r
            // fix rounding misalignment
            if x > 1.0 - step {
                e--
            }
            sig <- fmt.Sprintf("%*s%*s%*s", l+offset, "*",
                r+offset, "*", e+offset, "")
        }
    }
}

func main() {
    // parse arguments
    period := flag.Int("p", 25, "Period: length (lines) of wave")
    amp := flag.Int("a", 25, "Amplitude: width (chars) of wave")
    num := flag.Int("n", 2, "Number of waves")
    freq := flag.Int("f", 40, "Frequency of waves (Hz)")
    flag.Parse()
    sigs := make([]chan string, *num)
    // start up the waves
    for i, _ := range sigs {
        sigs[i] = make(chan string)
        go Wave(float64(*period), float64(*amp), 0, sigs[i])
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
