package playbooks

const (
	getPlaybookQuery = `
		query ($playbookId: ID!) {
		  	playbook(playbookId: $playbookId) {
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
		  	  	sequence
		  	}
		}
	`

	getPlaybooksQuery = `
		query ($categoryId: ID, $tags: Tags) {
			playbooks (categoryId: $categoryId, tags: $tags) {
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
		  	  	sequence
		  	}
		}
	`

	getPlaybookInstanceQuery = `
		query ($playbookInstanceId: ID!) {
		    playbookInstance (playbookInstanceId: $playbookInstanceId) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        playbook {
		            id
					title
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            tags
		            requires {
		                id
		            }
		            categories {
		                id
		            }
		        }
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
		              name
		              description
		              createdAt
		              updatedAt
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
		        sequence
		    }
		}
	`
	getPlaybookInstancesQuery = `
		query($playbookId: ID) {
			playbookInstances (playbookId: $playbookId) {
		        id
		        createdAt
		        updatedAt
		        tenant
		        name
		        description
		        tags
		        playbook {
		            id
					title
		            createdAt
		            updatedAt
		            tenant
		            name
		            description
		            tags
		            requires {
		                id
		            }
		            categories {
		                id
		            }
		        }
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
		              name
		              description
		              createdAt
		              updatedAt
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
		        sequence
		    }
		}
	`
	getPlaybookExecutionQuery = `
		query ($playbookExecutionId: ID!) {
		    playbookExecution (playbookExecutionId: $playbookExecutionId) {
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
		            enabled
		            inputs
		        }
		        state
		        inputs
		        outputs
		    }
		}
	`

	getPlaybookExecutionsQuery = `
		query ($playbookInstanceId: ID!, $pagination: Pagination!) {
		  	playbookExecutions(playbookInstanceId: $playbookInstanceId, pagination: $pagination) {
		    	totalCount
		    	executions {
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
		    	  	  	enabled
		    	  	  	inputs
		    	  	}
		    	  	state
		    	  	inputs
		    	  	outputs
		    	}
		  	}
		}
	`
	getPlaybookTriggerQuery = `
		query($playbookTriggerId: ID!) {
			playbookTrigger (playbookTriggerId: $playbookTriggerId) {
        		id
        		createdAt
        		updatedAt
        		tenant
        		name
        		description
        		instance {
        		    id
        		    createdAt
        		    updatedAt
        		    tenant
        		    name
        		    description
        		    tags
        		    enabled
        		    sequence
        		}
        		type {
        		    id
        		    createdAt
        		    updatedAt
        		    name
        		    description
        		    parameters
        		}
        		config
    		}
		}
	`
	getPlaybookTriggersQuery = `
		query($playbookTriggerTypeIds: IDs!) {
			playbookTriggers (playbookTriggerTypeIds: $playbookTriggerTypeIds) {
        		id
        		createdAt
        		updatedAt
        		tenant
        		name
        		description
        		instance {
        		    id
        		    createdAt
        		    updatedAt
        		    tenant
        		    name
        		    description
        		    tags
        		    enabled
        		    sequence
        		}
        		type {
        		    id
        		    createdAt
        		    updatedAt
        		    name
        		    description
        		    parameters
        		}
        		config
    		}
		}
	`

	getTriggerTypeQuery = `
		query($playbookTriggerTypeId: ID, $playbookTriggerTypeName: String) {
			playbookTriggerType(playbookTriggerTypeId: $playbookTriggerTypeId, playbookTriggerTypeName: $playbookTriggerTypeName) {
				id
				createdAt
				updatedAt
				name
				description
				parameters
			 }
		}
	`

	getTriggerTypesQuery = `
		query {
			playbookTriggerTypes {
				id
				createdAt
				updatedAt
				name
				description
				parameters
			 }
		}
	`
)
