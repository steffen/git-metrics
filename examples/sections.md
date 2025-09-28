# RUN
## Run Metadata
General metadata about when and where the report was generated (start time, host machine, versions). Useful to contextualize measurements and reproduce runs.

# REPOSITORY
## Repository Info
Origin information: repository path, remote URL, most recent commit and age. Helps anchor the report to a specific revision set.

# HISTORIC AND ESTIMATED GROWTH
## Historic and Estimated Growth
Shows yearly totals of core Git objects (commits, trees, blobs) along with on-disk size. Past years are actual; rows with ^ are current totals; rows with * are projections extrapolated from recent growth.

# RATE OF CHANGES
## Rate of Changes
Focuses on commit cadence to the default branch. P95/P99/P100 peaks per day/hour/minute reveal burstiness and scaling of integration workflow.

# LARGEST DIRECTORIES
## Largest Directories
Identifies directories contributing ≥1% of repository storage. Highlights translation files, tests, docs, and core source areas for optimization or pruning.

# LARGEST FILES
## Largest Files
Top individual files by cumulative blob storage, signaling hotspots for size bloat and potential candidates for history rewriting or splitting.

# LARGEST FILE EXTENSIONS
## Largest File Extensions
Distribution of blob count and size by file extension. Useful to see language / artifact composition and track shifts over time.

# AUTHORS WITH MOST COMMITS
## Authors With Most Commits
Per-year top authors by authored commits plus totals. Shows contributor concentration and evolution of community participation.

# COMMITTERS WITH MOST COMMITS
## Committers With Most Commits
Committer stats (who integrated patches). High centralization can indicate a gatekeeping pattern or strong maintainer oversight.

# FOOTER
## Footer / Summary
Runtime performance of the metrics tool itself (execution time, memory footprint).
