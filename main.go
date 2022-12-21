package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	Path      string
	Frequency int   // dir used how many times?
	Time      int64 // seconds since epoch
}

var (
	exitCodeCD = 3
	data       = make([]Entry, 0, 100)

	version  = "dev"
	helpText = `WarpDrive - Warp across the filesystem instantly

Usage: wd [<pattern> ...]
       wd --list/-l | --help/-h | --version/-v
       wd {--add/-a | --remove/-r} <path>
Options:
   --list/-l     List currently tracked paths along with their frecency scores
   --add/-a      Add a path to the data file (paths will be added automatically)
   --remove/-r   Remove a path from the data file
   --help/-h     Print this help message
   --version/-v  Print the version of WarpDrive installed
Examples:
   wd -l                          # list all tracked paths and scores
   wd                             # cd to home directory
   wd s                           # tries to match 's'
   wd someDir                     # tries to match 'someDir'
   wd some subDir                 # ensures matched path also contains 'some'
   wd /absolute/path/to/someDir   # absolute paths work too
Note:
   When specifying multiple patterns, order does not matter, except for the last pattern.
   WarpDrive will always cd to a directory that matches the last pattern.`
)

func main() {
	if ok, _ := ArgsHaveOption("help", "h"); ok {
		fmt.Println(helpText)
		return
	}

	if ok, _ := ArgsHaveOption("version", "v"); ok {
		fmt.Println("WarpDrive " + version)
		return
	}

	dataFile, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.MkdirAll(filepath.Join(dataFile, ".config", "warpdrive"), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	dataFile = filepath.Join(dataFile, ".config", "warpdrive", "data.json")

	data, err = readFromFileAsJSON(dataFile)
	data = normalizeData()

	retCode := 0 // used to exit with non-zero code later on if needed
	defer func() {
		err = writeToFileAsJSON(data, dataFile)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(retCode)
	}()

	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	if ok, _ := ArgsHaveOption("list", "l"); ok {
		listAllPaths()
		return
	}

	if ok, i := ArgsHaveOption("add", "a"); ok {
		if len(os.Args) < i+2 {
			fmt.Println("option --add requires an argument")
			return
		}
		path := normalizePath(os.Args[i+1])
		if s, err := os.Stat(path); err == nil {
			if !s.IsDir() {
				fmt.Println(path + " is not a directory")
				return
			}
		} else {
			fmt.Println(err)
			return
		}
		addPath(path)
		return
	}

	if ok, i := ArgsHaveOption("remove", "r"); ok {
		if len(os.Args) < i+2 {
			fmt.Println("option --remove requires an argument")
			return
		}
		data = removePath(normalizePath(os.Args[i+1]))
		return
	}

	if len(os.Args) == 1 { // go to home dir
		retCode = exitCodeCD
		return
	}

	if os.Args[1] == "-" { // don't match on -, simply let cd handle it (will go to previous pwd)
		fmt.Println("-")
		retCode = exitCodeCD
		return
	}
	patternToFind := strings.Join(os.Args[1:], " ")
	match := getBestMatch(patternToFind)
	addPath(match)
	fmt.Println(match)
	retCode = exitCodeCD // deferred os.Exit() call will use this exit code.
}

func normalizeData() []Entry {
	sortData()
	temp := data[:0] // uses the same underlying array: https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	for _, entry := range data {
		exists, err := pathExists(entry.Path)
		if score(entry) > 1 && exists {
			temp = append(temp, entry)
		}
		if err != nil {
			errPrintln(err)
		}
	}
	return data
}

