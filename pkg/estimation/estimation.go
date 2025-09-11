package estimation

import (
	"fmt"
	"math"
	"strings"
	"git-metrics/pkg/models"
)

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

// CompareModels returns the best estimation result along with both model fit scores.
// This allows callers to display a comparison like:
// "Linear fit score is x and exponential fit score is y, using <model> model."
func CompareModels(current models.GrowthStatistics, average models.GrowthStatistics, yearlyData []models.GrowthStatistics) (models.EstimationResult, float64, float64) {
	linearResult := CalculateLinearEstimation(current, average)
	exponentialResult := CalculateExponentialEstimation(current, yearlyData)

	linearResult.FitScore = calculateLinearFitScore(yearlyData, average)
	linearFit := linearResult.FitScore
	exponentialFit := exponentialResult.FitScore

	var best models.EstimationResult
	if exponentialResult.FitScore > linearResult.FitScore ||
		(exponentialResult.FitScore >= linearResult.FitScore-0.1 && len(yearlyData) >= 3) {
		best = exponentialResult
	} else {
		best = linearResult
	}

	return best, linearFit, exponentialFit
}

// GenerateFitScoreDebug builds human-readable explanations for linear and exponential fit scores (commits only)
// Returned strings may span multiple lines and are intended for debug output.
func GenerateFitScoreDebug(yearlyData []models.GrowthStatistics, average models.GrowthStatistics) (string, string) {
	if len(yearlyData) == 0 {
		return "[debug] No yearly data available for fit score calculation\n", "[debug] No yearly data available for fit score calculation\n"
	}

	linearBuilder := &strings.Builder{}
	expBuilder := &strings.Builder{}

	// Shared commit series
	fmt.Fprintf(linearBuilder, "[debug] Years & commits: ")
	fmt.Fprintf(expBuilder, "[debug] Years & commits: ")
	for i, d := range yearlyData {
		if i > 0 {
			linearBuilder.WriteString(", ")
			expBuilder.WriteString(", ")
		}
		fmt.Fprintf(linearBuilder, "%d=%d", d.Year, d.Commits)
		fmt.Fprintf(expBuilder, "%d=%d", d.Year, d.Commits)
	}
	linearBuilder.WriteString("\n")
	expBuilder.WriteString("\n")

	// Linear model explanation
	if len(yearlyData) < 2 {
		linearBuilder.WriteString("[debug] Linear: insufficient data (<2 years) => fit score 0.00\n")
	} else {
		avgInc := average.Commits
		fmt.Fprintf(linearBuilder, "[debug] Linear: average increment (Δ commits / year) = %d\n", avgInc)
		// Mean
		var sum float64
		for _, d := range yearlyData { sum += float64(d.Commits) }
		mean := sum / float64(len(yearlyData))
		fmt.Fprintf(linearBuilder, "[debug] Linear: mean commits = %.2f\n", mean)
		// Predicted & residuals
		ssRes := 0.0
		ssTot := 0.0
		linearBuilder.WriteString("[debug] Linear: predicted sequence: ")
		for i, d := range yearlyData {
			pred := float64(yearlyData[0].Commits) + float64(i*avgInc)
			if i > 0 { linearBuilder.WriteString(", ") }
			fmt.Fprintf(linearBuilder, "%d→%.0f", d.Year, pred)
			obs := float64(d.Commits)
			ssRes += (obs - pred) * (obs - pred)
			ssTot += (obs - mean) * (obs - mean)
		}
		linearBuilder.WriteString("\n")
		var r2 float64
		if ssTot == 0 {
			r2 = 1.0
		} else {
			r2 = 1.0 - ssRes/ssTot
			if r2 < 0 { r2 = 0 }
		}
		fmt.Fprintf(linearBuilder, "[debug] Linear: SS_res=%.2f SS_tot=%.2f R²=%.4f\n", ssRes, ssTot, r2)
	}

	// Exponential model explanation
	if len(yearlyData) < 3 { // fit function returns 0 when <3
		expBuilder.WriteString("[debug] Exponential: insufficient data (<3 years) => fit score 0.00\n")
	} else {
		// Pair growth rates
		var pairRates []float64
		for i := 1; i < len(yearlyData); i++ {
			prev := yearlyData[i-1].Commits
			curr := yearlyData[i].Commits
			if prev > 0 {
				pairRates = append(pairRates, float64(curr-prev)/float64(prev))
			}
		}
		fmt.Fprintf(expBuilder, "[debug] Exponential: pair growth rates = [")
		for i, r := range pairRates {
			if i > 0 { expBuilder.WriteString(", ") }
			fmt.Fprintf(expBuilder, "%0.4f", r)
		}
		expBuilder.WriteString("]\n")
		avgRate := calculateExponentialGrowthRate(yearlyData, "commits")
		fmt.Fprintf(expBuilder, "[debug] Exponential: average growth rate = %0.4f (%.2f%%)\n", avgRate, avgRate*100)
		// Mean
		var sum float64
		for _, d := range yearlyData { sum += float64(d.Commits) }
		mean := sum / float64(len(yearlyData))
		fmt.Fprintf(expBuilder, "[debug] Exponential: mean commits = %.2f\n", mean)
		// Predicted & residuals replicating fit logic (recomputes growth for prefix)
		ssRes := 0.0
		ssTot := 0.0
		expBuilder.WriteString("[debug] Exponential: predicted sequence: ")
		first := float64(yearlyData[0].Commits)
		for i, d := range yearlyData {
			var pred float64
			if i == 0 {
				pred = first
			} else {
				g := calculateExponentialGrowthRate(yearlyData[:i+1], "commits")
				pred = first * math.Pow(1.0+g, float64(i))
			}
			if i > 0 { expBuilder.WriteString(", ") }
			fmt.Fprintf(expBuilder, "%d→%.0f", d.Year, pred)
			obs := float64(d.Commits)
			ssRes += (obs - pred) * (obs - pred)
			ssTot += (obs - mean) * (obs - mean)
		}
		expBuilder.WriteString("\n")
		var r2 float64
		if ssTot == 0 {
			r2 = 1.0
		} else {
			r2 = 1.0 - ssRes/ssTot
			if r2 < 0 { r2 = 0 }
		}
		fmt.Fprintf(expBuilder, "[debug] Exponential: SS_res=%.2f SS_tot=%.2f R²=%.4f\n", ssRes, ssTot, r2)
	}

	linearBuilder.WriteString("[debug] Note: Fit scores use commits only.\n")
	expBuilder.WriteString("[debug] Note: Fit scores use commits only.\n")

	return linearBuilder.String(), expBuilder.String()
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
