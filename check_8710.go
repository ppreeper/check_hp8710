package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

var host = flag.String("H", "", "Printer to query")

// toners struct
type toners struct {
	cyanMax    string
	cyanLvl    string
	magentaMax string
	magentaLvl string
	yellowMax  string
	yellowLvl  string
	kromaMax   string
	kromaLvl   string
}

var hp = &toners{
	".1.3.6.1.2.1.43.11.1.1.8.1.2", ".1.3.6.1.2.1.43.11.1.1.9.1.2",
	".1.3.6.1.2.1.43.11.1.1.8.1.3", ".1.3.6.1.2.1.43.11.1.1.9.1.3",
	".1.3.6.1.2.1.43.11.1.1.8.1.1", ".1.3.6.1.2.1.43.11.1.1.9.1.1",
	".1.3.6.1.2.1.43.11.1.1.8.1.4", ".1.3.6.1.2.1.43.11.1.1.9.1.4",
}

func tonerOutput(color string, maxValue string, lvlValue string) string {
	color = strings.ToUpper(color)
	max, errm := strconv.Atoi(maxValue)
	lvl, errl := strconv.Atoi(lvlValue)
	var output string

	if errm == nil && errl == nil {
		if max != 0 {
			level := 100 * float64(float64(lvl)/float64(max))
			tLevels := "\tToner Levels: " + lvlValue + " of " + maxValue + " \t| "
			tLevels += strconv.FormatFloat(level, 'f', 0, 64) + "\n"
			if level >= 10.0 {
				output = color + " Toner\tOK " + tLevels
			} else {
				output = color + " Toner\tLOW " + tLevels
			}
		}
	}
	return output
}

func tonerLevel(color string) string {
	color = strings.ToUpper(color)
	var tonerColor string
	var t toners
	var output string
	var max, lvl string

	t = *hp

	if color == "C" {
		tonerColor = "CYAN"
		max, _ = getSNMPValue(t.cyanMax)
		lvl, _ = getSNMPValue(t.cyanLvl)
	}
	if color == "M" {
		tonerColor = "MAGENTA"
		max, _ = getSNMPValue(t.magentaMax)
		lvl, _ = getSNMPValue(t.magentaLvl)
	}
	if color == "Y" {
		tonerColor = "YELLOW"
		max, _ = getSNMPValue(t.yellowMax)
		lvl, _ = getSNMPValue(t.yellowLvl)
	}
	if color == "K" {
		tonerColor = "BLACK"
		max, _ = getSNMPValue(t.kromaMax)
		lvl, _ = getSNMPValue(t.kromaLvl)
	}
	output = tonerOutput(tonerColor, max, lvl)
	return output
}

func getSNMPValue(oid string) (string, error) {
	g.Default.Target = *host
	err := g.Default.Connect()
	if err != nil {
		return "", fmt.Errorf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()
	oids := []string{oid}
	result, err := g.Default.Get(oids)
	if err != nil {
		return "", fmt.Errorf("Get() err: %v", err)
	}
	return fmt.Sprintf("%s", g.ToBigInt(result.Variables[0].Value)), err
}

type levels struct {
	level [4]string
}

// main function
func main() {
	flag.Parse()
	var colors [4]string
	colors[0] = "C"
	colors[1] = "M"
	colors[2] = "Y"
	colors[3] = "K"
	if *host == "" {
		fmt.Fprintf(os.Stdout, "Host not set\n")
	} else {
		for i := 0; i < len(colors); i++ {
			r := tonerLevel(colors[i])
			fmt.Fprintf(os.Stdout, "%s", r)
		}
	}
	return
}
