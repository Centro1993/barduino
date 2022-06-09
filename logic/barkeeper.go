package logic

const SERVING_SIZE_IN_ML uint = 300

type DrinkStatus struct{
	ProgressInPercent	uint
	Served				bool
	Canceled			bool
	Interrupted			bool
}