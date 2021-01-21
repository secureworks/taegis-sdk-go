package playbooks

const (
	createPlaybookMutation = `
		mutation ($playbook: PlaybookInput!) {
		    createPlaybook (playbook: $playbook) {
		  	  	id
				title
		  	  	createdAt
		  	  	updatedAt
		  	  	tenant
		  	  	name
		  	  	description
		  	  	tags
		  	  	categories {
		  	  	  	id
		  	  	}
		  	  	requires {
		  	  	  	id
		  	  	}
		  	  	head {
		  	  		id
		  	  		createdAt
		  	  		createdBy
		  	  		inputs
		  	  		outputs
		  	  		dsl
		  	  	}
		  	  	versions {
		  	  	  	id
		  	  	  	createdAt
		  	  	  	createdBy
		  	  	  	requires {
		  	  	  	  	id
		  	  	  	}
		  	  	  	inputs
		  	  	  	outputs
		  	  	  	dsl
		  	  	}
		  	}
		}
	`

	clonePlaybookMutation = `
		mutation ($input: ClonePlaybookInput!) {
		    clonePlaybook (input: $input) {
		  	  	id
				title
		  	  	createdAt
		  	  	updatedAt
		  	  	tenant
		  	  	name
		  	  	description
		  	  	tags
		  	  	categories {
		  	  	  	id
		  	  	}
		  	  	requires {
		  	  	  	id
		  	  	}
		  	  	head {
		  	  		id
		  	  		createdAt
		  	  		createdBy
		  	  		inputs
		  	  		outputs
		  	  		dsl
		  	  	}
		  	  	versions {
		  	  	  	id
		  	  	  	createdAt
		  	  	  	createdBy
		  	  	  	requires {
		  	  	  	  	id
		  	  	  	}
		  	  	  	inputs
		  	  	  	outputs
		  	  	  	dsl
		  	  	}
		  	}
		}
	`

	updatePlaybookMutation = `
		mutation ($playbookId: ID!, $playbook: PlaybookInput!) {
		    updatePlaybook (playbookId: $playbookId, playbook: $playbook) {
				title
		  	  	id
		  	  	createdAt
		  	  	updatedAt
		  	  	tenant
		  	  	name
		  	  	description
		  	  	tags
		  	  	categories {
		  	  	  	id
		  	  	}
		  	  	requires {
		  	  	  	id
		  	  	}
		  	  	head {
		  	  		id
		  	  		createdAt
		  	  		createdBy
		  	  		inputs
		  	  		outputs
		  	  		dsl
		  	  	}
		  	  	versions {
		  	  	  	id
		  	  	  	createdAt
		  	  	  	createdBy
		  	  	  	requires {
		  	  	  	  	id
		  	  	  	}
		  	  	  	inputs
		  	  	  	outputs
		  	  	  	dsl
		  	  	}
		  	}
		}
	`

	deletePlaybookMutation = `
		mutation ($playbookId: ID!) {
		    deletePlaybook (playbookId: $playbookId) {
		  	  	id
		  	  	createdAt
		  	  	updatedAt
		  	  	tenant
		  	  	name
		  	  	description
		  	  	tags
		  	  	categories {
		  	  	  	id
		  	  	}
		  	  	requires {
		  	  	  	id
		  	  	}
		  	  	head {
		  	  		id
		  	  		createdAt
		  	  		createdBy
		  	  		inputs
		  	  		outputs
		  	  		dsl
		  	  	}
		  	  	versions {
		  	  	  	id
		  	  	  	createdAt
		  	  	  	createdBy
		  	  	  	requires {
		  	  	  	  	id
		  	  	  	}
		  	  	  	inputs
		  	  	  	outputs
		  	  	  	dsl
		  	  	}
		  	}
		}
	`

	executePlaybookMutation = `
		mutation ($playbookId: ID!, $parameters: JSONObject) {
		    executePlaybook (playbookId: $playbookId, parameters: $parameters) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        createdBy
		        instance {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            tags
		            version {
		                id
		                createdAt
		                createdBy
		                requires {
		                    id
		                }
		                inputs
		                outputs
		                dsl
		            }
		            enabled
		            inputs
					retries {
		                InitialInterval
		                MaximumInterval
		                BackoffCoefficient
		                MaximumRetries
		                MaximumDuration
		            }
		            connections {
		                id
		            }
		        }
		        state
		        inputs
		        outputs
		        events {
		            id
		            object
		            state
		            name
		            timestamp
		            inputs
		            outputs
		            reason
		            attempt
		        }
		    }
		}
	`

	createPlaybookInstanceMutation = `
		mutation ($playbookId: ID!, $instance: PlaybookInstanceInput!) {
		    createPlaybookInstance (playbookId: $playbookId, instance: $instance) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        version {
		            id
		            createdAt
		            createdBy
		            inputs
		            outputs
		            dsl
		        }
		        trigger {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            config
					type {
						id
						createdAt
						updatedAt
						name
						description
						parameters
					}
		        }
		        enabled
		        inputs
		        retries {
		            InitialInterval
		            MaximumInterval
		            BackoffCoefficient
		            MaximumRetries
		            MaximumDuration
		        }
		        connections {
		            id
		        }
		    }
		}
	`

	updatePlaybookInstanceMutation = `
		mutation ($playbookInstanceId: ID!, $instance: PlaybookInstanceInput!) {
		    updatePlaybookInstance (playbookInstanceId: $playbookInstanceId, instance: $instance) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        version {
		            id
		            createdAt
		            createdBy
		            inputs
		            outputs
		            dsl
		        }
		        trigger {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            config
					type {
						id
						createdAt
						updatedAt
						name
						description
						parameters
					}
		        }
		        enabled
		        inputs
		        retries {
		            InitialInterval
		            MaximumInterval
		            BackoffCoefficient
		            MaximumRetries
		            MaximumDuration
		        }
		        connections {
		            id
		        }
		    }
		}
	`

	deletePlaybookInstanceMutation = `
		mutation ($playbookInstanceId: ID!) {
		    deletePlaybookInstance (playbookInstanceId: $playbookInstanceId) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        version {
		            id
		            createdAt
		            createdBy
		            inputs
		            outputs
		            dsl
		        }
		        trigger {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            config
					type {
						id
						createdAt
						updatedAt
						name
						description
						parameters
					}
		        }
		        enabled
		        inputs
		        retries {
		            InitialInterval
		            MaximumInterval
		            BackoffCoefficient
		            MaximumRetries
		            MaximumDuration
		        }
		        connections {
		            id
		        }
		    }
		}
	`

	executePlaybookInstanceMutation = `
		mutation ($playbookInstanceId: ID!, $parameters: JSONObject) {
		    executePlaybookInstance (playbookInstanceId: $playbookInstanceId, parameters: $parameters) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        createdBy
		        state
		        inputs
		        outputs
		        instance {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            tags
		            enabled
		            inputs
		        }
		        events {
		            id
		            object
		            state
		            name
		            timestamp
		            inputs
		            outputs
		            reason
		            attempt
		        }
		    }
		}
	`

	setPlaybookInstanceStateMutation = `
		mutation ($playbookInstanceId: ID!, $enabled: Boolean!) {
		    setPlaybookInstanceState (playbookInstanceId: $playbookInstanceId, enabled: $enabled) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        version {
		            id
		            createdAt
		            createdBy
		            inputs
		            outputs
		            dsl
		        }
		        trigger {
		            id
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            config
					type {
						id
						createdAt
						updatedAt
						name
						description
						parameters
					}
		        }
		        enabled
		        inputs
		        retries {
		            InitialInterval
		            MaximumInterval
		            BackoffCoefficient
		            MaximumRetries
		            MaximumDuration
		        }
		        connections {
		            id
		        }
		    }
		}
	`
)
