package connectors

const (
	connectorCreatedSubQuery = `
		subscription ($%[1]s: IDs, $%[2]s: Boolean!) {
    		connectorCreated(connectorMethodIds: $%[1]s, allTenants: $%[2]s) {
				id
				title
                name
				actions {
					id
					action {
						id
					}
					config
					name
				}
				parameters
				authTypes
				sequence
			}
		}
	`
	connectorUpdatedSubQuery = `
		subscription ($%[1]s: IDs, $%[2]s: Boolean!) {
    		connectorUpdated(connectorMethodIds: $%[1]s, allTenants: $%[2]s){
				id
				title
                name
				actions {
					id
					action {
						id
					}
					config
					name
				}
				parameters
				authTypes
				sequence
			}
		}
	`

	connectorDeletedSubQuery = `
		subscription ($%[1]s: IDs, $%[2]s: Boolean!) {
    		connectorDeleted(connectorMethodIds: $%[1]s, allTenants: $%[2]s) {
				id
				sequence
				name
				description
			}
		}
	`
)
