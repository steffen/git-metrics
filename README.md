<h1 align="center">git-metrics</h1>

A powerful Git repository analysis tool that provides detailed metrics, growth statistics, future projections, and contributor insights for your Git repositories.

## Overview

`git-metrics` is a command-line utility that analyzes Git repositories to provide comprehensive insights about repository history, structure, and growth patterns. The tool examines your repository's Git object database to reveal historical trends, identify storage-heavy components, and visualize contributor activity over time. With this data, it generates projections for future repository growth and helps identify optimization opportunities.

Key features include:
- Repository metadata analysis (first commit, age)
- Year-by-year growth statistics for Git objects (commits, trees, blobs) and their on-disk size
- Future growth projections based on historical trends
- Directory structure analysis with size impact indicators
- Identification of largest files in the repository
- File extension distribution analysis
- Contributor statistics showing top committers and authors over time

## Examples

<p align="center"><a href="https://steffen.github.io/git-metrics/">Interactive examples</a></p>

## Installation

### Prerequisites
- Git

### Download prebuilt binaries

The easiest way to install `git-metrics` is to download a prebuilt binary from the [GitHub releases page](https://github.com/steffen/git-metrics/releases).

#### Linux

1. Download the latest release for Linux (64-bit):
   ```bash
   curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-linux-amd64.tar.gz -o git-metrics.tar.gz
   ```

2. Unpack the archive:
   ```bash
   tar -xzf git-metrics.tar.gz
   ```

3. _Optional:_ Move it to a directory in your PATH:
   ```bash
   sudo mv git-metrics /usr/local/bin/
   ```

#### macOS

1. Download the latest release for macOS (Intel or Apple Silicon):
   ```bash
   curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-darwin-arm64.zip -o git-metrics.zip
   ```

2. Unzip the archive:
   ```bash
   unzip git-metrics.zip
   ```

3. _Optional:_ Move it to a directory in your PATH:
   ```bash
   sudo mv git-metrics /usr/local/bin/
   ```

4. When downloaded via browser you may need to remove the quarantine attribute in order to run the tool:
   ```bash
   xattr -d com.apple.quarantine git-metrics
   ```

### Running the tool

* Analyze the current directory as a Git repository:
  ```bash
  git-metrics
  ```

* Analyze a specific repository:
  ```bash
  git-metrics -r /path/to/repository
  ```

## Command line options

| Option | Description |
|--------|-------------|
| `-r`, `--repository` | Path to Git repository (default: current directory) |
| `--debug` | Enable debug output |
| `--no-progress` | Disable progress indicators |
| `--version` | Display version information and exit |


## Understanding the output

`git-metrics` provides several sections of output:

1. **Run information**: Details about when, where, and with which versions the tool was executed.
2. **Repository information**: Basic metadata about your repository including path, remote URL, age, and commit history.
3. **Historic & estimated growth**: Year-by-year breakdown of Git object growth (commits, trees, blobs) and disk usage, with future projections based on historical trends.
4. **Largest directories**: Hierarchical view of directory sizes and their impact on repository size, showing both absolute and percentage values.
5. **Largest files**: Identification of the largest files in your repository by compressed size, along with their last commit year.
6. **File extensions**: Analysis of file extensions and their contribution to repository size.
7. **Contributors**: Statistics on authors and committers over time, showing who has contributed the most commits by year.

### Important metrics explained

- **Commits, Trees, Blobs**: These columns show the cumulative count of Git objects. Commits represent saved changes, trees represent folder snapshots, and blobs represent file versions.
- **On-disk size**: Shows the compressed size of Git objects as stored in Git's database (`.git/objects`). Objects are often stored using delta compression (storing only changes between similar objects). 
- **Percentages (%)**: In the growth table, percentages show estimated yearly growth relative to current totals. In directory and file listings, percentages show the proportion of total repository objects or size.
- **Growth projections**: Future estimates (marked with `*`) are calculated based on growth patterns from the last five years.
- **Directory markers**: Files or directories marked with `*` are not present in the latest commit (they were moved, renamed, or removed).

## Use cases

- Track repository growth over time
- Identify large files that may impact clone and fetch times
- Project future storage requirements for Git repositories
- Optimize repository size by identifying problematic files

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE.md)

## Building from source

If you prefer to build `git-metrics` from source, follow these steps:

### Prerequisites
- Git
- Go 1.23.2 or newer

1. Clone the repository:
   ```bash
   git clone https://github.com/steffen/git-metrics.git
   cd git-metrics
   ```

2. Build the binary:
   ```bash
   go build
   ```

After building, you can run the tool as described in the "Running the Tool" section.
