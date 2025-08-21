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

## Output examples

### [`git/git`](https://github.com/git/git)

```
RUN ############################################################################################

Start time                 Thu, 21 Aug 2025 17:56 CEST
Machine                    10 CPU cores with 64 GB memory (macOS 15.6 on Apple M1 Max)
Git metrics version        1.3.0
Git version                2.46.0

REPOSITORY #####################################################################################

Git directory              /Users/steffen/GitHub/oss/git/.git
Remote                     https://github.com/git/git.git
Most recent fetch          Thu, 21 Aug 2025 17:56 CEST
Most recent commit         Sun, 17 Aug 2025 (c44beea485)
First commit               Thu, 07 Apr 2005 (e83c51)
Age                        20 years 4 months 14 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005          3,215    +4 %          4,056    +3 %          5,922    +4 %         4.1 MB    +1 %
2006          7,816    +6 %         10,459    +4 %         13,181    +5 %         8.4 MB    +2 %
2007         13,312    +7 %         19,425    +6 %         22,503    +6 %        15.2 MB    +2 %
2008         17,440    +5 %         26,655    +5 %         30,759    +5 %        21.0 MB    +2 %
2009         21,267    +5 %         33,673    +4 %         37,227    +4 %        26.6 MB    +2 %
2010         25,150    +5 %         41,060    +5 %         44,099    +4 %        32.6 MB    +2 %
2011         28,673    +4 %         48,136    +4 %         50,231    +4 %        38.2 MB    +2 %
2012         32,455    +5 %         55,937    +5 %         55,808    +4 %        45.4 MB    +3 %
2013         36,772    +5 %         65,037    +6 %         62,775    +4 %        53.8 MB    +3 %
2014         39,875    +4 %         71,108    +4 %         68,121    +3 %        62.6 MB    +3 %
2015         43,161    +4 %         77,534    +4 %         73,436    +3 %        71.9 MB    +3 %
2016         47,031    +5 %         85,302    +5 %         79,873    +4 %        80.5 MB    +3 %
2017         51,617    +6 %         94,468    +6 %         88,493    +6 %        94.3 MB    +5 %
2018         56,098    +6 %        103,740    +6 %         99,098    +7 %       115.3 MB    +7 %
2019         59,867    +5 %        111,620    +5 %        106,638    +5 %       136.9 MB    +8 %
2020         63,562    +5 %        119,536    +5 %        114,111    +5 %       158.3 MB    +7 %
2021         67,578    +5 %        128,086    +5 %        121,937    +5 %       187.2 MB   +10 %
2022         71,225    +4 %        136,302    +5 %        129,471    +5 %       216.0 MB   +10 %
2023         74,172    +4 %        142,980    +4 %        139,020    +6 %       239.8 MB    +8 %
2024         78,394    +5 %        152,556    +6 %        149,393    +7 %       266.0 MB    +9 %
------------------------------------------------------------------------------------------------
2025^        81,072    +3 %        158,302    +4 %        155,623    +4 %       286.1 MB    +7 %
------------------------------------------------------------------------------------------------
2025*        82,099    +5 %        160,743    +5 %        157,944    +5 %       291.9 MB    +9 %
2026*        85,804    +5 %        168,930    +5 %        166,495    +5 %       317.7 MB    +9 %
2027*        89,509    +5 %        177,117    +5 %        175,046    +5 %       343.5 MB    +9 %
2028*        93,214    +5 %        185,304    +5 %        183,597    +5 %       369.4 MB    +9 %
2029*        96,919    +5 %        193,491    +5 %        192,148    +5 %       395.2 MB    +9 %
2030*       100,624    +5 %        201,678    +5 %        200,699    +5 %       421.0 MB    +9 %
------------------------------------------------------------------------------------------------

^ Current totals as of the most recent fetch on Thu, 21 Aug
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

RATE OF CHANGES ################################################################################

Commits to default branch (master)

              Commits             Peak per day            Peak per hour          Peak per minute
Year         per year        P95    P99   P100        P95    P99   P100        P95    P99   P100
------------------------------------------------------------------------------------------------
2005            3,132   │     25     33     87   │      7     13     64   │      4      7     17
2006            4,499   │     32     40     44   │      9     15     37   │      3      7     34
2007            5,388   │     32     55    132   │      8     16    113   │      3      7    113
2008            4,009   │     27     34     46   │      9     16     29   │      4      8     21
2009            3,676   │     30     40     58   │     10     17     27   │      4      8     18
2010            3,661   │     36     51     93   │     13     23     62   │      5     13     32
2011            3,371   │     31     59     77   │     12     22     72   │      4     12     71
2012            3,605   │     31     47     58   │     12     19     54   │      5     10     54
2013            4,180   │     42     54     82   │     14     23     61   │      5     12     45
2014            3,006   │     36     52     71   │     16     29     57   │      6     15     34
2015            3,176   │     41     54     68   │     16     30     56   │      8     20     56
2016            3,745   │     37     54     62   │     18     32     46   │      9     22     41
2017            4,469   │     46     56     78   │     21     34     50   │     12     23     43
2018            4,421   │     46     62     85   │     22     37     82   │     12     30     77
2019            3,700   │     45     70     97   │     19     34     66   │     12     24     44
2020            3,600   │     33     52     70   │     14     27     50   │     10     20     41
2021            3,911   │     39     51     82   │     16     28     78   │     11     21     76
2022            3,520   │     35     47     64   │     15     25     55   │     11     18     34
2023            2,833   │     32     47     60   │     15     27     55   │     10     22     54
2024            3,973   │     34     48     54   │     18     27     51   │     11     22     49
2025            2,186   │     38     50     66   │     14     21     66   │     10     16     35

LARGEST DIRECTORIES ############################################################################

Showing directories and files that contribute more than 1% of total on-disk size.

Path                                                        Blobs           On-disk size
------------------------------------------------------------------------------------------------
./                                                        155,623 100.0 %       181.5 MB 100.0 %
├─ po/                                                      1,412   0.9 %        51.2 MB  28.2 %
│  ├─ fr.po                                                   142   0.1 %         5.3 MB   2.9 %
│  ├─ zh_CN.po                                                156   0.1 %         5.2 MB   2.8 %
│  ├─ de.po                                                   186   0.1 %         4.9 MB   2.7 %
│  ├─ sv.po                                                   121   0.1 %         4.7 MB   2.6 %
│  ├─ ca.po                                                    79   0.1 %         4.5 MB   2.5 %
│  ├─ vi.po                                                   105   0.1 %         4.1 MB   2.2 %
│  ├─ bg.po                                                    76   0.0 %         3.9 MB   2.1 %
│  ├─ tr.po                                                    43   0.0 %         3.2 MB   1.8 %
│  ├─ zh_TW.po                                                 36   0.0 %         2.9 MB   1.6 %
│  ├─ git.pot*                                                108   0.1 %         2.1 MB   1.2 %
│  ├─ pt_PT.po                                                 38   0.0 %         2.0 MB   1.1 %
│  └─ id.po                                                    31   0.0 %         1.9 MB   1.0 %
├─ builtin/                                                18,974  12.2 %        22.0 MB  12.1 %
│  └─ pack-objects.c                                          653   0.4 %         2.1 MB   1.1 %
├─ t/                                                      29,580  19.0 %        17.3 MB   9.5 %
├─ Documentation/                                          19,976  12.8 %        14.3 MB   7.9 %
│  ├─ RelNotes/                                             1,850   1.2 %         3.0 MB   1.7 %
│  └─ config.txt*                                           1,459   0.9 %         1.8 MB   1.0 %
├─ contrib/                                                 3,780   2.4 %         3.3 MB   1.8 %
├─ compat/                                                  1,241   0.8 %         2.2 MB   1.2 %
├─ refs/                                                    1,196   0.8 %         2.1 MB   1.2 %
├─ whats-cooking.txt*                                       1,387   0.9 %         4.6 MB   2.5 %
├─ sequencer.c                                              1,058   0.7 %         4.1 MB   2.3 %
├─ Makefile                                                 3,298   2.1 %         3.0 MB   1.7 %
├─ diff.c                                                   1,627   1.0 %         2.4 MB   1.3 %
└─ read-cache.c                                               842   0.5 %         2.2 MB   1.2 %

* File or directory not present in latest commit of master branch (moved, renamed or removed)

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
po/fr.po                                      2025            142   0.1 %         5.3 MB   2.9 %
po/zh_CN.po                                   2025            156   0.1 %         5.2 MB   2.8 %
po/de.po                                      2025            186   0.1 %         4.9 MB   2.7 %
po/sv.po                                      2025            121   0.1 %         4.7 MB   2.6 %
whats-cooking.txt                             0001          1,387   0.9 %         4.6 MB   2.5 %
po/ca.po                                      2025             79   0.1 %         4.5 MB   2.5 %
sequencer.c                                   2025          1,058   0.7 %         4.1 MB   2.3 %
po/vi.po                                      2025            105   0.1 %         4.1 MB   2.2 %
po/bg.po                                      2025             76   0.0 %         3.9 MB   2.1 %
po/tr.po                                      2025             43   0.0 %         3.2 MB   1.8 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                   3,353   2.2 %        44.4 MB  24.5 %
└─ Out of 6,287                                           155,623 100.0 %       181.5 MB 100.0 %

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                     967  15.4 %         69,150  44.4 %        71.6 MB  25.0 %
.po                                     63   1.0 %          1,387   0.9 %        49.8 MB  17.4 %
.txt                                 1,323  21.0 %         20,482  13.2 %        17.8 MB   6.2 %
.sh                                  1,538  24.5 %         30,562  19.6 %        17.0 MB   5.9 %
No Extension                           704  11.2 %         10,183   6.5 %         7.7 MB   2.7 %
.h                                     382   6.1 %         12,481   8.0 %         5.9 MB   2.1 %
.perl                                   60   1.0 %          3,051   2.0 %         2.4 MB   0.8 %
.pot                                     5   0.1 %            128   0.1 %         2.2 MB   0.8 %
.py                                     30   0.5 %            534   0.3 %         1.4 MB   0.5 %
.bash                                    2   0.0 %          1,115   0.7 %         1.3 MB   0.5 %
------------------------------------------------------------------------------------------------
├─ Top 10                            5,074  80.7 %        149,073  95.8 %       177.1 MB  97.6 %
└─ Out of 364                        6,287 100.0 %        155,622 100.0 %       181.5 MB 100.0 %

AUTHORS WITH MOST COMMITS ######################################################################

Year     Author (#1)    Commits        Author (#2)    Commits        Author (#3)    Commits
------------------------------------------------------------------------------------------------
2005   │ Junio C Hamano   1,368  43% │ Linus Torvalds     680  21% │ Kay Sievers        139   4%
2006   │ Junio C Hamano   2,202  48% │ Shawn O. Pe...     243   5% │ Jakub Narebski     220   5%
2007   │ Junio C Hamano   1,606  29% │ Shawn O. Pe...     841  15% │ Johannes Sc...     244   4%
2008   │ Junio C Hamano   1,365  33% │ Shawn O. Pe...     163   4% │ Jeff King          152   4%
2009   │ Junio C Hamano   1,439  38% │ Jeff King          134   3% │ Johannes Sc...     104   3%
2010   │ Junio C Hamano   1,537  40% │ Jonathan Ni...     333   9% │ Ævar Arnfjö...     121   3%
2011   │ Junio C Hamano   1,588  45% │ Jeff King          237   7% │ Jonathan Ni...     190   5%
2012   │ Junio C Hamano   1,743  46% │ Jeff King          314   8% │ Nguyễn Thái...     204   5%
2013   │ Junio C Hamano   1,787  41% │ Felipe Cont...     283   7% │ Jeff King          243   6%
2014   │ Junio C Hamano   1,220  39% │ Jeff King          340  11% │ Nguyễn Thái...     140   5%
2015   │ Junio C Hamano   1,440  44% │ Jeff King          368  11% │ Michael Hag...     173   5%
2016   │ Junio C Hamano   1,604  41% │ Jeff King          386  10% │ Johannes Sc...     193   5%
2017   │ Junio C Hamano   1,739  38% │ Jeff King          405   9% │ Brandon Wil...     217   5%
2018   │ Junio C Hamano   1,230  27% │ Nguyễn Thái...     499  11% │ Jeff King          255   6%
2019   │ Junio C Hamano   1,095  29% │ Johannes Sc...     321   9% │ Jeff King          271   7%
2020   │ Junio C Hamano   1,157  31% │ Jeff King          270   7% │ Johannes Sc...     208   6%
2021   │ Junio C Hamano   1,215  30% │ Ævar Arnfjö...     606  15% │ Elijah Newren      201   5%
2022   │ Junio C Hamano   1,187  33% │ Ævar Arnfjö...     555  15% │ Taylor Blau        181   5%
2023   │ Junio C Hamano   1,041  35% │ Jeff King          330  11% │ Elijah Newren      195   7%
2024   │ Junio C Hamano   1,449  34% │ Patrick Ste...   1,026  24% │ Jeff King          259   6%
2025   │ Junio C Hamano     974  36% │ Patrick Ste...     434  16% │ Jeff King           95   4%
------------------------------------------------------------------------------------------------
Total  │ Junio C Hamano  29,986  37% │ Jeff King        4,598   6% │ Johannes Sc...   2,380   3%

COMMITTERS WITH MOST COMMITS ###################################################################

Year     Committer (#1) Commits        Committer (#2) Commits        Committer (#3) Commits
------------------------------------------------------------------------------------------------
2005   │ Junio C Hamano   1,794  56% │ Linus Torvalds   1,036  32% │ Kay Sievers        136   4%
2006   │ Junio C Hamano   4,310  94% │ Shawn O. Pe...     158   3% │ Paul Mackerras      82   2%
2007   │ Junio C Hamano   3,964  72% │ Shawn O. Pe...     842  15% │ Simon Hausmann     247   4%
2008   │ Junio C Hamano   3,532  86% │ Shawn O. Pe...     372   9% │ Paul Mackerras     123   3%
2009   │ Junio C Hamano   3,514  92% │ Eric Wong           96   3% │ Avery Pennarun      65   2%
2010   │ Junio C Hamano   3,665  94% │ Eric Wong           49   1% │ Pat Thoyts          47   1%
2011   │ Junio C Hamano   3,328  95% │ Jonathan Ni...      80   2% │ Pat Thoyts          70   2%
2012   │ Junio C Hamano   3,348  89% │ Jeff King          153   4% │ Jiang Xin           91   2%
2013   │ Junio C Hamano   4,093  95% │ Jonathan Ni...      81   2% │ Jiang Xin           24   1%
2014   │ Junio C Hamano   2,967  96% │ Jiang Xin           27   1% │ Eric Wong           23   1%
2015   │ Junio C Hamano   2,970  90% │ Jeff King          111   3% │ Jiang Xin           63   2%
2016   │ Junio C Hamano   3,652  94% │ Jiang Xin           45   1% │ Michael Hag...      34   1%
2017   │ Junio C Hamano   4,444  97% │ Jiang Xin           64   1% │ Jean-Noel A...      15   0%
2018   │ Junio C Hamano   4,327  97% │ Jiang Xin           64   1% │ Jeff King           21   0%
2019   │ Junio C Hamano   3,540  94% │ Johannes Sc...      67   2% │ Jiang Xin           53   1%
2020   │ Junio C Hamano   3,474  94% │ Jiang Xin           62   2% │ Pratyush Yadav      41   1%
2021   │ Junio C Hamano   3,796  95% │ Jiang Xin           81   2% │ Johannes Sc...      30   1%
2022   │ Junio C Hamano   3,239  89% │ Taylor Blau        240   7% │ Jiang Xin           58   2%
2023   │ Junio C Hamano   2,765  94% │ Johannes Sc...      85   3% │ Jiang Xin           36   1%
2024   │ Junio C Hamano   3,792  90% │ Taylor Blau        218   5% │ Johannes Sc...      86   2%
2025   │ Junio C Hamano   2,400  90% │ Taylor Blau         66   2% │ Mark Levedahl       58   2%
------------------------------------------------------------------------------------------------
Total  │ Junio C Hamano  72,914  90% │ Shawn O. Pe...   1,457   2% │ Linus Torvalds   1,041   1%

Finished in 14s with a memory footprint of 118.7 MB.
```

