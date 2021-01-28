package events

const (
	getEvents = `
		query($%[1]s: [ID!]!) {
		  	events(ids: $%[1]s) {
		  	  	id
				values
		  	}
		}
	`
	getEventQuery = `
		query($%[1]s: ID!) {
		  	eventQuery(id: $%[1]s) {
				id,
				query,
				status,
                reasons {
                    id,
                    type,
                    backend,
                    status,
                    reason,
                    submitted,
                    completed
                },
				submitted,
				completed,
				expires,
				types,
				metadata
		  	}
		}
	`
	getEventQueries = `
		query {
		  	eventQueries {
				id,
				query
				status
                reasons {
                    id,
                    type,
                    backend,
                    status,
                    reason,
                    submitted,
                    completed
                },
				submitted
				completed
				expires
				types
				metadata
		  	}
		}
	`
	deleteEventQuery = `
		mutation($%[1]s: ID!) {
		  	deleteEventQuery(id: $%[1]s)
		}
	`
)
