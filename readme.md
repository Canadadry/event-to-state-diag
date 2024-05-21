State Transition Diagram Generator

This project provides a tool for generating state transition diagrams in Mermaid format from CSV files. The tool allows filtering transitions based on a minimum value, enabling users to focus on significant transitions.

## Features

- **CSV to Mermaid Conversion**: Converts a CSV file containing state transition data into a Mermaid diagram.
- **Transition Filtering**: Filters out transitions below a specified threshold to simplify the output diagram.

## Installation

To use this tool, you will need Go installed on your system. You can download and install Go from [the official Go website](https://golang.org/dl/).

Once Go is installed, clone this repository to your local machine using:

git clone https://github.com/yourusername/yourrepository.git
cd yourrepository

Build the project using:

```bash
go build -o state-diagram-generator
```

This will compile the source code into an executable named `state-diagram-generator`.

## Usage

To run the application, you need to provide it with a CSV file and optionally set a minimum transition value. Below is the syntax for running the tool:

```bash
./state-diagram-generator -in <path-to-csv-file> -min <minimum-transition-value>
```

### Flags

- **-in**: Specifies the path to the CSV file that contains the state transition data.
- **-min**: (Optional) Specifies the minimum transition value to include in the diagram. Transitions with values below this will be omitted from the output. The default value is 0, which includes all transitions.

## Output

The output will be printed directly to the standard output (stdout) in Mermaid format, which you can copy and paste into any Mermaid-compatible viewer to visualize the state transition diagram.