### [`torvalds/linux`](https://github.com/torvalds/linux)

```
RUN ############################################################################################

Start time                 Thu, 21 Aug 2025 18:03 CEST
Machine                    10 CPU cores with 64 GB memory (macOS 15.6 on Apple M1 Max)
Git metrics version        1.3.0
Git version                2.46.0

REPOSITORY #####################################################################################

Git directory              /Users/steffen/GitHub/oss/linux/.git
Remote                     https://github.com/torvalds/linux.git
Most recent fetch          Thu, 21 Aug 2025 18:01 CEST
Most recent commit         Thu, 21 Aug 2025 (1c656b1efde6)
First commit               Sat, 16 Apr 2005 (1da177)
Age                        20 years 4 months 5 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005         15,862    +1 %         71,850    +1 %         63,135    +2 %       127.4 MB    +2 %
2006         45,307    +2 %        204,857    +2 %        147,863    +3 %       193.9 MB    +1 %
2007         75,872    +2 %        339,474    +2 %        234,445    +3 %       271.5 MB    +1 %
2008        126,734    +4 %        562,425    +3 %        351,965    +4 %       386.9 MB    +2 %
2009        179,269    +4 %        800,802    +4 %        474,393    +4 %       526.2 MB    +2 %
2010        228,892    +4 %      1,030,868    +3 %        596,418    +4 %       657.5 MB    +2 %
2011        284,002    +4 %      1,290,525    +4 %        730,819    +4 %       845.8 MB    +3 %
2012        348,964    +5 %      1,605,119    +5 %        882,529    +5 %         1.1 GB    +4 %
2013        420,326    +5 %      1,938,504    +5 %      1,027,272    +5 %         1.2 GB    +3 %
2014        496,286    +5 %      2,296,647    +5 %      1,177,798    +5 %         1.4 GB    +3 %
2015        571,740    +5 %      2,654,428    +5 %      1,328,031    +5 %         1.6 GB    +4 %
2016        648,805    +6 %      3,025,934    +6 %      1,485,413    +5 %         1.9 GB    +4 %
2017        729,675    +6 %      3,432,943    +6 %      1,675,158    +6 %         2.3 GB    +7 %
2018        810,065    +6 %      3,825,660    +6 %      1,832,024    +5 %         2.6 GB    +5 %
2019        892,609    +6 %      4,243,177    +6 %      2,021,235    +6 %         2.9 GB    +6 %
2020        983,052    +7 %      4,693,652    +7 %      2,203,820    +6 %         3.4 GB    +8 %
2021      1,069,168    +6 %      5,117,022    +6 %      2,368,647    +6 %         3.9 GB   +10 %
2022      1,155,031    +6 %      5,544,636    +6 %      2,539,526    +6 %         4.5 GB   +11 %
2023      1,246,336    +7 %      5,991,105    +7 %      2,721,194    +6 %         5.1 GB    +9 %
2024      1,329,852    +6 %      6,403,293    +6 %      2,885,958    +6 %         5.5 GB    +8 %
------------------------------------------------------------------------------------------------
2025^     1,381,903    +4 %      6,661,078    +4 %      2,989,020    +3 %         5.8 GB    +4 %
------------------------------------------------------------------------------------------------
2025*     1,417,300    +6 %      6,835,316    +6 %      3,058,902    +6 %         6.1 GB    +9 %
2026*     1,504,748    +6 %      7,267,339    +6 %      3,231,846    +6 %         6.6 GB    +9 %
2027*     1,592,196    +6 %      7,699,362    +6 %      3,404,790    +6 %         7.1 GB    +9 %
2028*     1,679,644    +6 %      8,131,385    +6 %      3,577,734    +6 %         7.6 GB    +9 %
2029*     1,767,092    +6 %      8,563,408    +6 %      3,750,678    +6 %         8.2 GB    +9 %
2030*     1,854,540    +6 %      8,995,431    +6 %      3,923,622    +6 %         8.7 GB    +9 %
------------------------------------------------------------------------------------------------

^ Current totals as of the most recent fetch on Thu, 21 Aug
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

RATE OF CHANGES ################################################################################

Commits to default branch (master)

              Commits             Peak per day            Peak per hour          Peak per minute
Year         per year        P95    P99   P100        P95    P99   P100        P95    P99   P100
------------------------------------------------------------------------------------------------
2005           15,860   │    226    390    612   │     24     87    347   │      5     27    290
2006           29,446   │    321    559    733   │     35    129    412   │     12     48    354
2007           30,564   │    346    634   1218   │     34    118    836   │     12     50    417
2008           50,836   │    408   1029   1381   │     40    111    926   │     10     46    393
2009           52,613   │    401    758   1155   │     30     89    715   │      8     31    456
2010           49,636   │    334    635    957   │     33     75    461   │      8     28    281
2011           55,112   │    333    555    642   │     35     79    437   │      8     25    163
2012           64,959   │    376    519    560   │     38     77    348   │      7     23    181
2013           71,348   │    403    559    689   │     41     83    453   │      8     23    369
2014           75,933   │    452    534    870   │     42     91    303   │      8     24    253
2015           75,440   │    398    530    787   │     39     88    405   │      8     22    404
2016           77,138   │    425    536    665   │     40     84    288   │      8     20    286
2017           80,858   │    429    585    735   │     42     79    510   │      8     21    153
2018           80,413   │    435    599    664   │     42     81    373   │      8     21    167
2019           82,544   │    463    562    753   │     42     74    338   │      8     20    217
2020           90,450   │    485    624    735   │     45     86    366   │      9     22    363
2021           86,112   │    498    622    911   │     46     85    293   │      9     24    257
2022           85,823   │    496    617    870   │     44     90    345   │      9     24    249
2023           91,340   │    520    731   2812   │     44     90   2765   │      9     24   1298
2024           83,520   │    465    646    784   │     44     80    285   │      9     21    211
2025           52,052   │    487    628    735   │     42     84    236   │      9     23    196

LARGEST DIRECTORIES ############################################################################

Showing directories and files that contribute more than 1% of total on-disk size.

Path                                                        Blobs           On-disk size
------------------------------------------------------------------------------------------------
./                                                      2,989,020 100.0 %         4.0 GB 100.0 %
├─ drivers/                                             1,425,418  47.7 %         2.1 GB  52.0 %
│  ├─ net/                                                295,444   9.9 %       503.3 MB  12.5 %
│  │  ├─ ethernet/                                        115,093   3.9 %       246.1 MB   6.1 %
│  │  │  └─ intel/                                         21,564   0.7 %        52.7 MB   1.3 %
│  │  └─ wireless/                                        110,534   3.7 %       163.1 MB   4.0 %
│  │     └─ ath/                                           23,619   0.8 %        40.4 MB   1.0 %
│  ├─ gpu/                                                267,723   9.0 %       455.8 MB  11.3 %
│  │  └─ drm/                                             266,291   8.9 %       454.4 MB  11.3 %
│  │     ├─ amd/                                           74,200   2.5 %       205.5 MB   5.1 %
│  │     │  ├─ include/                                     1,610   0.1 %        71.1 MB   1.8 %
│  │     │  │  └─ asic_reg/                                   880   0.0 %        68.4 MB   1.7 %
│  │     │  ├─ display/                                    27,936   0.9 %        57.8 MB   1.4 %
│  │     │  │  └─ dc/                                      22,574   0.8 %        44.0 MB   1.1 %
│  │     │  └─ amdgpu/                                     31,619   1.1 %        51.3 MB   1.3 %
│  │     └─ i915/                                          74,362   2.5 %       123.5 MB   3.1 %
│  ├─ staging/                                            152,990   5.1 %       139.4 MB   3.5 %
│  ├─ scsi/                                                53,040   1.8 %       123.9 MB   3.1 %
│  ├─ media/                                               81,795   2.7 %       103.0 MB   2.6 %
│  ├─ usb/                                                 46,867   1.6 %        66.6 MB   1.7 %
│  └─ infiniband/                                          31,627   1.1 %        55.5 MB   1.4 %
├─ arch/                                                  525,473  17.6 %       436.4 MB  10.8 %
│  ├─ arm/                                                153,262   5.1 %        95.1 MB   2.4 %
│  ├─ x86/                                                 88,921   3.0 %        93.1 MB   2.3 %
│  ├─ arm64/                                               57,518   1.9 %        62.5 MB   1.5 %
│  └─ powerpc/                                             62,050   2.1 %        58.1 MB   1.4 %
├─ fs/                                                    211,763   7.1 %       287.1 MB   7.1 %
│  └─ btrfs/                                               29,499   1.0 %        51.3 MB   1.3 %
├─ net/                                                   162,838   5.4 %       214.0 MB   5.3 %
├─ include/                                               194,541   6.5 %       168.2 MB   4.2 %
│  └─ linux/                                              102,143   3.4 %        89.1 MB   2.2 %
├─ sound/                                                  91,016   3.0 %       127.4 MB   3.2 %
│  └─ soc/                                                 54,408   1.8 %        74.2 MB   1.8 %
├─ kernel/                                                 63,862   2.1 %       120.5 MB   3.0 %
├─ Documentation/                                          93,766   3.1 %       112.5 MB   2.8 %
├─ tools/                                                  92,830   3.1 %       107.0 MB   2.7 %
│  └─ perf/                                                43,818   1.5 %        46.1 MB   1.1 %
├─ mm/                                                     34,450   1.2 %        60.6 MB   1.5 %
└─ MAINTAINERS                                             19,611   0.7 %       197.0 MB   4.9 %

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
MAINTAINERS                                   2025         19,611   0.7 %       197.0 MB   4.9 %
kernel/bpf/verifier.c                         2025          1,535   0.1 %        15.6 MB   0.4 %
drivers/gpu/drm/i9...ay/intel_display.c [1]   2025          1,532   0.1 %        10.0 MB   0.2 %
drivers/gpu/drm/i915/intel_display.c          2019          4,155   0.1 %         9.5 MB   0.2 %
drivers/gpu/drm/i915/i915_reg.h               2025          2,482   0.1 %         9.4 MB   0.2 %
drivers/gpu/drm/am...gpu_dm/amdgpu_dm.c [2]   2025          1,970   0.1 %         8.8 MB   0.2 %
arch/x86/kvm/x86.c                            2025          3,136   0.1 %         8.7 MB   0.2 %
fs/io_uring.c                                 2022          2,054   0.1 %         8.3 MB   0.2 %
crypto/testmgr.h                              2025            228   0.0 %         7.5 MB   0.2 %
fs/btrfs/inode.c                              2025          2,860   0.1 %         7.4 MB   0.2 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                  39,563   1.3 %       282.2 MB   7.0 %
└─ Out of 157,671                                       2,989,020 100.0 %         4.0 GB 100.0 %

[1] drivers/gpu/drm/i915/display/intel_display.c
[2] drivers/gpu/drm/amd/display/amdgpu_dm/amdgpu_dm.c

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                  60,398  38.3 %      1,941,028  64.9 %         2.9 GB  50.4 %
.h                                  51,823  32.9 %        611,450  20.5 %       584.6 MB  10.2 %
No Extension                        11,261   7.1 %        201,477   6.7 %       287.5 MB   5.0 %
.dtsi                                3,410   2.2 %         45,715   1.5 %        52.6 MB   0.9 %
.rst                                 5,408   3.4 %         26,781   0.9 %        47.8 MB   0.8 %
.txt                                 6,142   3.9 %         33,853   1.1 %        37.8 MB   0.7 %
.dts                                 4,713   3.0 %         35,776   1.2 %        26.5 MB   0.5 %
.S                                   2,993   1.9 %         29,531   1.0 %        26.1 MB   0.5 %
.yaml                                5,303   3.4 %         24,483   0.8 %        16.9 MB   0.3 %
.json                                  880   0.6 %          3,428   0.1 %         7.7 MB   0.1 %
------------------------------------------------------------------------------------------------
├─ Top 10                          152,331  96.6 %      2,953,522  98.8 %         4.0 GB  98.9 %
└─ Out of 449                      157,671 100.0 %      2,989,020 100.0 %         4.0 GB 100.0 %

AUTHORS WITH MOST COMMITS ######################################################################

Year     Author (#1)    Commits        Author (#2)    Commits        Author (#3)    Commits
------------------------------------------------------------------------------------------------
2005   │ Linus Torvalds     775   5% │ Jeff Garzik        392   2% │ Russell King       344   2%
2006   │ Linus Torvalds   1,108   4% │ Al Viro            765   3% │ David S. Mi...     612   2%
2007   │ Linus Torvalds   1,394   5% │ Ralf Baechle       506   2% │ Thomas Glei...     484   2%
2008   │ Linus Torvalds   1,912   4% │ Ingo Molnar      1,271   2% │ David S. Mi...     928   2%
2009   │ Linus Torvalds   2,124   4% │ Ingo Molnar      1,088   2% │ Takashi Iwai       952   2%
2010   │ Linus Torvalds   1,884   4% │ Joe Perches        546   1% │ Chris Wilson       519   1%
2011   │ Linus Torvalds   2,080   4% │ Mark Brown       1,047   2% │ David S. Mi...     743   1%
2012   │ Linus Torvalds   2,271   3% │ H Hartley S...   1,447   2% │ Mark Brown       1,224   2%
2013   │ Linus Torvalds   2,044   3% │ H Hartley S...   1,582   2% │ Mark Brown       1,506   2%
2014   │ Linus Torvalds   2,085   3% │ H Hartley S...   1,620   2% │ David S. Mi...     922   1%
2015   │ Linus Torvalds   2,009   3% │ David S. Mi...     987   1% │ H Hartley S...     784   1%
2016   │ Linus Torvalds   2,270   3% │ Arnd Bergmann    1,186   2% │ David S. Mi...   1,151   1%
2017   │ Linus Torvalds   2,303   3% │ David S. Mi...   1,420   2% │ Arnd Bergmann    1,121   1%
2018   │ Linus Torvalds   2,172   3% │ David S. Mi...   1,405   2% │ Arnd Bergmann      904   1%
2019   │ Linus Torvalds   2,384   3% │ David S. Mi...   1,207   1% │ Chris Wilson     1,174   1%
2020   │ Linus Torvalds   2,635   3% │ Mauro Carva...   1,215   1% │ Christoph H...   1,200   1%
2021   │ Linus Torvalds   2,415   3% │ David S. Mi...   1,042   1% │ Arnd Bergmann      978   1%
2022   │ Linus Torvalds   2,601   3% │ Krzysztof K...   1,371   2% │ Jakub Kicinski     981   1%
2023   │ Kent Overst...   2,813   3% │ Linus Torvalds   2,505   3% │ Uwe Kleine-...   2,235   2%
2024   │ Linus Torvalds   2,889   3% │ Kent Overst...   1,368   2% │ Krzysztof K...   1,180   1%
2025   │ Linus Torvalds   1,805   3% │ Jakub Kicinski   1,024   2% │ Kent Overst...     698   1%
------------------------------------------------------------------------------------------------
Total  │ Linus Torvalds  43,665   3% │ David S. Mi...  15,662   1% │ Arnd Bergmann   11,963   1%

COMMITTERS WITH MOST COMMITS ###################################################################

Year     Committer (#1) Commits        Committer (#2) Commits        Committer (#3) Commits
------------------------------------------------------------------------------------------------
2005   │ Linus Torvalds   6,398  40% │ David S. Mi...   1,384   9% │ Jeff Garzik      1,384   9%
2006   │ Linus Torvalds   9,384  32% │ David S. Mi...   3,115  11% │ Greg Kroah-...   1,842   6%
2007   │ Linus Torvalds   7,419  24% │ David S. Mi...   3,337  11% │ Jeff Garzik      2,040   7%
2008   │ Linus Torvalds   7,305  14% │ Ingo Molnar      5,983  12% │ David S. Mi...   5,730  11%
2009   │ Linus Torvalds   5,731  11% │ David S. Mi...   5,367  10% │ Ingo Molnar      4,362   8%
2010   │ David S. Mi...   5,073  10% │ Greg Kroah-...   4,497   9% │ Linus Torvalds   4,126   8%
2011   │ Greg Kroah-...   5,939  11% │ David S. Mi...   4,286   8% │ Linus Torvalds   4,071   7%
2012   │ Greg Kroah-...   7,131  11% │ Linus Torvalds   4,632   7% │ David S. Mi...   4,407   7%
2013   │ Greg Kroah-...   8,285  12% │ David S. Mi...   5,523   8% │ Linus Torvalds   4,955   7%
2014   │ Greg Kroah-...  10,701  14% │ David S. Mi...   6,504   9% │ Linus Torvalds   4,786   6%
2015   │ Greg Kroah-...   9,942  13% │ David S. Mi...   6,823   9% │ Linus Torvalds   4,100   5%
2016   │ David S. Mi...   8,045  10% │ Greg Kroah-...   7,815  10% │ Linus Torvalds   4,713   6%
2017   │ David S. Mi...  10,441  13% │ Greg Kroah-...   5,759   7% │ Linus Torvalds   4,190   5%
2018   │ David S. Mi...   9,411  12% │ Greg Kroah-...   5,909   7% │ Linus Torvalds   3,884   5%
2019   │ David S. Mi...   8,679  11% │ Greg Kroah-...   5,176   6% │ Linus Torvalds   4,159   5%
2020   │ David S. Mi...   8,132   9% │ Linus Torvalds   5,286   6% │ Greg Kroah-...   4,558   5%
2021   │ David S. Mi...   7,643   9% │ Greg Kroah-...   6,116   7% │ Linus Torvalds   4,671   5%
2022   │ Greg Kroah-...   4,446   5% │ David S. Mi...   4,023   5% │ Mark Brown       3,877   5%
2023   │ Mark Brown       4,167   5% │ Greg Kroah-...   4,137   5% │ Alex Deucher     3,916   4%
2024   │ Jakub Kicinski   4,933   6% │ Alex Deucher     3,784   5% │ Greg Kroah-...   3,197   4%
2025   │ Jakub Kicinski   3,999   8% │ Andrew Morton    2,202   4% │ Mark Brown       1,942   4%
------------------------------------------------------------------------------------------------
Total  │ David S. Mi... 113,448   8% │ Greg Kroah-... 104,411   8% │ Linus Torvalds 100,675   7%

Finished in 4m29s with a memory footprint of 2.5 GB.
```

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
