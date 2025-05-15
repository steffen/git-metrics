<h1 align="center">git-metrics</h1>

A powerful Git repository analysis tool that provides detailed metrics, growth statistics, and future projections for your Git repositories.

## Overview

`git-metrics` is a command-line utility that analyzes Git repositories to provide comprehensive insights about repository growth, object statistics, and file usage. The tool gathers historical data and provides projections for future repository growth.

Key features include:
- Repository metadata analysis (first commit, age)
- Year-by-year growth statistics for Git objects (commits, trees, blobs) and their on-disk size
- Identification of largest files in the repository
- File extension distribution analysis
- Future growth projections based on historical trends

## Installation

### Prerequisites
- Git

### Download prebuilt binaries

The easiest way to install `git-metrics` is to download a prebuilt binary from the [GitHub releases page](https://github.com/steffen/git-metrics/releases).

#### Linux
```bash
# Download the latest release for Linux (64-bit)
curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-linux-amd64.tar.gz -o git-metrics.tar.gz

# Unpack the archive
tar -xzf git-metrics.tar.gz

# Optional: Move it to a directory in your PATH
sudo mv git-metrics /usr/local/bin/
```

#### macOS
```bash
# Download the latest release for macOS (Intel or Apple Silicon)
curl -L https://github.com/steffen/git-metrics/releases/latest/download/git-metrics-darwin-arm64.zip -o git-metrics.zip

# Unzip the archive
unzip git-metrics.zip

# Optional: Move it to a directory in your PATH
sudo mv git-metrics /usr/local/bin/

# When downloaded via browser: Remove quarantine attribute
xattr -d com.apple.quarantine git-metrics
```

### Running the tool

```bash
# Analyze the current directory as a Git repository
git-metrics

# Analyze a specific repository
git-metrics -r /path/to/repository
```

## Command line options

| Option | Description |
|--------|-------------|
| `-r`, `--repository` | Path to Git repository (default: current directory) |
| `--debug` | Enable debug output |
| `--no-progress` | Disable progress indicators |

## Output examples

### [`git/git`](https://github.com/git/git)

```
RUN ############################################################################################

Start time                 Thu, 15 May 2025 19:45 CEST
Machine                    10 CPU cores with 64 GB memory (macOS 15.4.1 on Apple M1 Max)
Git metrics version        1.1.0
Git version                2.46.0

REPOSITORY #####################################################################################

Git directory              /Users/steffen/GitHub/oss/git/.git
Remote                     https://github.com/git/git.git
Most recent fetch          Mon, 12 May 2025 19:42 CEST
Most recent commit         Fri, 09 May 2025 (7a1d2bd0a5)
First commit               Thu, 07 Apr 2005 (e83c51)
Age                        20 years 1 months 8 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005          3,215    +4 %          4,056    +3 %          5,922    +4 %         3.9 MB    +1 %
2006          7,816    +6 %         10,459    +4 %         13,181    +5 %         8.0 MB    +2 %
2007         13,312    +7 %         19,425    +6 %         22,503    +6 %        14.5 MB    +2 %
2008         17,440    +5 %         26,655    +5 %         30,759    +5 %        20.0 MB    +2 %
2009         21,267    +5 %         33,673    +4 %         37,227    +4 %        25.3 MB    +2 %
2010         25,150    +5 %         41,060    +5 %         44,099    +4 %        31.1 MB    +2 %
2011         28,673    +4 %         48,136    +5 %         50,231    +4 %        36.4 MB    +2 %
2012         32,455    +5 %         55,937    +5 %         55,808    +4 %        43.3 MB    +3 %
2013         36,772    +5 %         65,037    +6 %         62,775    +5 %        51.3 MB    +3 %
2014         39,875    +4 %         71,108    +4 %         68,121    +3 %        59.7 MB    +3 %
2015         43,161    +4 %         77,534    +4 %         73,436    +3 %        68.6 MB    +3 %
2016         47,031    +5 %         85,302    +5 %         79,873    +4 %        76.8 MB    +3 %
2017         51,617    +6 %         94,468    +6 %         88,493    +6 %        89.9 MB    +5 %
2018         56,098    +6 %        103,740    +6 %         99,098    +7 %       110.0 MB    +8 %
2019         59,867    +5 %        111,620    +5 %        106,638    +5 %       130.5 MB    +8 %
2020         63,562    +5 %        119,536    +5 %        114,111    +5 %       151.0 MB    +8 %
2021         67,578    +5 %        128,086    +5 %        121,937    +5 %       178.5 MB   +10 %
2022         71,225    +5 %        136,302    +5 %        129,471    +5 %       206.0 MB   +10 %
2023         74,172    +4 %        142,980    +4 %        139,020    +6 %       228.7 MB    +9 %
2024         78,574    +5 %        152,860    +6 %        149,414    +7 %       254.0 MB   +10 %
------------------------------------------------------------------------------------------------
2025^        80,261    +2 %        156,505    +2 %        153,235    +2 %       265.9 MB    +4 %
------------------------------------------------------------------------------------------------
2025*        82,315    +5 %        161,108    +5 %        157,969    +6 %       278.7 MB    +9 %
2026*        86,056    +5 %        169,356    +5 %        166,524    +6 %       303.4 MB    +9 %
2027*        89,797    +5 %        177,604    +5 %        175,079    +6 %       328.1 MB    +9 %
2028*        93,538    +5 %        185,852    +5 %        183,634    +6 %       352.8 MB    +9 %
2029*        97,279    +5 %        194,100    +5 %        192,189    +6 %       377.5 MB    +9 %
2030*       101,020    +5 %        202,348    +5 %        200,744    +6 %       402.2 MB    +9 %
------------------------------------------------------------------------------------------------

^ Current totals as of the most recent fetch on Mon, 12 May
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

LARGEST DIRECTORIES ############################################################################

Path                                                        Blobs           On-disk size
------------------------------------------------------------------------------------------------
(root files)                                               72,221  47.1 %        59.1 MB  35.2 %
├─ whats-cooking.txt*                                       1,354   0.9 %         4.3 MB   2.5 %
├─ sequencer.c                                              1,042   0.7 %         3.8 MB   2.2 %
├─ Makefile                                                 3,271   2.1 %         2.8 MB   1.7 %
├─ diff.c                                                   1,626   1.1 %         2.3 MB   1.4 %
├─ read-cache.c                                               836   0.5 %         2.0 MB   1.2 %
├─ cache.h*                                                 2,208   1.4 %         1.4 MB   0.9 %
├─ refs.c                                                   1,285   0.8 %         1.4 MB   0.8 %
├─ merge-recursive.c*                                         775   0.5 %         1.2 MB   0.7 %
├─ config.c                                                   963   0.6 %         1.2 MB   0.7 %
└─ gitk*                                                      568   0.4 %         1.2 MB   0.7 %
------------------------------------------------------------------------------------------------
po                                                          1,389   0.9 %        46.9 MB  27.9 %
├─ fr.po                                                      140   0.1 %         4.9 MB   2.9 %
├─ zh_CN.po                                                   154   0.1 %         4.7 MB   2.8 %
├─ de.po                                                      185   0.1 %         4.6 MB   2.7 %
├─ sv.po                                                      120   0.1 %         4.4 MB   2.6 %
├─ ca.po                                                       78   0.1 %         4.1 MB   2.4 %
├─ vi.po                                                      104   0.1 %         3.7 MB   2.2 %
├─ bg.po                                                       72   0.0 %         3.4 MB   2.1 %
├─ tr.po                                                       41   0.0 %         2.9 MB   1.7 %
├─ zh_TW.po                                                    34   0.0 %         2.5 MB   1.5 %
└─ git.pot*                                                   108   0.1 %         2.0 MB   1.2 %
------------------------------------------------------------------------------------------------
builtin                                                    18,493  12.1 %        20.3 MB  12.1 %
├─ pack-objects.c                                             644   0.4 %         2.0 MB   1.2 %
├─ log.c                                                      596   0.4 %         1.3 MB   0.8 %
├─ submodule--helper.c                                        539   0.4 %         1.2 MB   0.7 %
├─ clone.c                                                    581   0.4 %       975.6 KB   0.6 %
├─ fetch.c                                                    651   0.4 %       960.6 KB   0.6 %
├─ receive-pack.c                                             470   0.3 %       926.2 KB   0.5 %
├─ rebase.c                                                   432   0.3 %       867.7 KB   0.5 %
├─ gc.c                                                       387   0.3 %       689.8 KB   0.4 %
├─ commit.c                                                   683   0.4 %       686.6 KB   0.4 %
└─ index-pack.c                                               362   0.2 %       639.9 KB   0.4 %
------------------------------------------------------------------------------------------------
t                                                          29,241  19.1 %        16.2 MB   9.6 %
├─ test-lib.sh                                                760   0.5 %       847.0 KB   0.5 %
├─ test-lib-functions.sh                                      291   0.2 %       566.6 KB   0.3 %
├─ README                                                     249   0.2 %       484.5 KB   0.3 %
├─ helper                                                   1,453   0.9 %       470.1 KB   0.3 %
├─ t9300-fast-import.sh                                       206   0.1 %       376.6 KB   0.2 %
├─ t0013                                                        1   0.0 %       371.8 KB   0.2 %
├─ t3404-rebase-interactive.sh                                306   0.2 %       323.0 KB   0.2 %
├─ t9902-completion.sh                                        203   0.1 %       306.2 KB   0.2 %
├─ t4014-format-patch.sh                                      204   0.1 %       293.3 KB   0.2 %
└─ t9001-send-email.sh                                        238   0.2 %       290.8 KB   0.2 %
------------------------------------------------------------------------------------------------
Documentation                                              19,695  12.9 %        13.3 MB   7.9 %
├─ RelNotes                                                 1,804   1.2 %         2.8 MB   1.7 %
├─ config.txt*                                              1,459   1.0 %         1.7 MB   1.0 %
├─ technical                                                  810   0.5 %       692.5 KB   0.4 %
├─ git.txt*                                                   851   0.6 %       628.1 KB   0.4 %
├─ config                                                     671   0.4 %       391.5 KB   0.2 %
├─ git-rebase.txt*                                            302   0.2 %       332.6 KB   0.2 %
├─ rev-list-options.txt*                                      243   0.2 %       298.2 KB   0.2 %
├─ diff-options.txt*                                          279   0.2 %       241.8 KB   0.1 %
├─ user-manual.txt*                                           307   0.2 %       215.6 KB   0.1 %
└─ gitattributes.txt*                                         191   0.1 %       212.2 KB   0.1 %
------------------------------------------------------------------------------------------------
contrib                                                     3,764   2.5 %         3.1 MB   1.9 %
├─ completion                                               1,281   0.8 %         1.3 MB   0.8 %
├─ hooks                                                      123   0.1 %       402.6 KB   0.2 %
├─ fast-import                                                464   0.3 %       221.9 KB   0.1 %
├─ buildsystems                                               160   0.1 %       215.9 KB   0.1 %
├─ mw-to-git                                                  195   0.1 %       164.0 KB   0.1 %
├─ subtree                                                    181   0.1 %       122.8 KB   0.1 %
├─ emacs                                                      134   0.1 %       102.1 KB   0.1 %
├─ remote-helpers                                             320   0.2 %        97.3 KB   0.1 %
├─ credential                                                 116   0.1 %        85.3 KB   0.0 %
└─ examples                                                   105   0.1 %        65.5 KB   0.0 %
------------------------------------------------------------------------------------------------
compat                                                      1,227   0.8 %         2.1 MB   1.3 %
├─ mingw.c                                                    360   0.2 %         1.2 MB   0.7 %
├─ regex                                                       56   0.0 %       241.5 KB   0.1 %
├─ nedmalloc                                                   32   0.0 %       191.0 KB   0.1 %
├─ fsmonitor                                                   93   0.1 %        90.8 KB   0.1 %
├─ mingw.h                                                    189   0.1 %        89.9 KB   0.1 %
├─ winansi.c                                                   44   0.0 %        56.2 KB   0.0 %
├─ simple-ipc                                                  26   0.0 %        50.5 KB   0.0 %
├─ win32                                                       80   0.1 %        37.7 KB   0.0 %
├─ terminal.c                                                  31   0.0 %        29.6 KB   0.0 %
└─ vcbuild                                                     38   0.0 %        24.4 KB   0.0 %
------------------------------------------------------------------------------------------------
refs                                                        1,153   0.8 %         1.8 MB   1.1 %
├─ files-backend.c                                            534   0.3 %         1.1 MB   0.7 %
├─ packed-backend.c                                           170   0.1 %       295.6 KB   0.2 %
├─ reftable-backend.c                                         130   0.1 %       209.9 KB   0.1 %
├─ refs-internal.h                                            173   0.1 %       142.7 KB   0.1 %
├─ debug.c                                                     46   0.0 %        35.5 KB   0.0 %
├─ ref-cache.c                                                 41   0.0 %        32.3 KB   0.0 %
├─ iterator.c                                                  21   0.0 %        17.8 KB   0.0 %
├─ ref-cache.h                                                 22   0.0 %        10.1 KB   0.0 %
└─ packed-backend.h                                            16   0.0 %         2.5 KB   0.0 %
------------------------------------------------------------------------------------------------
gitweb                                                      1,131   0.7 %       974.2 KB   0.6 %
├─ gitweb.perl                                                861   0.6 %       807.4 KB   0.5 %
├─ static                                                      38   0.0 %        46.5 KB   0.0 %
├─ README                                                      43   0.0 %        34.4 KB   0.0 %
├─ gitweb.cgi*                                                 65   0.0 %        34.4 KB   0.0 %
├─ INSTALL                                                     27   0.0 %        23.8 KB   0.0 %
├─ Makefile                                                    32   0.0 %        11.1 KB   0.0 %
├─ gitweb.js*                                                   6   0.0 %         8.2 KB   0.0 %
├─ gitweb.css*                                                 50   0.0 %         6.0 KB   0.0 %
├─ meson.build                                                  2   0.0 %         1.0 KB   0.0 %
└─ generate-gitweb-cgi.sh                                       1   0.0 %         0.5 KB   0.0 %
------------------------------------------------------------------------------------------------
git-gui                                                       644   0.4 %       871.4 KB   0.5 %
├─ po                                                         135   0.1 %       503.7 KB   0.3 %
├─ lib                                                        344   0.2 %       190.9 KB   0.1 %
├─ git-gui.sh                                                  85   0.1 %       141.5 KB   0.1 %
├─ Makefile                                                    36   0.0 %        20.3 KB   0.0 %
├─ macosx                                                       7   0.0 %         5.9 KB   0.0 %
├─ README.md                                                    2   0.0 %         2.7 KB   0.0 %
├─ GIT-VERSION-GEN                                             17   0.0 %         1.8 KB   0.0 %
├─ git-gui--askpass                                             4   0.0 %         1.8 KB   0.0 %
├─ windows                                                      4   0.0 %         1.2 KB   0.0 %
└─ CREDITS-GEN*                                                 1   0.0 %         0.8 KB   0.0 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                 148,958  97.2 %       164.6 MB  98.1 %
└─ Out of 41                                              153,235 100.0 %       167.8 MB 100.0 %

* File or directory not present in latest commit of master branch (moved, renamed or removed)

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
po/fr.po                                      2025            140   0.1 %         4.9 MB   2.9 %
po/zh_CN.po                                   2025            154   0.1 %         4.7 MB   2.8 %
po/de.po                                      2025            185   0.1 %         4.6 MB   2.7 %
po/sv.po                                      2025            120   0.1 %         4.4 MB   2.6 %
whats-cooking.txt                             0001          1,354   0.9 %         4.3 MB   2.5 %
po/ca.po                                      2024             78   0.1 %         4.1 MB   2.4 %
sequencer.c                                   2025          1,042   0.7 %         3.8 MB   2.2 %
po/vi.po                                      2025            104   0.1 %         3.7 MB   2.2 %
po/bg.po                                      2025             72   0.0 %         3.4 MB   2.1 %
po/tr.po                                      2025             41   0.0 %         2.9 MB   1.7 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                   3,290   2.1 %        40.7 MB  24.3 %
└─ Out of 6,187                                           153,235 100.0 %       167.8 MB 100.0 %

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                     953  15.4 %         67,986  44.4 %        66.4 MB  25.0 %
.po                                     62   1.0 %          1,363   0.9 %        45.5 MB  17.1 %
.txt                                 1,323  21.4 %         20,441  13.3 %        16.8 MB   6.3 %
.sh                                  1,521  24.6 %         30,263  19.7 %        15.9 MB   6.0 %
No Extension                           699  11.3 %         10,025   6.5 %         7.0 MB   2.6 %
.h                                     380   6.1 %         12,310   8.0 %         5.5 MB   2.1 %
.perl                                   57   0.9 %          3,033   2.0 %         2.2 MB   0.8 %
.pot                                     5   0.1 %            128   0.1 %         2.1 MB   0.8 %
.py                                     30   0.5 %            534   0.3 %         1.4 MB   0.5 %
.bash                                    2   0.0 %          1,117   0.7 %         1.3 MB   0.5 %
------------------------------------------------------------------------------------------------
├─ Top 10                            5,032  81.3 %        147,200  96.1 %       164.0 MB  97.8 %
└─ Out of 363                        6,187 100.0 %        153,234 100.0 %       167.8 MB 100.0 %

AUTHORS & COMMITTERS WITH MOST COMMITS #########################################################

Year    Author                   Commits                Committer                Commits
------------------------------------------------------------------------------------------------

2005    Junio C Hamano             1,368   42.6%        Junio C Hamano             1,794   55.8%
        Linus Torvalds               680   21.2%        Linus Torvalds             1,036   32.2%
        Kay Sievers                  139    4.3%        Kay Sievers                  136    4.2%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,187   68.0%        ├─ Top 3                   2,966   92.3%
        └─ Out of 131              3,215  100.0%        └─ Out of 13               3,215  100.0%

2006    Junio C Hamano             2,202   47.9%        Junio C Hamano             4,310   93.7%
        Shawn O. Pearce              243    5.3%        Shawn O. Pearce              158    3.4%
        Jakub Narebski               220    4.8%        Paul Mackerras                82    1.8%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,665   57.9%        ├─ Top 3                   4,550   98.9%
        └─ Out of 163              4,601  100.0%        └─ Out of 7                4,601  100.0%

2007    Junio C Hamano             1,606   29.2%        Junio C Hamano             3,964   72.1%
        Shawn O. Pearce              841   15.3%        Shawn O. Pearce              842   15.3%
        Johannes Schindelin          244    4.4%        Simon Hausmann               247    4.5%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,691   49.0%        ├─ Top 3                   5,053   91.9%
        └─ Out of 285              5,496  100.0%        └─ Out of 10               5,496  100.0%

2008    Junio C Hamano             1,365   33.1%        Junio C Hamano             3,532   85.7%
        Shawn O. Pearce              163    4.0%        Shawn O. Pearce              372    9.0%
        Jeff King                    152    3.7%        Paul Mackerras               123    3.0%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,680   40.8%        ├─ Top 3                   4,027   97.7%
        └─ Out of 311              4,120  100.0%        └─ Out of 9                4,120  100.0%

2009    Junio C Hamano             1,439   37.5%        Junio C Hamano             3,514   91.6%
        Jeff King                    134    3.5%        Eric Wong                     96    2.5%
        Johannes Schindelin          104    2.7%        Avery Pennarun                65    1.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,677   43.7%        ├─ Top 3                   3,675   95.8%
        └─ Out of 295              3,835  100.0%        └─ Out of 10               3,835  100.0%

2010    Junio C Hamano             1,537   39.6%        Junio C Hamano             3,665   94.4%
        Jonathan Nieder              333    8.6%        Eric Wong                     49    1.3%
        Ævar Arnfjörð Bjarmason      121    3.1%        Pat Thoyts                    47    1.2%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,991   51.3%        ├─ Top 3                   3,761   96.9%
        └─ Out of 254              3,883  100.0%        └─ Out of 10               3,883  100.0%

2011    Junio C Hamano             1,588   45.1%        Junio C Hamano             3,328   94.5%
        Jeff King                    237    6.7%        Jonathan Nieder               80    2.3%
        Jonathan Nieder              190    5.4%        Pat Thoyts                    70    2.0%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,015   57.2%        ├─ Top 3                   3,478   98.8%
        └─ Out of 191              3,521  100.0%        └─ Out of 7                3,521  100.0%

2012    Junio C Hamano             1,743   46.1%        Junio C Hamano             3,348   88.5%
        Jeff King                    314    8.3%        Jeff King                    153    4.0%
        Nguyễn Thái Ngọc Duy         204    5.4%        Jiang Xin                     91    2.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,261   59.8%        ├─ Top 3                   3,592   95.0%
        └─ Out of 214              3,782  100.0%        └─ Out of 16               3,782  100.0%

2013    Junio C Hamano             1,787   41.4%        Junio C Hamano             4,093   94.8%
        Felipe Contreras             283    6.6%        Jonathan Nieder               81    1.9%
        Jeff King                    243    5.6%        Jiang Xin                     24    0.6%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,313   53.6%        ├─ Top 3                   4,198   97.2%
        └─ Out of 238              4,319  100.0%        └─ Out of 13               4,319  100.0%

2014    Junio C Hamano             1,220   39.3%        Junio C Hamano             2,967   95.6%
        Jeff King                    340   11.0%        Jiang Xin                     27    0.9%
        Nguyễn Thái Ngọc Duy         140    4.5%        Eric Wong                     23    0.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,700   54.8%        ├─ Top 3                   3,017   97.2%
        └─ Out of 209              3,103  100.0%        └─ Out of 12               3,103  100.0%

2015    Junio C Hamano             1,440   43.8%        Junio C Hamano             2,970   90.4%
        Jeff King                    368   11.2%        Jeff King                    111    3.4%
        Michael Haggerty             173    5.3%        Jiang Xin                     63    1.9%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,981   60.3%        ├─ Top 3                   3,144   95.7%
        └─ Out of 191              3,286  100.0%        └─ Out of 18               3,286  100.0%

2016    Junio C Hamano             1,604   41.4%        Junio C Hamano             3,652   94.4%
        Jeff King                    386   10.0%        Jiang Xin                     45    1.2%
        Johannes Schindelin          193    5.0%        Michael Haggerty              34    0.9%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,183   56.4%        ├─ Top 3                   3,731   96.4%
        └─ Out of 196              3,870  100.0%        └─ Out of 19               3,870  100.0%

2017    Junio C Hamano             1,739   37.9%        Junio C Hamano             4,444   96.9%
        Jeff King                    405    8.8%        Jiang Xin                     64    1.4%
        Brandon Williams             217    4.7%        Jean-Noel Avila               15    0.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,361   51.5%        ├─ Top 3                   4,523   98.6%
        └─ Out of 204              4,585  100.0%        └─ Out of 18               4,585  100.0%

2018    Junio C Hamano             1,230   27.4%        Junio C Hamano             4,327   96.5%
        Nguyễn Thái Ngọc Duy         499   11.1%        Jiang Xin                     64    1.4%
        Jeff King                    255    5.7%        Jeff King                     21    0.5%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,984   44.3%        ├─ Top 3                   4,412   98.4%
        └─ Out of 190              4,482  100.0%        └─ Out of 16               4,482  100.0%

2019    Junio C Hamano             1,095   29.1%        Junio C Hamano             3,540   93.9%
        Johannes Schindelin          321    8.5%        Johannes Schindelin           67    1.8%
        Jeff King                    271    7.2%        Jiang Xin                     53    1.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,687   44.8%        ├─ Top 3                   3,660   97.1%
        └─ Out of 217              3,768  100.0%        └─ Out of 18               3,768  100.0%

2020    Junio C Hamano             1,157   31.3%        Junio C Hamano             3,474   94.0%
        Jeff King                    270    7.3%        Jiang Xin                     62    1.7%
        Johannes Schindelin          208    5.6%        Pratyush Yadav                41    1.1%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,635   44.2%        ├─ Top 3                   3,577   96.8%
        └─ Out of 210              3,696  100.0%        └─ Out of 20               3,696  100.0%

2021    Junio C Hamano             1,215   30.3%        Junio C Hamano             3,796   94.5%
        Ævar Arnfjörð Bjarmason      606   15.1%        Jiang Xin                     81    2.0%
        Elijah Newren                201    5.0%        Johannes Schindelin           30    0.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,022   50.3%        ├─ Top 3                   3,907   97.3%
        └─ Out of 179              4,016  100.0%        └─ Out of 22               4,016  100.0%

2022    Junio C Hamano             1,187   32.5%        Junio C Hamano             3,239   88.8%
        Ævar Arnfjörð Bjarmason      555   15.2%        Taylor Blau                  240    6.6%
        Taylor Blau                  181    5.0%        Jiang Xin                     58    1.6%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,923   52.7%        ├─ Top 3                   3,537   97.0%
        └─ Out of 168              3,647  100.0%        └─ Out of 19               3,647  100.0%

2023    Junio C Hamano             1,041   35.3%        Junio C Hamano             2,765   93.8%
        Jeff King                    330   11.2%        Johannes Schindelin           85    2.9%
        Elijah Newren                195    6.6%        Jiang Xin                     36    1.2%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,566   53.1%        ├─ Top 3                   2,886   97.9%
        └─ Out of 150              2,947  100.0%        └─ Out of 17               2,947  100.0%

2024    Junio C Hamano             1,578   35.9%        Junio C Hamano             3,921   89.2%
        Patrick Steinhardt         1,026   23.4%        Taylor Blau                  262    6.0%
        Jeff King                    259    5.9%        Johannes Schindelin           86    2.0%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,863   65.2%        ├─ Top 3                   4,269   97.2%
        └─ Out of 173              4,394  100.0%        └─ Out of 19               4,394  100.0%

2025    Junio C Hamano               741   43.7%        Junio C Hamano             1,637   96.6%
        Patrick Steinhardt           299   17.6%        Jiang Xin                     21    1.2%
        Jeff King                     55    3.2%        Johannes Sixt                 12    0.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,095   64.6%        ├─ Top 3                   1,670   98.5%
        └─ Out of 101              1,695  100.0%        └─ Out of 16               1,695  100.0%

------------------------------------------------------------------------------------------------

TOTAL   Junio C Hamano            29,882   37.2%        Junio C Hamano            72,280   90.1%
        Jeff King                  4,558    5.7%        Shawn O. Pearce            1,457    1.8%
        Johannes Schindelin        2,361    2.9%        Linus Torvalds             1,041    1.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                  36,801   45.9%        ├─ Top 3                  74,778   93.2%
        └─ Out of 2,289           80,261  100.0%        └─ Out of 88              80,261  100.0%

Finished in 12s with a memory footprint of 104.6 MB.
```

### [`torvalds/linux`](https://github.com/torvalds/linux)

```
RUN ############################################################################################

Start time                 Thu, 15 May 2025 19:47 CEST
Machine                    10 CPU cores with 64 GB memory (macOS 15.4.1 on Apple M1 Max)
Git metrics version        1.1.0
Git version                2.46.0

REPOSITORY #####################################################################################

Git directory              /Users/steffen/GitHub/oss/linux/.git
Remote                     https://github.com/torvalds/linux.git
Most recent fetch          Mon, 12 May 2025 19:36 CEST
Most recent commit         Sun, 11 May 2025 (627277ba7c23)
First commit               Sat, 16 Apr 2005 (1da177)
Age                        20 years 29 days

HISTORIC & ESTIMATED GROWTH ####################################################################

Year        Commits                  Trees                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
2005         15,862    +1 %         71,850    +1 %         63,135    +2 %       121.5 MB    +2 %
2006         45,307    +2 %        204,857    +2 %        147,863    +3 %       184.9 MB    +1 %
2007         75,872    +2 %        339,474    +2 %        234,445    +3 %       258.9 MB    +1 %
2008        126,734    +4 %        562,425    +3 %        351,965    +4 %       369.0 MB    +2 %
2009        179,269    +4 %        800,802    +4 %        474,393    +4 %       501.9 MB    +2 %
2010        228,892    +4 %      1,030,868    +4 %        596,418    +4 %       627.0 MB    +2 %
2011        284,002    +4 %      1,290,525    +4 %        730,819    +5 %       806.6 MB    +3 %
2012        348,966    +5 %      1,605,127    +5 %        882,531    +5 %      1006.2 MB    +4 %
2013        420,327    +5 %      1,938,511    +5 %      1,027,275    +5 %         1.1 GB    +3 %
2014        496,286    +6 %      2,296,647    +5 %      1,177,798    +5 %         1.3 GB    +3 %
2015        571,741    +6 %      2,654,431    +5 %      1,328,032    +5 %         1.5 GB    +4 %
2016        648,806    +6 %      3,025,936    +6 %      1,485,414    +5 %         1.8 GB    +5 %
2017        729,675    +6 %      3,432,943    +6 %      1,675,158    +6 %         2.1 GB    +7 %
2018        810,066    +6 %      3,825,666    +6 %      1,832,028    +5 %         2.4 GB    +5 %
2019        892,609    +6 %      4,243,177    +6 %      2,021,235    +6 %         2.7 GB    +6 %
2020        983,052    +7 %      4,693,652    +7 %      2,203,820    +6 %         3.1 GB    +8 %
2021      1,069,168    +6 %      5,117,022    +6 %      2,368,647    +6 %         3.7 GB   +10 %
2022      1,155,031    +6 %      5,544,636    +7 %      2,539,526    +6 %         4.2 GB   +11 %
2023      1,246,546    +7 %      5,991,760    +7 %      2,721,932    +6 %         4.7 GB    +9 %
2024      1,329,852    +6 %      6,403,293    +6 %      2,885,958    +6 %         5.2 GB    +8 %
------------------------------------------------------------------------------------------------
2025^     1,352,887    +2 %      6,515,858    +2 %      2,930,206    +2 %         5.3 GB    +2 %
------------------------------------------------------------------------------------------------
2025*     1,417,300    +6 %      6,835,316    +7 %      3,058,902    +6 %         5.7 GB    +9 %
2026*     1,504,748    +6 %      7,267,339    +7 %      3,231,846    +6 %         6.1 GB    +9 %
2027*     1,592,196    +6 %      7,699,362    +7 %      3,404,790    +6 %         6.6 GB    +9 %
2028*     1,679,644    +6 %      8,131,385    +7 %      3,577,734    +6 %         7.1 GB    +9 %
2029*     1,767,092    +6 %      8,563,408    +7 %      3,750,678    +6 %         7.6 GB    +9 %
2030*     1,854,540    +6 %      8,995,431    +7 %      3,923,622    +6 %         8.1 GB    +9 %
------------------------------------------------------------------------------------------------

^ Current totals as of the most recent fetch on Mon, 12 May
* Estimated growth based on the last five years
% Percentages show the increase relative to the current total (^)

LARGEST DIRECTORIES ############################################################################

Path                                                        Blobs           On-disk size
------------------------------------------------------------------------------------------------
drivers                                                 1,398,094  47.7 %         1.9 GB  52.1 %
├─ net                                                    288,961   9.9 %       473.8 MB  12.5 %
├─ gpu                                                    259,529   8.9 %       426.3 MB  11.2 %
├─ staging                                                152,344   5.2 %       132.6 MB   3.5 %
├─ scsi                                                    52,762   1.8 %       117.9 MB   3.1 %
├─ media                                                   80,961   2.8 %        97.8 MB   2.6 %
├─ usb                                                     46,497   1.6 %        63.2 MB   1.7 %
├─ infiniband                                              31,284   1.1 %        52.7 MB   1.4 %
├─ clk                                                     19,746   0.7 %        25.5 MB   0.7 %
├─ md                                                      15,078   0.5 %        24.9 MB   0.7 %
└─ video                                                   17,406   0.6 %        23.7 MB   0.6 %
------------------------------------------------------------------------------------------------
arch                                                      517,362  17.7 %       410.7 MB  10.8 %
├─ arm                                                    152,543   5.2 %        90.4 MB   2.4 %
├─ x86                                                     86,387   2.9 %        87.1 MB   2.3 %
├─ arm64                                                   54,753   1.9 %        56.8 MB   1.5 %
├─ powerpc                                                 61,740   2.1 %        55.3 MB   1.5 %
├─ mips                                                    33,315   1.1 %        24.2 MB   0.6 %
├─ s390                                                    19,367   0.7 %        16.1 MB   0.4 %
├─ sparc                                                    9,203   0.3 %         7.9 MB   0.2 %
├─ sh                                                      12,890   0.4 %         7.5 MB   0.2 %
├─ ia64*                                                    7,033   0.2 %         6.2 MB   0.2 %
└─ m68k                                                     6,787   0.2 %         6.1 MB   0.2 %
------------------------------------------------------------------------------------------------
fs                                                        206,793   7.1 %       269.6 MB   7.1 %
├─ btrfs                                                   28,606   1.0 %        47.9 MB   1.3 %
├─ xfs                                                     27,349   0.9 %        31.1 MB   0.8 %
├─ ext4                                                     8,866   0.3 %        15.2 MB   0.4 %
├─ nfs                                                     12,100   0.4 %        14.0 MB   0.4 %
├─ cifs*                                                    9,482   0.3 %        13.4 MB   0.4 %
├─ f2fs                                                     8,523   0.3 %        11.2 MB   0.3 %
├─ ocfs2                                                    5,885   0.2 %         8.5 MB   0.2 %
├─ nfsd                                                     6,755   0.2 %         8.4 MB   0.2 %
├─ bcachefs                                                14,091   0.5 %         8.4 MB   0.2 %
└─ io_uring.c*                                              2,054   0.1 %         7.9 MB   0.2 %
------------------------------------------------------------------------------------------------
net                                                       160,938   5.5 %       202.4 MB   5.3 %
├─ ipv4                                                    22,108   0.8 %        28.8 MB   0.8 %
├─ core                                                    11,514   0.4 %        22.4 MB   0.6 %
├─ mac80211                                                14,105   0.5 %        18.0 MB   0.5 %
├─ ipv6                                                    13,341   0.5 %        15.4 MB   0.4 %
├─ netfilter                                               13,427   0.5 %        14.5 MB   0.4 %
├─ bluetooth                                                7,710   0.3 %        11.0 MB   0.3 %
├─ wireless                                                 4,791   0.2 %         9.2 MB   0.2 %
├─ sched                                                    7,651   0.3 %         8.4 MB   0.2 %
├─ sunrpc                                                   7,553   0.3 %         7.7 MB   0.2 %
└─ sctp                                                     3,788   0.1 %         7.1 MB   0.2 %
------------------------------------------------------------------------------------------------
(root files)                                               22,706   0.8 %       188.1 MB   5.0 %
├─ MAINTAINERS                                             18,933   0.6 %       183.9 MB   4.8 %
├─ Makefile                                                 2,552   0.1 %         2.6 MB   0.1 %
├─ CREDITS                                                    330   0.0 %       919.9 KB   0.0 %
├─ .mailmap                                                   551   0.0 %       574.3 KB   0.0 %
├─ .clang-format                                               54   0.0 %        59.2 KB   0.0 %
├─ README                                                      46   0.0 %        45.3 KB   0.0 %
├─ .gitignore                                                 136   0.0 %        26.2 KB   0.0 %
├─ Kbuild                                                      53   0.0 %        13.6 KB   0.0 %
├─ REPORTING-BUGS*                                             15   0.0 %         9.0 KB   0.0 %
└─ COPYING                                                      4   0.0 %         8.0 KB   0.0 %
------------------------------------------------------------------------------------------------
include                                                   191,451   6.5 %       158.8 MB   4.2 %
├─ linux                                                  100,372   3.4 %        84.1 MB   2.2 %
├─ net                                                     20,842   0.7 %        19.1 MB   0.5 %
├─ uapi                                                    12,295   0.4 %        17.5 MB   0.5 %
├─ drm                                                      5,620   0.2 %         5.4 MB   0.1 %
├─ sound                                                    3,987   0.1 %         3.3 MB   0.1 %
├─ trace                                                    3,326   0.1 %         2.9 MB   0.1 %
├─ media                                                    3,238   0.1 %         2.6 MB   0.1 %
├─ rdma                                                     1,622   0.1 %         2.2 MB   0.1 %
├─ acpi                                                     2,893   0.1 %         2.1 MB   0.1 %
└─ dt-bindings                                              2,882   0.1 %         2.0 MB   0.1 %
------------------------------------------------------------------------------------------------
sound                                                      89,391   3.1 %       120.2 MB   3.2 %
├─ soc                                                     53,474   1.8 %        70.1 MB   1.8 %
├─ pci                                                     18,069   0.6 %        28.1 MB   0.7 %
├─ usb                                                      3,917   0.1 %         5.6 MB   0.1 %
├─ core                                                     3,788   0.1 %         5.0 MB   0.1 %
├─ oss                                                      1,150   0.0 %         2.7 MB   0.1 %
├─ isa                                                      2,027   0.1 %         2.2 MB   0.1 %
├─ firewire                                                 2,651   0.1 %         1.8 MB   0.0 %
├─ drivers                                                    962   0.0 %         1.1 MB   0.0 %
├─ hda                                                        528   0.0 %       537.8 KB   0.0 %
└─ sparc                                                      237   0.0 %       516.8 KB   0.0 %
------------------------------------------------------------------------------------------------
kernel                                                     62,365   2.1 %       112.6 MB   3.0 %
├─ bpf                                                      5,275   0.2 %        21.9 MB   0.6 %
├─ sched                                                    6,523   0.2 %        18.4 MB   0.5 %
├─ trace                                                    9,272   0.3 %        15.0 MB   0.4 %
├─ events                                                   1,757   0.1 %         7.1 MB   0.2 %
├─ rcu                                                      3,920   0.1 %         5.9 MB   0.2 %
├─ time                                                     3,253   0.1 %         3.6 MB   0.1 %
├─ cgroup                                                     924   0.0 %         3.0 MB   0.1 %
├─ workqueue.c                                                946   0.0 %         2.7 MB   0.1 %
├─ locking                                                  1,393   0.0 %         2.5 MB   0.1 %
└─ irq                                                      2,345   0.1 %         2.0 MB   0.1 %
------------------------------------------------------------------------------------------------
Documentation                                              91,115   3.1 %       105.5 MB   2.8 %
├─ devicetree                                              37,784   1.3 %        22.4 MB   0.6 %
├─ admin-guide                                              3,601   0.1 %        11.5 MB   0.3 %
├─ networking                                               3,169   0.1 %         6.7 MB   0.2 %
├─ translations                                             2,136   0.1 %         6.0 MB   0.2 %
├─ filesystems                                              2,758   0.1 %         5.4 MB   0.1 %
├─ DocBook*                                                 2,444   0.1 %         3.0 MB   0.1 %
├─ media*                                                   4,291   0.1 %         3.0 MB   0.1 %
├─ ABI                                                      3,128   0.1 %         2.5 MB   0.1 %
├─ virt                                                       662   0.0 %         2.5 MB   0.1 %
└─ userspace-api                                            2,131   0.1 %         2.5 MB   0.1 %
------------------------------------------------------------------------------------------------
tools                                                      89,628   3.1 %        99.8 MB   2.6 %
├─ perf                                                    42,715   1.5 %        43.0 MB   1.1 %
├─ testing                                                 31,446   1.1 %        32.6 MB   0.9 %
├─ lib                                                      4,541   0.2 %        10.8 MB   0.3 %
├─ power                                                    2,051   0.1 %         3.7 MB   0.1 %
├─ bpf                                                      1,773   0.1 %         2.2 MB   0.1 %
├─ include                                                  1,299   0.0 %         1.6 MB   0.0 %
├─ objtool                                                  1,375   0.0 %         1.5 MB   0.0 %
├─ net                                                        721   0.0 %       843.1 KB   0.0 %
├─ memory-model                                               387   0.0 %       497.5 KB   0.0 %
└─ tracing                                                    339   0.0 %       482.5 KB   0.0 %
------------------------------------------------------------------------------------------------
├─ Top 10                                               2,829,843  96.6 %         3.6 GB  96.0 %
└─ Out of 28                                            2,930,206 100.0 %         3.7 GB 100.0 %

* File or directory not present in latest commit of master branch (moved, renamed or removed)

LARGEST FILES ##################################################################################

File path                              Last commit          Blobs           On-disk size
------------------------------------------------------------------------------------------------
MAINTAINERS                                   2025         18,933   0.6 %       183.9 MB   4.8 %
kernel/bpf/verifier.c                         2025          1,469   0.1 %        14.3 MB   0.4 %
drivers/gpu/drm/i9...ay/intel_display.c [1]   2025          1,475   0.1 %         9.3 MB   0.2 %
drivers/gpu/drm/i915/intel_display.c          2019          4,155   0.1 %         9.1 MB   0.2 %
drivers/gpu/drm/i915/i915_reg.h               2025          2,469   0.1 %         8.9 MB   0.2 %
drivers/gpu/drm/am...gpu_dm/amdgpu_dm.c [2]   2025          1,904   0.1 %         8.2 MB   0.2 %
arch/x86/kvm/x86.c                            2025          3,058   0.1 %         8.1 MB   0.2 %
fs/io_uring.c                                 2022          2,054   0.1 %         7.9 MB   0.2 %
crypto/testmgr.h                              2025            226   0.0 %         7.1 MB   0.2 %
kernel/sched/fair.c                           2025          1,549   0.1 %         7.0 MB   0.2 %
------------------------------------------------------------------------------------------------
├─ Top 10                                                  37,292   1.3 %       263.8 MB   7.0 %
└─ Out of 154,839                                       2,930,206 100.0 %         3.7 GB 100.0 %

[1] drivers/gpu/drm/i915/display/intel_display.c
[2] drivers/gpu/drm/amd/display/amdgpu_dm/amdgpu_dm.c

LARGEST FILE EXTENSIONS ########################################################################

Extension                            Files                  Blobs           On-disk size
------------------------------------------------------------------------------------------------
.c                                  59,508  38.4 %      1,904,819  65.0 %         2.7 GB  50.7 %
.h                                  51,191  33.1 %        599,969  20.5 %       551.3 MB  10.2 %
No Extension                        11,095   7.2 %        197,547   6.7 %       269.3 MB   5.0 %
.dtsi                                3,312   2.1 %         44,457   1.5 %        48.4 MB   0.9 %
.rst                                 5,301   3.4 %         25,924   0.9 %        45.0 MB   0.8 %
.txt                                 6,128   4.0 %         33,716   1.2 %        35.9 MB   0.7 %
.S                                   2,943   1.9 %         29,354   1.0 %        24.9 MB   0.5 %
.dts                                 4,556   2.9 %         34,867   1.2 %        24.6 MB   0.5 %
.yaml                                4,858   3.1 %         22,890   0.8 %        15.5 MB   0.3 %
.json                                  866   0.6 %          3,221   0.1 %         6.9 MB   0.1 %
------------------------------------------------------------------------------------------------
├─ Top 10                          149,758  96.7 %      2,896,764  98.9 %         3.7 GB  98.9 %
└─ Out of 441                      154,839 100.0 %      2,930,206 100.0 %         3.7 GB 100.0 %

AUTHORS & COMMITTERS WITH MOST COMMITS #########################################################

Year    Author                   Commits                Committer                Commits
------------------------------------------------------------------------------------------------

2005    Linus Torvalds               775    4.9%        Linus Torvalds             6,398   40.3%
        Jeff Garzik                  392    2.5%        David S. Miller            1,384    8.7%
        Russell King                 344    2.2%        Jeff Garzik                1,384    8.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,511    9.5%        ├─ Top 3                   9,166   57.8%
        └─ Out of 1,278           15,862  100.0%        └─ Out of 48              15,862  100.0%

2006    Linus Torvalds             1,108    3.8%        Linus Torvalds             9,384   31.9%
        Al Viro                      765    2.6%        David S. Miller            3,115   10.6%
        David S. Miller              612    2.1%        Greg Kroah-Hartman         1,842    6.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,485    8.4%        ├─ Top 3                  14,341   48.7%
        └─ Out of 1,935           29,445  100.0%        └─ Out of 82              29,445  100.0%

2007    Linus Torvalds             1,394    4.6%        Linus Torvalds             7,419   24.3%
        Ralf Baechle                 506    1.7%        David S. Miller            3,337   10.9%
        Thomas Gleixner              484    1.6%        Jeff Garzik                2,040    6.7%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,384    7.8%        ├─ Top 3                  12,796   41.9%
        └─ Out of 2,028           30,565  100.0%        └─ Out of 92              30,565  100.0%

2008    Linus Torvalds             1,912    3.8%        Linus Torvalds             7,305   14.4%
        Ingo Molnar                1,271    2.5%        Ingo Molnar                5,983   11.8%
        David S. Miller              928    1.8%        David S. Miller            5,730   11.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,111    8.1%        ├─ Top 3                  19,018   37.4%
        └─ Out of 2,506           50,869  100.0%        └─ Out of 137             50,869  100.0%

2009    Linus Torvalds             2,124    4.0%        Linus Torvalds             5,731   10.9%
        Ingo Molnar                1,088    2.1%        David S. Miller            5,367   10.2%
        Takashi Iwai                 952    1.8%        Ingo Molnar                4,362    8.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,164    7.9%        ├─ Top 3                  15,460   29.4%
        └─ Out of 2,786           52,581  100.0%        └─ Out of 163             52,581  100.0%

2010    Linus Torvalds             1,884    3.8%        David S. Miller            5,073   10.2%
        Joe Perches                  546    1.1%        Greg Kroah-Hartman         4,497    9.1%
        Chris Wilson                 519    1.0%        Linus Torvalds             4,126    8.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   2,949    5.9%        ├─ Top 3                  13,696   27.6%
        └─ Out of 2,942           49,634  100.0%        └─ Out of 195             49,634  100.0%

2011    Linus Torvalds             2,080    3.8%        Greg Kroah-Hartman         5,939   10.8%
        Mark Brown                 1,047    1.9%        David S. Miller            4,286    7.8%
        David S. Miller              743    1.3%        Linus Torvalds             4,071    7.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   3,870    7.0%        ├─ Top 3                  14,296   25.9%
        └─ Out of 3,028           55,111  100.0%        └─ Out of 219             55,111  100.0%

2012    Linus Torvalds             2,271    3.5%        Greg Kroah-Hartman         7,131   11.0%
        H Hartley Sweeten          1,447    2.2%        Linus Torvalds             4,632    7.1%
        Mark Brown                 1,224    1.9%        David S. Miller            4,407    6.8%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,942    7.6%        ├─ Top 3                  16,170   24.9%
        └─ Out of 3,090           64,959  100.0%        └─ Out of 250             64,959  100.0%

2013    Linus Torvalds             2,044    2.9%        Greg Kroah-Hartman         8,285   11.6%
        H Hartley Sweeten          1,582    2.2%        David S. Miller            5,523    7.7%
        Mark Brown                 1,506    2.1%        Linus Torvalds             4,955    6.9%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   5,132    7.2%        ├─ Top 3                  18,763   26.3%
        └─ Out of 3,403           71,356  100.0%        └─ Out of 253             71,356  100.0%

2014    Linus Torvalds             2,085    2.7%        Greg Kroah-Hartman        10,701   14.1%
        H Hartley Sweeten          1,620    2.1%        David S. Miller            6,504    8.6%
        David S. Miller              922    1.2%        Linus Torvalds             4,786    6.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,627    6.1%        ├─ Top 3                  21,991   29.0%
        └─ Out of 3,739           75,959  100.0%        └─ Out of 263             75,959  100.0%

2015    Linus Torvalds             2,009    2.7%        Greg Kroah-Hartman         9,942   13.2%
        David S. Miller              987    1.3%        David S. Miller            6,823    9.0%
        H Hartley Sweeten            784    1.0%        Linus Torvalds             4,100    5.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   3,780    5.0%        ├─ Top 3                  20,865   27.7%
        └─ Out of 4,040           75,415  100.0%        └─ Out of 260             75,415  100.0%

2016    Linus Torvalds             2,270    2.9%        David S. Miller            8,045   10.4%
        Arnd Bergmann              1,186    1.5%        Greg Kroah-Hartman         7,815   10.1%
        David S. Miller            1,151    1.5%        Linus Torvalds             4,713    6.1%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,607    6.0%        ├─ Top 3                  20,573   26.7%
        └─ Out of 4,056           77,129  100.0%        └─ Out of 273             77,129  100.0%

2017    Linus Torvalds             2,303    2.8%        David S. Miller           10,441   12.9%
        David S. Miller            1,420    1.8%        Greg Kroah-Hartman         5,759    7.1%
        Arnd Bergmann              1,121    1.4%        Linus Torvalds             4,190    5.2%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,844    6.0%        ├─ Top 3                  20,390   25.2%
        └─ Out of 4,443           80,866  100.0%        └─ Out of 275             80,866  100.0%

2018    Linus Torvalds             2,172    2.7%        David S. Miller            9,411   11.7%
        David S. Miller            1,405    1.7%        Greg Kroah-Hartman         5,909    7.3%
        Arnd Bergmann                904    1.1%        Linus Torvalds             3,884    4.8%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,481    5.6%        ├─ Top 3                  19,204   23.9%
        └─ Out of 4,406           80,411  100.0%        └─ Out of 283             80,411  100.0%

2019    Linus Torvalds             2,384    2.9%        David S. Miller            8,679   10.5%
        David S. Miller            1,207    1.5%        Greg Kroah-Hartman         5,176    6.3%
        Chris Wilson               1,174    1.4%        Linus Torvalds             4,159    5.0%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,765    5.8%        ├─ Top 3                  18,014   21.8%
        └─ Out of 4,385           82,540  100.0%        └─ Out of 301             82,540  100.0%

2020    Linus Torvalds             2,635    2.9%        David S. Miller            8,132    9.0%
        Mauro Carvalho Chehab      1,215    1.3%        Linus Torvalds             5,286    5.8%
        Christoph Hellwig          1,200    1.3%        Greg Kroah-Hartman         4,558    5.0%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   5,050    5.6%        ├─ Top 3                  17,976   19.9%
        └─ Out of 4,636           90,448  100.0%        └─ Out of 289             90,448  100.0%

2021    Linus Torvalds             2,415    2.8%        David S. Miller            7,643    8.9%
        David S. Miller            1,042    1.2%        Greg Kroah-Hartman         6,116    7.1%
        Arnd Bergmann                978    1.1%        Linus Torvalds             4,671    5.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,435    5.2%        ├─ Top 3                  18,430   21.4%
        └─ Out of 4,800           86,115  100.0%        └─ Out of 295             86,115  100.0%

2022    Linus Torvalds             2,601    3.0%        Greg Kroah-Hartman         4,446    5.2%
        Krzysztof Kozlowski        1,371    1.6%        David S. Miller            4,023    4.7%
        Jakub Kicinski               981    1.1%        Mark Brown                 3,877    4.5%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   4,953    5.8%        ├─ Top 3                  12,346   14.4%
        └─ Out of 4,931           85,820  100.0%        └─ Out of 312             85,820  100.0%

2023    Kent Overstreet            2,813    3.1%        Mark Brown                 4,167    4.6%
        Linus Torvalds             2,505    2.7%        Greg Kroah-Hartman         4,137    4.5%
        Uwe Kleine-König           2,235    2.4%        Alex Deucher               3,916    4.3%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   7,553    8.3%        ├─ Top 3                  12,220   13.4%
        └─ Out of 4,938           91,340  100.0%        └─ Out of 319             91,340  100.0%

2024    Linus Torvalds             2,889    3.5%        Jakub Kicinski             4,933    5.9%
        Kent Overstreet            1,368    1.6%        Alex Deucher               3,784    4.5%
        Krzysztof Kozlowski        1,180    1.4%        Greg Kroah-Hartman         3,197    3.8%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   5,437    6.5%        ├─ Top 3                  11,914   14.3%
        └─ Out of 5,109           83,526  100.0%        └─ Out of 337             83,526  100.0%

2025    Linus Torvalds             1,007    4.4%        Jakub Kicinski             1,877    8.1%
        Jakub Kicinski               458    2.0%        Andrew Morton              1,181    5.1%
        Kent Overstreet              324    1.4%        Linus Torvalds             1,026    4.5%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                   1,789    7.8%        ├─ Top 3                   4,084   17.7%
        └─ Out of 2,659           23,036  100.0%        └─ Out of 287             23,036  100.0%

------------------------------------------------------------------------------------------------

TOTAL   Linus Torvalds            42,867    3.2%        David S. Miller          113,351    8.4%
        David S. Miller           15,651    1.2%        Greg Kroah-Hartman       103,411    7.6%
        Arnd Bergmann             11,674    0.9%        Linus Torvalds            99,864    7.4%
        ┌───────────────────────────────────────        ┌───────────────────────────────────────
        ├─ Top 3                  70,192    5.2%        ├─ Top 3                 316,626   23.4%
        └─ Out of 30,128       1,352,987  100.0%        └─ Out of 895          1,352,987  100.0%

Finished in 4m7s with a memory footprint of 2.2 GB.
```

## Understanding the output

`git-metrics` provides several sections of output:

1. **Repository information**: Basic metadata about your repository including path, remote URL, and commit history.

2. **Growth statistics**: Year-by-year breakdown of Git object growth (commits, trees, blobs) and disk usage.

3. **Growth projections**: Estimation of future repository growth based on historical trends.

4. **Largest files**: Identification of the largest files in your repository by compressed size.

5. **File extensions**: Analysis of file extensions and their impact on repository size.

### "On-disk size" explained

The on-disk size in `git-metrics`'s output shows the compressed size of commits (saved changes), trees (folder snapshots) and blobs (file versions) as stored in Git's object database (`.git/objects`). These objects are often stored using deltas (storing only changes between similar objects). Repacking the repository (e.g. `git gc`) can alter on-disk sizes of these objects by changing compression and deltas. `git-metrics` does not include the on-disk size of metadata such as pack file indexes (`.git/objects/pack/*.idx`), refs, or other auxiliary files which accounts for 5% to 10% of the overall on-disk size of a repository in most cases.


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

```bash
# Clone the repository
git clone https://github.com/steffen/git-metrics.git
cd git-metrics

# Build the binary
go build
```

After building, you can run the tool as described in the "Running the Application" section.
