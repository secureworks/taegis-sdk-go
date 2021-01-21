package connectorLogger

const (
	getAllConnectorLogsQuery = `
		query($args: ConnectorLogQueryInput, $pagination: Pagination){
			getAllConnectorLogs(args: $args, pagination: $pagination){
				totalCount
				entries {
					id
					connector
					message
					raw_error
					created_at
				}
			}
		}
	`
)
