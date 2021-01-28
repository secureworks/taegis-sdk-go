package connectors

const (
	connectionMethodQuery = `
		query ($%[1]s: String!) {
			connectionMethod(connectionMethodName: $%[1]s) {
				id
				name
				description
				parameters
			 }
		}
	`
	connectorsQuery = `
		query ($%[1]s: IDs, $%[2]s: IDs, $%[3]s: IDs, $%[4]s: Tags) {
			connectors(connectionMethodIds: $%[1]s, connectorInterfaceIds: $%[2]s, connectorCategoryIds: $%[3]s, tags: $%[4]s) {
				id
                name
				title
				method {
					id
					name
            		description
					parameters
				}
                implements {
					id
					name
					description
					actions {
						id
						name
						description
						inputs
						outputs
					}
                }
				actions {
					id
					action {
						id
						name
						description
						inputs
						outputs
					}
					config
					name
				}
				parameters
				authTypes
				sequence
				tenant
			 }
		}
	`
	// TODO: update the connections query to use the GetConnectionsInput
	connectionsQuery = `
		query ($%[1]s: IDs, $%[2]s: IDs, $%[3]s: IDs) {
			connections(connectionIds: $%[1]s, connectorIds: $%[2]s, connectorInterfaceIds: $%[3]s) {
				id
                name
				connector {
					id
					title
					name
					description
                    categories {
						id
						name
						description
					}
					method {
						id
						name
						description
						parameters
					}
					implements {
						id
						name
						description
						actions {
							id
							name
							description
							inputs
							outputs
						}
					}
					actions {
						id
						action {
							id
							name
							description
							inputs
							outputs
						}
						name
					}
					parameters
					authTypes
					sequence
					tenant
				}
				authType
				authUrl
				config
				credentials
				sequence
			 }
		}
	`
	getContextQuery = `
		query ($%[1]s: IDs, $%[2]s: IDs, $%[3]s: IDs) {
			connections(connectionIds: $%[1]s, connectorIds: $%[2]s, connectorInterfaceIds: $%[3]s) {
				id
                name
				connector {
					id
					title
					name
					method {
						id
						name
					}
					implements {
						id
						name
						actions {
							id
							name
							inputs
							outputs
						}
					}
					actions {
						id
						action {
							id
							name
							interface {
								id
								name
							}
						}
					}
					parameters
					tenant
				}
				authType
				authUrl
				config
				credentials
			 }
		}
	`
	connectorCategoryQuery = `
	query ($%[1]s: ID!) {
			connectorCategory(connectorCategoryId: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
	}
`
	connectorInterfaceQuery = `
	query ($%[1]s: ID!) {
			connectorInterface(connectorInterfaceId: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				tenantId
				categories
				actions {
					id
					name
				}
	}
`
	connectorQuery = `
	query ($%[1]s: ID!) {
			connector(connectorId: $%[1]s) {
				id
                name
				title
				method {
					id
					name
            		description
					parameters
				}
                implements {
					id
					name
					description
					actions {
						id
						name
						description
						inputs
						outputs
					}
                }
				actions {
					id
					action {
						id
						name
						description
						inputs
						outputs
					}
					config
					name
				}
				sequence
				parameters
				authTypes
				tenant
			 }
	}
`

	connectionQuery = `
	query ($%[1]s: ID!) {
		connection(connectionId: $%[1]s) {
			id
			name
			connector {
				id
				name
				title
				method {
					id
					name
				}
				implements {
					id
					name
					actions {
						id
						name
						inputs
						outputs
					}
				}
				actions {
					id
					name
					action {
						id
						name
						interface {
							id
							name
						}
					}
				}
				parameters
				sequence
			}
			authType
			authUrl
			config
			credentials
			sequence
		 }
	}`
)
