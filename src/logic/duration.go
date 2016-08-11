//dsp video duration directional
package logic

type Duration struct {
	Min int
	Max int
}

func (d Duration) UnderDuration(duration int) bool {
	return (d.Min == 0 || duration >= d.Min) && (d.Max == 0 || duration < d.Max)
}
