package playbooks

const (
	playbookInstanceCreateSubQuery = `
		subscription ($%[1]s: IDs) {
		    playbookInstanceCreated(playbookIds: $%[1]s){
				id
				name
		        enabled
				sequence
				trigger {
					id
					config
					tenant
					type {
						id
						name
					}
				}		
			}
		}
	`

	playbookInstanceUpdateSubQuery = `
		subscription ($%[1]s: IDs) {
		    playbookInstanceUpdated(playbookIds: $%[1]s) {
				id
				name
				enabled
				sequence
				trigger {
					id
					config
					tenant
					type {
						id
						name
					}
				}	
			}
		}
	`

	playbookInstanceDeleteSubQuery = `
		subscription ($%[1]s: IDs) {
		    playbookInstanceDeleted(playbookIds: $%[1]s) {
				id
				sequence
			}
		}
	`
)
