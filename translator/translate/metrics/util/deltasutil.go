package util

const (
	Report_deltas_Key            = "report_deltas"
	Tags_Key                     = "tags"
	True_value                   = "true"
	Ignored_fields_for_delta     = "iops_in_progress"
	Ignored_fields_for_delta_Key = "ignored_fields_for_delta"
)

func addReportDeltasTag(inputMap map[string]interface{}, result map[string]interface{}) bool {
	reportDelta := true //default to be true if not specified
	if val, ok := inputMap[Report_deltas_Key]; ok {
		reportDelta = val.(bool)
	}

	if !reportDelta {
		return false
	}

	if result[Tags_Key] == nil {
		result[Tags_Key] = map[string]interface{}{}
	}

	tagsMap := result[Tags_Key].(map[string]interface{})
	tagsMap[Report_deltas_Key] = True_value
	return true
}

func ProcessReportDeltasForDiskIO(input interface{}, result map[string]interface{}) {
	m := input.(map[string]interface{})
	if !addReportDeltasTag(m, result) {
		return
	}
	//add the ignored field to tags if it exists
	if val, ok := m[Measurement_Key]; ok {
		fieldsPass := val.([]interface{})
		for _, field := range fieldsPass {
			if field.(string) == Ignored_fields_for_delta {
				tagsMap := result[Tags_Key].(map[string]interface{})
				tagsMap[Ignored_fields_for_delta_Key] = field
				break
			}
		}
	}
}

func ProcessReportDeltasForNet(input interface{}, result map[string]interface{}) {
	m := input.(map[string]interface{})
	addReportDeltasTag(m, result)
}