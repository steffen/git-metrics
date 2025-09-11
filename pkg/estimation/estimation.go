package estimation

import "git-metrics/pkg/models"

// CalculateEstimate calculates estimated future growth based on current statistics and average growth
func CalculateEstimate(current models.GrowthStatistics, average models.GrowthStatistics) models.GrowthStatistics {
	return models.GrowthStatistics{
		Year:         current.Year + 1,
		Commits:      current.Commits + average.Commits,
		Trees:        current.Trees + average.Trees,
		Blobs:        current.Blobs + average.Blobs,
		Compressed:   current.Compressed + average.Compressed,
		LargestFiles: []models.FileInformation{},
	}
}

// CalculateLinearEstimation calculates linear growth estimation
func CalculateLinearEstimation(current models.GrowthStatistics, average models.GrowthStatistics) models.EstimationResult {
	return models.EstimationResult{
		Method: models.EstimationMethodLinear,
		Statistics: models.GrowthStatistics{
			Year:         current.Year + 1,
			Commits:      current.Commits + average.Commits,
			Trees:        current.Trees + average.Trees,
			Blobs:        current.Blobs + average.Blobs,
			Compressed:   current.Compressed + average.Compressed,
			LargestFiles: []models.FileInformation{},
		},
		FitScore:   0.0,
		GrowthRate: 0.0,
	}
}

// CalculateExponentialEstimation calculates exponential growth estimation
func CalculateExponentialEstimation(current models.GrowthStatistics, yearlyData []models.GrowthStatistics) models.EstimationResult {
	if len(yearlyData) < 2 {
		return models.EstimationResult{
			Method: models.EstimationMethodExponential,
			Statistics: models.GrowthStatistics{
				Year:         current.Year + 1,
				Commits:      current.Commits,
				Trees:        current.Trees,
				Blobs:        current.Blobs,
				Compressed:   current.Compressed,
				LargestFiles: []models.FileInformation{},
			},
			FitScore:   0.0,
			GrowthRate: 0.0,
		}
	}

	commitGrowthRate := calculateExponentialGrowthRate(yearlyData, "commits")
	treeGrowthRate := calculateExponentialGrowthRate(yearlyData, "trees")
	blobGrowthRate := calculateExponentialGrowthRate(yearlyData, "blobs")
	compressedGrowthRate := calculateExponentialGrowthRate(yearlyData, "compressed")

	averageGrowthRate := (commitGrowthRate + treeGrowthRate + blobGrowthRate + compressedGrowthRate) / 4.0

	return models.EstimationResult{
		Method: models.EstimationMethodExponential,
		Statistics: models.GrowthStatistics{
			Year:         current.Year + 1,
			Commits:      int(float64(current.Commits) * (1.0 + commitGrowthRate)),
			Trees:        int(float64(current.Trees) * (1.0 + treeGrowthRate)),
			Blobs:        int(float64(current.Blobs) * (1.0 + blobGrowthRate)),
			Compressed:   int64(float64(current.Compressed) * (1.0 + compressedGrowthRate)),
			LargestFiles: []models.FileInformation{},
		},
		FitScore:   calculateExponentialFitScore(yearlyData),
		GrowthRate: averageGrowthRate,
	}
}

// SelectBestEstimationMethod chooses between linear and exponential based on fit quality
func SelectBestEstimationMethod(current models.GrowthStatistics, average models.GrowthStatistics, yearlyData []models.GrowthStatistics) models.EstimationResult {
	linearResult := CalculateLinearEstimation(current, average)
	exponentialResult := CalculateExponentialEstimation(current, yearlyData)

	linearResult.FitScore = calculateLinearFitScore(yearlyData, average)

	if exponentialResult.FitScore > linearResult.FitScore ||
		(exponentialResult.FitScore >= linearResult.FitScore-0.1 && len(yearlyData) >= 3) {
		return exponentialResult
	}
	return linearResult
}

// calculateExponentialGrowthRate calculates the average annual growth rate for a specific metric
func calculateExponentialGrowthRate(yearlyData []models.GrowthStatistics, metric string) float64 {
	if len(yearlyData) < 2 {
		return 0.0
	}
	var totalGrowthRate float64
	var validPairs int
	for i := 1; i < len(yearlyData); i++ {
		var current, previous float64
		switch metric {
		case "commits":
			current, previous = float64(yearlyData[i].Commits), float64(yearlyData[i-1].Commits)
		case "trees":
			current, previous = float64(yearlyData[i].Trees), float64(yearlyData[i-1].Trees)
		case "blobs":
			current, previous = float64(yearlyData[i].Blobs), float64(yearlyData[i-1].Blobs)
		case "compressed":
			current, previous = float64(yearlyData[i].Compressed), float64(yearlyData[i-1].Compressed)
		}
		if previous > 0 {
			growthRate := (current - previous) / previous
			totalGrowthRate += growthRate
			validPairs++
		}
	}
	if validPairs > 0 {
		return totalGrowthRate / float64(validPairs)
	}
	return 0.0
}

// calculateExponentialFitScore calculates how well exponential model fits the data
func calculateExponentialFitScore(yearlyData []models.GrowthStatistics) float64 {
	if len(yearlyData) < 3 {
		return 0.0
	}
	n := len(yearlyData)
	var sum float64
	for _, data := range yearlyData {
		sum += float64(data.Commits)
	}
	mean := sum / float64(n)
	var ssRes, ssTot float64
	for i, data := range yearlyData {
		observed := float64(data.Commits)
		var predicted float64
		if i > 0 {
			growthRate := calculateExponentialGrowthRate(yearlyData[:i+1], "commits")
			predicted = float64(yearlyData[0].Commits) * (1.0 + growthRate)
			for j := 1; j < i; j++ {
				predicted *= (1.0 + growthRate)
			}
		} else {
			predicted = observed
		}
		ssRes += (observed - predicted) * (observed - predicted)
		ssTot += (observed - mean) * (observed - mean)
	}
	if ssTot == 0 {
		return 1.0
	}
	r2 := 1.0 - (ssRes / ssTot)
	if r2 < 0 {
		return 0.0
	}
	return r2
}

// calculateLinearFitScore calculates how well linear model fits the data
func calculateLinearFitScore(yearlyData []models.GrowthStatistics, average models.GrowthStatistics) float64 {
	if len(yearlyData) < 2 {
		return 0.0
	}
	n := len(yearlyData)
	var sum float64
	for _, data := range yearlyData {
		sum += float64(data.Commits)
	}
	mean := sum / float64(n)
	var ssRes, ssTot float64
	for i, data := range yearlyData {
		observed := float64(data.Commits)
		predicted := float64(yearlyData[0].Commits) + float64(i)*float64(average.Commits)
		ssRes += (observed - predicted) * (observed - predicted)
		ssTot += (observed - mean) * (observed - mean)
	}
	if ssTot == 0 {
		return 1.0
	}
	r2 := 1.0 - (ssRes / ssTot)
	if r2 < 0 {
		return 0.0
	}
	return r2
}
