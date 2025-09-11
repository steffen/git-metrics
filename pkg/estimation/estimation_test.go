package estimation

import (
	"git-metrics/pkg/models"
	"testing"
)

func TestCalculateEstimate(t *testing.T) {
	current := models.GrowthStatistics{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 4000}
	average := models.GrowthStatistics{Commits: 100, Trees: 200, Blobs: 300, Compressed: 400}
	expected := models.GrowthStatistics{Year: 2024, Commits: 1100, Trees: 2200, Blobs: 3300, Compressed: 4400, LargestFiles: []models.FileInformation{}}
	result := CalculateEstimate(current, average)
	if result.Year != expected.Year {
		t.Errorf("Year = %v, want %v", result.Year, expected.Year)
	}
	if result.Commits != expected.Commits {
		t.Errorf("Commits = %v, want %v", result.Commits, expected.Commits)
	}
	if result.Trees != expected.Trees {
		t.Errorf("Trees = %v, want %v", result.Trees, expected.Trees)
	}
	if result.Blobs != expected.Blobs {
		t.Errorf("Blobs = %v, want %v", result.Blobs, expected.Blobs)
	}
	if result.Compressed != expected.Compressed {
		t.Errorf("Compressed = %v, want %v", result.Compressed, expected.Compressed)
	}
	if len(result.LargestFiles) != len(expected.LargestFiles) {
		t.Errorf("LargestFiles length = %v, want %v", len(result.LargestFiles), len(expected.LargestFiles))
	}
}

func TestCalculateLinearEstimation(t *testing.T) {
	current := models.GrowthStatistics{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 4000}
	average := models.GrowthStatistics{Commits: 100, Trees: 200, Blobs: 300, Compressed: 400}
	result := CalculateLinearEstimation(current, average)
	if result.Method != models.EstimationMethodLinear {
		t.Errorf("Method = %v, want %v", result.Method, models.EstimationMethodLinear)
	}
	if result.Statistics.Year != 2024 {
		t.Errorf("Year = %v, want %v", result.Statistics.Year, 2024)
	}
	if result.Statistics.Commits != 1100 {
		t.Errorf("Commits = %v, want %v", result.Statistics.Commits, 1100)
	}
}

func TestCalculateExponentialEstimation(t *testing.T) {
	yearlyData := []models.GrowthStatistics{
		{Year: 2021, Commits: 100, Trees: 200, Blobs: 300, Compressed: 1000},
		{Year: 2022, Commits: 120, Trees: 240, Blobs: 360, Compressed: 1200},
		{Year: 2023, Commits: 144, Trees: 288, Blobs: 432, Compressed: 1440},
	}
	current := yearlyData[len(yearlyData)-1]
	result := CalculateExponentialEstimation(current, yearlyData)
	if result.Method != models.EstimationMethodExponential {
		t.Errorf("Method = %v, want %v", result.Method, models.EstimationMethodExponential)
	}
	if result.Statistics.Year != 2024 {
		t.Errorf("Year = %v, want %v", result.Statistics.Year, 2024)
	}
	if result.Statistics.Commits < 160 || result.Statistics.Commits > 180 {
		t.Errorf("Commits = %v, expected around 173", result.Statistics.Commits)
	}
}

func TestSelectBestEstimationMethod(t *testing.T) {
	current := models.GrowthStatistics{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 4000}
	average := models.GrowthStatistics{Commits: 100, Trees: 200, Blobs: 300, Compressed: 400}
	yearlyData := []models.GrowthStatistics{{Year: 2023, Commits: 1000, Trees: 2000, Blobs: 3000, Compressed: 4000}}
	result := SelectBestEstimationMethod(current, average, yearlyData)
	if result.Method != models.EstimationMethodLinear {
		t.Errorf("Method = %v, want %v", result.Method, models.EstimationMethodLinear)
	}
	yearlyDataMore := []models.GrowthStatistics{
		{Year: 2021, Commits: 100, Trees: 200, Blobs: 300, Compressed: 1000},
		{Year: 2022, Commits: 200, Trees: 400, Blobs: 600, Compressed: 2000},
		{Year: 2023, Commits: 300, Trees: 600, Blobs: 900, Compressed: 3000},
	}
	result2 := SelectBestEstimationMethod(current, average, yearlyDataMore)
	if result2.Method != models.EstimationMethodLinear && result2.Method != models.EstimationMethodExponential {
		t.Errorf("Returned invalid method: %v", result2.Method)
	}
}
