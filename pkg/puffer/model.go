package puffer

type Info struct {
	HighTemp      float32
	MidTemp       float32
	LowTemp       float32
	CollectorTemp float32
}

type Options struct {
	Url string
	User string
	Password string
}