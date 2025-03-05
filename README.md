# git-metrics

A powerful Git repository analysis tool that provides detailed metrics, growth statistics, and future projections for your Git repositories.

## Overview

`git-metrics` is a command-line utility that analyzes Git repositories to provide comprehensive insights about repository growth, object statistics, and file usage. The tool gathers historical data and provides projections for future repository growth.

Key features include:
- Repository metadata analysis (first commit, age)
- Year-by-year growth statistics for Git objects (commits, trees, blobs) and disk usage
- Identification of largest files in the repository
- File extension distribution analysis
- Future growth projections based on historical trends

## Installation

### Prerequisites
- Git

### Download Prebuilt Binaries

The easiest way to install `git-metrics` is to download a prebuilt binary from the [GitHub releases page](https://github.com/steffen/git-metrics/releases).

#### Linux
```bash
# Download the latest release for Linux (64-bit)
curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-linux-amd64.tar.gz -o git-metrics.tar.gz

# Unpack the archive
tar -xzf git-metrics.tar.gz

# Move it to a directory in your PATH (optional)
sudo mv git-metrics /usr/local/bin/
```

#### macOS
```bash
# Download the latest release for macOS (Intel or Apple Silicon)
curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-darwin-arm64.zip -o git-metrics.zip

# Unzip the archive
unzip git-metrics.zip

# Remove quarantine attribute (required for macOS security)
xattr -d com.apple.quarantine git-metrics

# Move it to a directory in your PATH (optional)
sudo mv git-metrics /usr/local/bin/
```

### Running the Application

```bash
# Analyze the current directory as a Git repository
git-metrics

# Analyze a specific repository
git-metrics -r /path/to/repository
```

## Command Line Options

| Option | Description |
|--------|-------------|
| `-r`, `--repository` | Path to Git repository (default: current directory) |
| `--debug` | Enable debug output |
| `--no-progress` | Disable progress indicators |

## Output Examples

`git-metrics` provides detailed output about your repository:

```

RUN ############################################################################################

Start time                 Mon, 24 Feb 2025 12:06 CET
Machine                    10 CPU cores with 64 GB memory (macOS 15.3.1 on Apple M1 Max)
Git version                2.46.0

REPOSITORY #####################################################################################

Path                       /Users/steffen/GitHub/oss/git
Remote                     https://github.com/git/git.git
Most recent fetch          Wed, 05 Feb 2025 13:10 CET
Most recent commit         Mon, 03 Feb 2025 (bc204b7427)
First commit               Thu, 07 Apr 2005 (e83c51)
Age                        19 years 10 months 17 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005          3,215    +4 %          4,056    +3 %          5,922    +4 %         3.9 MB    +2 %
2006          7,816    +6 %         10,459    +4 %         13,181    +5 %         8.0 MB    +2 %
2007         13,312    +7 %         19,425    +6 %         22,503    +6 %        14.5 MB    +3 %
2008         17,432    +5 %         26,640    +5 %         30,747    +5 %        20.0 MB    +2 %
2009         21,267    +5 %         33,673    +5 %         37,227    +4 %        25.3 MB    +2 %
2010         25,150    +5 %         41,060    +5 %         44,099    +5 %        31.1 MB    +2 %
2011         28,673    +4 %         48,136    +5 %         50,231    +4 %        36.4 MB    +2 %
2012         32,454    +5 %         55,936    +5 %         55,807    +4 %        43.3 MB    +3 %
2013         36,772    +5 %         65,037    +6 %         62,775    +5 %        51.3 MB    +3 %
2014         39,875    +4 %         71,108    +4 %         68,121    +4 %        59.7 MB    +3 %
2015         43,161    +4 %         77,534    +4 %         73,436    +4 %        68.6 MB    +3 %
2016         47,031    +5 %         85,302    +5 %         79,873    +4 %        76.8 MB    +3 %
2017         51,617    +6 %         94,468    +6 %         88,493    +6 %        89.9 MB    +5 %
2018         56,098    +6 %        103,740    +6 %         99,098    +7 %       110.0 MB    +8 %
2019         59,867    +5 %        111,620    +5 %        106,638    +5 %       130.5 MB    +8 %
2020         63,562    +5 %        119,536    +5 %        114,111    +5 %       151.0 MB    +8 %
2021         67,578    +5 %        128,086    +6 %        121,937    +5 %       178.5 MB   +11 %
2022         71,225    +5 %        136,302    +5 %        129,471    +5 %       206.0 MB   +11 %
2023         74,172    +4 %        142,980    +4 %        139,020    +6 %       228.7 MB    +9 %
2024         78,581    +6 %        152,873    +6 %        149,434    +7 %       254.0 MB   +10 %
------------------------------------------------------------------------------------------------
2025^        79,103    +1 %        154,046    +1 %        150,546    +1 %       257.6 MB    +1 %
------------------------------------------------------------------------------------------------
2025*        82,323    +5 %        161,123    +5 %        157,993    +6 %       278.7 MB   +10 %
2026*        86,065    +5 %        169,373    +5 %        166,552    +6 %       303.3 MB   +10 %
2027*        89,807    +5 %        177,623    +5 %        175,111    +6 %       328.0 MB   +10 %
2028*        93,549    +5 %        185,873    +5 %        183,670    +6 %       352.7 MB   +10 %
2029*        97,291    +5 %        194,123    +5 %        192,229    +6 %       377.4 MB   +10 %
2030*       101,033    +5 %        202,373    +5 %        200,788    +6 %       402.1 MB   +10 %
------------------------------------------------------------------------------------------------

^ Current totals as of the last fetch on Wed, 05 Feb
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
po/fr.po                                      2024            139   0.1 %         4.8 MB   3.0 %
po/zh_CN.po                                   2025            153   0.1 %         4.6 MB   2.8 %
po/de.po                                      2025            184   0.1 %         4.4 MB   2.7 %
po/sv.po                                      2025            118   0.1 %         4.3 MB   2.7 %
whats-cooking.txt                             0001          1,328   0.9 %         4.2 MB   2.6 %
po/ca.po                                      2024             78   0.1 %         4.1 MB   2.5 %
po/vi.po                                      2025            103   0.1 %         3.6 MB   2.2 %
sequencer.c                                   2024          1,024   0.7 %         3.6 MB   2.2 %
po/bg.po                                      2024             71   0.0 %         3.3 MB   2.1 %
po/tr.po                                      2025             40   0.0 %         2.8 MB   1.7 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                   3,238   2.2 %        39.6 MB  24.6 %
└─ Out of 6,088                                           150,546 100.0 %       161.4 MB 100.0 %

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                     941  15.5 %         66,561  44.2 %        63.0 MB  24.5 %
.po                                     61   1.0 %          1,349   0.9 %        44.3 MB  17.2 %
.txt                                 1,323  21.7 %         20,381  13.5 %        16.7 MB   6.5 %
.sh                                  1,513  24.9 %         29,859  19.8 %        15.4 MB   6.0 %
No Extension                           687  11.3 %          9,894   6.6 %         6.8 MB   2.6 %
.h                                     377   6.2 %         12,047   8.0 %         5.3 MB   2.1 %
.perl                                   57   0.9 %          3,025   2.0 %         2.2 MB   0.8 %
.pot                                     5   0.1 %            128   0.1 %         2.1 MB   0.8 %
.py                                     30   0.5 %            534   0.4 %         1.4 MB   0.5 %
.bash                                    2   0.0 %          1,113   0.7 %         1.2 MB   0.5 %
------------------------------------------------------------------------------------------------
├─ Top 10                            4,996  82.1 %        144,891  96.2 %       158.3 MB  98.1 %
└─ Out of 363                        6,088 100.0 %        150,545 100.0 %       161.4 MB 100.0 %

Finished in 11s with a memory footprint of 96.6 MB.
```

```
RUN ############################################################################################

Start time                 Mon, 24 Feb 2025 12:06 CET
Machine                    10 CPU cores with 64 GB memory (macOS 15.3.1 on Apple M1 Max)
Git version                2.46.0

REPOSITORY #####################################################################################

Path                       /Users/steffen/GitHub/oss/linux
Remote                     https://github.com/torvalds/linux.git
Most recent fetch          Wed, 20 Nov 2024 15:23 CET
Most recent commit         Tue, 19 Nov 2024 (bf9aa14fc523)
First commit               Sat, 16 Apr 2005 (1da177)
Age                        19 years 10 months 8 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005         15,862    +1 %         71,850    +1 %         63,135    +2 %       121.5 MB    +2 %
2006         45,307    +2 %        204,857    +2 %        147,863    +3 %       184.9 MB    +1 %
2007         75,872    +2 %        339,474    +2 %        234,445    +3 %       258.9 MB    +1 %
2008        126,732    +4 %        562,419    +4 %        351,963    +4 %       369.0 MB    +2 %
2009        179,267    +4 %        800,794    +4 %        474,389    +4 %       501.8 MB    +3 %
2010        228,892    +4 %      1,030,868    +4 %        596,418    +4 %       627.0 MB    +2 %
2011        284,001    +4 %      1,290,522    +4 %        730,818    +5 %       806.6 MB    +3 %
2012        348,960    +5 %      1,605,101    +5 %        882,525    +5 %      1006.2 MB    +4 %
2013        420,318    +5 %      1,938,462    +5 %      1,027,261    +5 %         1.1 GB    +3 %
2014        496,275    +6 %      2,296,602    +6 %      1,177,785    +5 %         1.3 GB    +3 %
2015        571,739    +6 %      2,654,422    +6 %      1,328,028    +5 %         1.5 GB    +4 %
2016        648,804    +6 %      3,025,928    +6 %      1,485,410    +6 %         1.8 GB    +5 %
2017        729,672    +6 %      3,432,929    +6 %      1,675,154    +7 %         2.1 GB    +7 %
2018        810,063    +6 %      3,825,649    +6 %      1,832,020    +6 %         2.4 GB    +5 %
2019        892,604    +6 %      4,243,151    +7 %      2,021,226    +7 %         2.7 GB    +6 %
2020        983,051    +7 %      4,693,646    +7 %      2,203,817    +6 %         3.1 GB    +9 %
2021      1,069,164    +7 %      5,117,001    +7 %      2,368,641    +6 %         3.7 GB   +10 %
2022      1,155,023    +7 %      5,544,603    +7 %      2,539,517    +6 %         4.2 GB   +11 %
2023      1,246,324    +7 %      5,991,045    +7 %      2,721,177    +6 %         4.7 GB    +9 %
2024      1,312,943    +5 %      6,317,253    +5 %      2,847,880    +4 %         5.1 GB    +7 %
------------------------------------------------------------------------------------------------
2025^     1,312,943    +0 %      6,317,253    +0 %      2,847,880    +0 %         5.1 GB    +0 %
------------------------------------------------------------------------------------------------
2025*     1,397,010    +6 %      6,732,073    +7 %      3,013,210    +6 %         5.5 GB    +9 %
2026*     1,481,077    +6 %      7,146,893    +7 %      3,178,540    +6 %         6.0 GB    +9 %
2027*     1,565,144    +6 %      7,561,713    +7 %      3,343,870    +6 %         6.5 GB    +9 %
2028*     1,649,211    +6 %      7,976,533    +7 %      3,509,200    +6 %         6.9 GB    +9 %
2029*     1,733,278    +6 %      8,391,353    +7 %      3,674,530    +6 %         7.4 GB    +9 %
2030*     1,817,345    +6 %      8,806,173    +7 %      3,839,860    +6 %         7.9 GB    +9 %
------------------------------------------------------------------------------------------------

^ Current totals as of the last fetch on Wed, 20 Nov
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
MAINTAINERS                                   2024         18,098   0.6 %       177.9 MB   4.9 %
kernel/bpf/verifier.c                         2024          1,378   0.0 %        13.6 MB   0.4 %
drivers/gpu/drm/i915/intel_display.c          2019          4,155   0.1 %         9.1 MB   0.2 %
drivers/gpu/drm/i915/i915_reg.h               2024          2,428   0.1 %         8.8 MB   0.2 %
drivers/gpu/drm/i915/display/intel_display.c  2024          1,334   0.0 %         8.8 MB   0.2 %
fs/io_uring.c                                 2022          2,054   0.1 %         7.9 MB   0.2 %
drivers/gpu/drm/amd/...mdgpu_dm/amdgpu_dm.c   2024          1,783   0.1 %         7.8 MB   0.2 %
arch/x86/kvm/x86.c                            2024          2,955   0.1 %         7.7 MB   0.2 %
kernel/sched/fair.c                           2024          1,500   0.1 %         6.7 MB   0.2 %
crypto/testmgr.h                              2024            220   0.0 %         6.6 MB   0.2 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                  35,905   1.3 %       255.0 MB   7.0 %
└─ Out of 151,479                                       2,847,880 100.0 %         3.6 GB 100.0 %

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                  58,426  38.6 %      1,852,539  65.0 %         2.6 GB  50.6 %
.h                                  50,446  33.3 %        584,094  20.5 %       532.6 MB  10.3 %
No Extension                        10,914   7.2 %        192,484   6.8 %       260.8 MB   5.0 %
.dtsi                                3,100   2.0 %         42,702   1.5 %        45.3 MB   0.9 %
.rst                                 5,166   3.4 %         24,800   0.9 %        42.5 MB   0.8 %
.txt                                 6,118   4.0 %         33,490   1.2 %        35.5 MB   0.7 %
.S                                   2,928   1.9 %         29,077   1.0 %        24.6 MB   0.5 %
.dts                                 4,265   2.8 %         33,416   1.2 %        23.0 MB   0.4 %
.yaml                                4,600   3.0 %         21,260   0.7 %        14.4 MB   0.3 %
.json                                  768   0.5 %          2,966   0.1 %         6.6 MB   0.1 %
------------------------------------------------------------------------------------------------
├─ Top 10                          146,731  96.9 %      2,816,828  98.9 %         3.5 GB  98.9 %
└─ Out of 435                      151,479 100.0 %      2,847,880 100.0 %         3.6 GB 100.0 %

Finished in 3m18s with a memory footprint of 2.1 GB.
```

## Understanding the Output

`git-metrics` provides several sections of output:

1. **Repository Information**: Basic metadata about your repository including path, remote URL, and commit history.

2. **Growth Statistics**: Year-by-year breakdown of Git object growth (commits, trees, blobs) and disk usage.

3. **Growth Projections**: Estimation of future repository growth based on historical trends.

4. **Largest Files**: Identification of the largest files in your repository by compressed size.

5. **File Extensions**: Analysis of file extensions and their impact on repository size.

## Use Cases

- Track repository growth over time
- Identify large files that may impact clone and fetch times
- Project future storage requirements for Git repositories
- Optimize repository size by identifying problematic files

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE.md)

## Building from Source

If you prefer to build `git-metrics` from source, follow these steps:

### Prerequisites
- Git
- Go 1.23.2 or newer

```bash
# Clone the repository
git clone https://github.com/steffen/git-metrics.git
cd git-metrics

# Build the binary
go build
```

After building, you can run the tool as described in the "Running the Application" section.
