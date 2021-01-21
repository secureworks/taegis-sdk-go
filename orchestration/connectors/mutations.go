package connectors

const (
	defineConnectionMethodMutation = `
		mutation ($%[1]s: ConnectionMethodInput!) {
			defineConnectionMethod(connectionMethod: $%[1]s) {
				id
				name
				description
				parameters
			}
		}
	`

	removeConnectionMethodMutation = `
		mutation ($%[1]s: ID!) {
			removeConnectionMethod(connectorMethodId: $%[1]s) {
				id
				name
				description
				parameters
			}
		}
	`

	createConnectorInterfaceMutation = `
		mutation ($%[1]s: ConnectorInterfaceInput!) {
			createConnectorInterface(connectorInterface: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				actions {
					id
					createdAt
					updatedAt
					name
					description
					inputs
					outputs
				}
			}
		}
	`
	updateConnectorInterfaceMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectorInterfaceInput!) {
			updateConnectorInterface(connectorInterfaceId: $%[1]s, connectorInterface: $%[2]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				tenantId
				actions {
					id
					createdAt
					updatedAt
					name
					description
					inputs
					outputs
				}
			}
		}
	`

	deleteConnectorInterfaceMutation = `
		mutation ($%[1]s: ID!) {
			deleteConnectorInterface(connectorInterfaceId: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
			}
		}
	`

	createConnectorMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectorInput!) {
			createConnector(connectionMethodId: $%[1]s, connector: $%[2]s) {
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
				}
				parameters
				authTypes
				tenant
			}
		}
	`
	updateConnectorMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectorUpdateInput!) {
			updateConnector(connectorId: $%[1]s, connector: $%[2]s) {
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
				}
				parameters
				authTypes
				tenant
			}
		}
	`

	deleteConnectorMutation = `
		mutation ($%[1]s: ID!) {
			deleteConnector(connectorId: $%[1]s) {
				id
                name
				method {
					id
					name
            		description
					parameters
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
				}
				parameters
				authTypes
				tenant
			}
		}
	`

	createConnectionMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectionInput!) {
			createConnection(connectorId: $%[1]s, connection: $%[2]s) {
				id
                name
				connector {
					id
					name
					title
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
					}
					parameters
					authTypes
					tenant
				}
				authType
				authUrl
				config
				credentials
			}
		}
	`

	updateConnectionMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectionInput!) {
			updateConnection(connectionId: $%[1]s, connection: $%[2]s) {
				id
                name
				connector {
					id
					name
					title
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
					}
					parameters
					authTypes
					tenant
				}
				authType
				authUrl
				config
				credentials
			}
		}
	`

	deleteConnectionMutation = `
		mutation ($%[1]s: ID!) {
			deleteConnection(connectionId: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				authType
				authUrl
				config
			}
		}
	`

	validateConnectionMutation = `
		mutation ($%[1]s: ID!) {
			validateConnection(connectionId: $%[1]s) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				authType
				authUrl
				config
			}
		}
	`

	validateConnectionInputMutation = `
		mutation ($%[1]s: ID!, $%[2]s: ConnectionInput!) {
			validateConnectionInput(connectionId: $%[1]s, connection: $%[2]s ) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				docs
				parameters
				authTypes
			}
		}
	`

	executeConnectionActionMutation = `
		mutation ($%[1]s: ID!, $%[1]s: ConnectionInput!, $%[2]s: String!, $%[3]s: Any) {
			executeConnectionAction(connectionId: $%[1]s, actionName: $%[2]s, inputs: $%[3]s ) {
				id
				createdAt
				updatedAt
				name
				description
				tags
				docs
				parameters
				authTypes
			}
		}
	`
)