func pathExists(path string) (bool, error) {
	// check if path exists
	_, err := os.Stat(path)
	//errPrintln(err, "s;osi")
	if os.IsNotExist(err) {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func errPrintln(a ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintln(a...))
}

func listAllPaths() {
	fmt.Println("Score\t\tPath")
	for _, e := range data {
		s := score(e)
		if s < 10 {
			fmt.Printf("%.1f\t\t%s\n", s, e.Path)
			continue
		}
		fmt.Printf("%.0f\t\t%s\n", s, e.Path)
	}
}

func addPath(path string) {
	for j, entry := range data {
		if entry.Path == path {
			data[j].Frequency++
			data[j].Time = time.Now().Unix()
			break
		}
		if j == len(data)-1 { // at the end of the list, and we still haven't added anything
			data = append(data, Entry{
				Path:      path,
				Frequency: 1,
				Time:      time.Now().Unix(),
			})
		}
	}
	if len(data) == 0 {
		data = append(data, Entry{
			Path:      path,
			Frequency: 1,
			Time:      time.Now().Unix(),
		})
	}
}

func removePath(path string) []Entry {
	for j, entry := range data {
		if entry.Path == path {
			data = append(data[:j], data[j+1:]...)
			return data
		}
	}
	return data
}

func normalizePath(path string) string {
	s, err := filepath.Abs(path)
	if err != nil { // errors only if the current working directory can't be gotten
		panic(err)
	}
	return s
}

func getBestMatch(pattern string) string {
	candidates := make([]Entry, 0, 10)
	for _, entry := range data {
		pattenSplitOnSpaces := strings.Split(pattern, " ")
		splitCandidate := strings.Split(entry.Path, "/")
		// path must have all the words in the pattern
		if stringContainsAllElemsInArr(entry.Path, pattenSplitOnSpaces) {
			// path's last element must have the pattern's last element
			patternLastElemSplitOnSlash := strings.Split(pattenSplitOnSpaces[len(pattenSplitOnSpaces)-1], string(os.PathSeparator))
			if strings.Contains(strings.ToLower(splitCandidate[len(splitCandidate)-1]),
				strings.ToLower(patternLastElemSplitOnSlash[len(patternLastElemSplitOnSlash)-1])) {
				candidates = append(candidates, entry)
			}
		}
	}
	sortData()
	if len(candidates) == 0 {
		return normalizePath(pattern)
	}
	return candidates[0].Path
}

func stringContainsAllElemsInArr(searchIn string, searchArr []string) bool {
	for _, s := range searchArr {
		if !strings.Contains(strings.ToLower(searchIn), strings.ToLower(s)) {
			return false
		}
	}
	return true
}

func sortData() {
	sort.Slice(data, func(i, j int) bool {
		return score(data[i]) > score(data[j])
	})
}

func score(e Entry) float64 {
	t := time.Unix(e.Time, 0)
	if time.Since(t) > time.Hour*24*30 { // haven't been there in a month?
		return 0
	}
	//log := math.Round(100 * float64(e.Frequency) * math.Exp(-0.0001*time.Since(t).Minutes()))
	result := 1000 * float64(e.Frequency) * math.Exp2(
		-(time.Since(t).Minutes()+58000)/(9*60*24), // halve every nine days
	)
	//os.Stderr.WriteString(fmt.Sprintln(e.Path, e.Frequency, time.Since(t).Minutes(), result))
	//if result < 0 {
	//	return 0
	//}
	return result
	//return math.Round(
	//	300 * math.Sqrt(
	//		float64(e.Frequency)/math.Sqrt(time.Since(t).Minutes()),
	//	),
	//)
}

// ArgsHaveOption checks command line arguments for an option
func ArgsHaveOption(long, short string) (hasOption bool, foundAt int) {
	for i, arg := range os.Args {
		if arg == "--"+long || arg == "-"+short {
			return true, i
		}
	}
	return false, 0
}

func writeToFileAsJSON(data []Entry, fileName string) error {
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, b, 0644)
}

func readFromFileAsJSON(fileName string) ([]Entry, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var dataRead []Entry
	err = json.Unmarshal(b, &dataRead)
	if err != nil {
		return nil, err
	}
	return dataRead, nil
}
