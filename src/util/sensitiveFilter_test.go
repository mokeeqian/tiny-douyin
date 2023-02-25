package util

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	InitFilter()

	content := "AV 123 色情"
	contentFiltered := Filter.Replace(content, '*')
	fmt.Println(contentFiltered)
}
