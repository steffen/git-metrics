# RUN AND REPOSITORY
## Run & repository information
Origin and execution context: when and where the metrics were generated (start time, host machine, tool versions) together with repository identity (local path, remote URL, most recent commit hash, commit date, repository age). This anchors the entire report to a reproducible environment and revision so later comparisons or audits know exactly which code and tooling produced subsequent sections.

# HISTORIC AND ESTIMATED GROWTH
## Historic & estimated growth
Shows yearly totals of core Git objects (commits, trees, blobs) along with on-disk size. Past years are actual; rows with ^ are current totals; rows with * are projections extrapolated from recent growth.

# LARGEST FILE EXTENSIONS
## Largest file extensions
Distribution of blob count and size by file extension. Useful to see language / artifact composition and track shifts over time.

# LARGEST DIRECTORIES
## Largest directories
Identifies directories contributing ≥1% of repository storage. Highlights translation files, tests, docs, and core source areas for optimization or pruning.

# LARGEST FILES
## Largest files
Top individual files by cumulative blob storage, signaling hotspots for size bloat and potential candidates for history rewriting or splitting.

# RATE OF CHANGES
## Rate of changes
Focuses on commit cadence to the default branch. P95/P99/P100 peaks per day/hour/minute reveal burstiness and scaling of integration workflow.

# AUTHORS WITH MOST COMMITS
## Authors with most commits
Per-year top authors by authored commits plus totals. Shows contributor concentration and evolution of community participation.

# COMMITTERS WITH MOST COMMITS
## Committers with most commits
Committer stats (who integrated patches). High centralization can indicate a gatekeeping pattern or strong maintainer oversight.

# FOOTER
## Run summary
Runtime performance of the metrics tool itself (execution time, memory footprint).
