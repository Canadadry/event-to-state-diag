package matrix

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const Name = "matrix"

type Event struct {
	Name  string
	RunID int
	Date  time.Time
}

// Function to group and sort events by runID and then by date
func sortEventsByRunAndDate(events []Event) map[int][]Event {
	groupedEvents := make(map[int][]Event)
	for _, event := range events {
		groupedEvents[event.RunID] = append(groupedEvents[event.RunID], event)
	}

	for runID, events := range groupedEvents {
		sort.Slice(events, func(i, j int) bool {
			return events[i].Date.Before(events[j].Date)
		})
		groupedEvents[runID] = events // Important to update the map after sorting
	}
	return groupedEvents
}

// Function to build the transition matrix for each run
func buildTransitionMatrix(events []Event) map[string]map[string]int {
	groupedEvents := sortEventsByRunAndDate(events)
	transitionMatrix := make(map[string]map[string]int)

	for _, events := range groupedEvents {
		var previousEvent Event
		for _, event := range events {
			if previousEvent.Name != "" {
				if _, ok := transitionMatrix[previousEvent.Name]; !ok {
					transitionMatrix[previousEvent.Name] = make(map[string]int)
				}
				transitionMatrix[previousEvent.Name][event.Name]++
			}
			previousEvent = event
		}
	}
	return transitionMatrix
}

// Function to read and parse the CSV file, grouped by kind_id
func readCSVByKind(filename string) (map[int][]Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.TrimLeadingSpace = true
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, err
	}

	eventsByKind := make(map[int][]Event)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		sentAt := strings.TrimSpace(record[10])
		if sentAt == "" {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05", sentAt)
		if err != nil {
			log.Printf("Skipping record with invalid date: %s", sentAt)
			continue
		}

		kindID, err := strconv.Atoi(record[6])
		if err != nil {
			log.Printf("Skipping record with invalid kind_id: %s", record[6])
			continue
		}

		experienceID, err := strconv.Atoi(record[5])
		if err != nil {
			log.Printf("Skipping record with invalid experience_id: %s", record[5])
			continue
		}

		event := Event{
			Name:  record[2],
			RunID: experienceID,
			Date:  date,
		}
		eventsByKind[kindID] = append(eventsByKind[kindID], event)
	}
	return eventsByKind, nil
}

// Save transition matrix to a CSV format using an io.Writer, for a more flexible output
func saveTransitionMatrix(writer io.Writer, matrix map[string]map[string]int) {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Gather all event names to create sorted headers
	eventNames := make(map[string]bool)
	for fromEvent, transitions := range matrix {
		eventNames[fromEvent] = true
		for toEvent := range transitions {
			eventNames[toEvent] = true
		}
	}

	// Convert map keys to a slice and sort them
	var headers []string
	for eventName := range eventNames {
		headers = append(headers, eventName)
	}
	sort.Strings(headers) // Sort headers for consistent output

	// Write headers to CSV
	csvWriter.Write(append([]string{"From/To"}, headers...)) // First column has "From/To"

	// Write data rows
	for _, fromEvent := range headers {
		row := make([]string, len(headers)+1)
		row[0] = fromEvent // First column of each row is the fromEvent name
		for j, toEvent := range headers {
			if transitions, ok := matrix[fromEvent]; ok {
				count, ok := transitions[toEvent]
				if ok {
					row[j+1] = strconv.Itoa(count) // j+1 because the first column is fromEvent
				} else {
					row[j+1] = "0" // Fill with zero if no transition exists
				}
			} else {
				row[j+1] = "0" // Fill with zero if no transitions exist
			}
		}
		csvWriter.Write(row)
	}
}

func Run(args []string) error {
	filename := ""
	fs := flag.NewFlagSet(Name, flag.ContinueOnError)
	fs.StringVar(&filename, "in", "", "csv of events")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("cannot read app args : %w", err)
	}

	eventsByKind, err := readCSVByKind(filename)
	if err != nil {
		return fmt.Errorf("cannot read csv : %w", err)
	}

	for kindID, events := range eventsByKind {
		fmt.Printf("Transition matrix for kind_id %d:\n", kindID)
		transitionMatrix := buildTransitionMatrix(events)
		saveTransitionMatrix(os.Stdout, transitionMatrix)
	}
	return nil
}
