package plotberry

import (
	"net/http"
	"testing"
)

// Tests whether the Plotly API is still functioning
func TestPlotlyAPI(t *testing.T) {

	resp, err := http.Get("https://plot.ly/v0/plotberries")
	if err != nil {
		t.Errorf("Connection error: %s", err)
	}
	if resp.StatusCode/100 != 2 {
		t.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
}

func TestGetPlotberry(t *testing.T) {

	con, _ := GetPlotberry()

	if con == nil {
		t.Errorf("Invalid plotberries content")
	}
}
