package matrix

import (
	"bytes"
	"encoding/csv"
	"reflect"
	"testing"
	"time"
)

func TestSortAndBuildTransitionMatrix(t *testing.T) {
	tests := map[string]struct {
		in             []Event
		startEventName string
		stopEventName  string
		out            map[string]map[string]int
	}{
		"simple transition": {
			startEventName: "start",
			stopEventName:  "stop",
			in: []Event{
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)},
				{Name: "Event2", RunID: 1, Date: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC)},
			},
			out: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event2": 1},
				"Event2": {"stop": 1},
			},
		},
		"single event name": {
			startEventName: "start",
			stopEventName:  "stop",
			in: []Event{
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)},
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC)},
			},
			out: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event1": 1, "stop": 1},
			},
		},
		"simple transition unordered": {
			startEventName: "start",
			stopEventName:  "stop",
			in: []Event{
				{Name: "Event2", RunID: 1, Date: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC)},
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)},
			},
			out: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event2": 1},
				"Event2": {"stop": 1},
			},
		},
		"return transition": {
			startEventName: "start",
			stopEventName:  "stop",
			in: []Event{
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)},
				{Name: "Event2", RunID: 1, Date: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC)},
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 15, 0, 0, 0, 0, time.UTC)},
			},
			out: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event2": 1, "stop": 1},
				"Event2": {"Event1": 1},
			},
		},
		"interleaved runs": {
			startEventName: "start",
			stopEventName:  "stop",
			in: []Event{
				{Name: "Event1", RunID: 1, Date: time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)},
				{Name: "Event2", RunID: 1, Date: time.Date(2022, 5, 11, 0, 0, 0, 0, time.UTC)},
				{Name: "Event1", RunID: 2, Date: time.Date(2022, 5, 12, 0, 0, 0, 0, time.UTC)},
				{Name: "Event2", RunID: 2, Date: time.Date(2022, 5, 13, 0, 0, 0, 0, time.UTC)},
				{Name: "Event3", RunID: 1, Date: time.Date(2022, 5, 14, 0, 0, 0, 0, time.UTC)},
			},
			out: map[string]map[string]int{
				"start":  {"Event1": 2},
				"Event1": {"Event2": 2},
				"Event2": {"Event3": 1, "stop": 1},
				"Event3": {"stop": 1},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := buildTransitionMatrix(tt.in, tt.startEventName, tt.stopEventName)
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("Test %s failed, expected %v, got %v", name, tt.out, got)
			}
		})
	}
}

// Test function for saveTransitionMatrix
func TestSaveTransitionMatrix(t *testing.T) {
	tests := map[string]struct {
		in  map[string]map[string]int
		out [][]string
	}{
		"empty matrix": {
			in:  map[string]map[string]int{},
			out: [][]string{{"From/To"}},
		},
		"single transition": {
			in: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event2": 1},
				"Event2": {"stop": 1},
			},
			out: [][]string{
				{"From/To", "Event1", "Event2", "start", "stop"},
				{"Event1", "0", "1", "0", "0"},
				{"Event2", "0", "0", "0", "1"},
				{"start", "1", "0", "0", "0"},
				{"stop", "0", "0", "0", "0"},
			},
		},
		"multiple transitions": {
			in: map[string]map[string]int{
				"start":  {"Event1": 1},
				"Event1": {"Event2": 2, "Event3": 3},
				"Event2": {"Event1": 1},
				"Event3": {"stop": 1},
			},
			out: [][]string{
				{"From/To", "Event1", "Event2", "Event3", "start", "stop"},
				{"Event1", "0", "2", "3", "0", "0"},
				{"Event2", "1", "0", "0", "0", "0"},
				{"Event3", "0", "0", "0", "0", "1"},
				{"start", "1", "0", "0", "0", "0"},
				{"stop", "0", "0", "0", "0", "0"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Use a buffer to capture output from io.Writer
			var buf bytes.Buffer
			saveTransitionMatrix(&buf, tt.in)

			r := csv.NewReader(&buf)
			result, err := r.ReadAll()
			if err != nil {
				t.Fatalf("Failed to read from buffer: %v", err)
			}

			// Compare expected and actual results
			if !reflect.DeepEqual(result, tt.out) {
				t.Errorf("Test %s failed, expected %v, got %v", name, tt.out, result)
			}
		})
	}
}
