package handleerr

import "fmt"

type InvalidPercentageErr struct {
	Percentage float64
}

func (e InvalidPercentageErr) Error() string {
	return fmt.Sprintf("Invalid percentage: %f. must be greater than 0", e.Percentage)
}
