package uni

func Float32ToFloat64(float32Slice []float32) []float64 {
	float64Slice := make([]float64, len(float32Slice))
	for i, v := range float32Slice {
		float64Slice[i] = float64(v)
	}
	return float64Slice
}

func Float64ToFloat32(float64Slice []float64) []float32 {
	float32Slice := make([]float32, len(float64Slice))
	for i, v := range float64Slice {
		float32Slice[i] = float32(v)
	}
	return float32Slice
}
