package api

// FetchPufferData is for getting the latest puffer data
func FetchPufferData() *PufferInfo {
	// Dummy info for now ...
	return &PufferInfo{
		HighTemp: 55.5,
		MidTemp:  35.4,
		LowTemp:  23,
	}
}
