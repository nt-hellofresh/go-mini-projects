package main

type Package struct {
	SKU        string   `json:"box_sku"`
	Items      []string `json:"box_recipes"`
	NrIcePacks int      `json:"ices"`
}

type Details struct {
	Packages []Package `json:"packages"`
}

type PackageResult struct {
	Id      string  `json:"main_package_id"`
	Details Details `json:"output"`
}

type Output map[string]interface{}

// factory method to configure an ETL to pipe input to output
func configurePipe() *ETLPipe[PackageResult, Output] {
	return NewETLPipe(handleTransformLogic)
}

func handleTransformLogic(input PackageResult) (Output, error) {
	// Custom ETL logic. For this example, we just count the number of boxes

	result := Output{
		"box_id":   input.Id,
		"nr_boxes": len(input.Details.Packages),
	}
	return result, nil
}
