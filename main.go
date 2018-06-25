package main

import (
	"fmt"
	"os"

	ansi "github.com/jhunt/go-ansi"
)

func main() {
	var client Client
	err := client.setup()
	if err != nil {
		bailWith("err setting up client: %s", err)
	}
	orgs, err := client.getOrgs()
	if err != nil {
		bailWith("error getting orgs: %s", err)
	}
	fmt.Printf("\n\n\n\n\n")
	fmt.Println("orgs before processing", orgs)
	fmt.Printf("\n\n\n\n\n")

	//associate app creates with orgs
	for index, org := range orgs {
		var response eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.create&q=organization_guid:"+org.guid, &response)
		if err != nil {
			bailWith("error associating app creates with orgs: %s", err)
		}
		resourceList, err := client.cfEventsResourcesFromResponse(response)
		if err != nil {
			bailWith("error getting resources out of api response %s", err)
		}
		orgs[index].associatedAppCreates = resourceList
	}

	//associate app starts with orgs
	for index, org := range orgs {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.start&q=organization_guid:"+org.guid, &returnStruct)
		if err != nil {
			bailWith("error associating app starts with orgs: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error getting resources out of api resp %s", err)
		}
		orgs[index].associatedAppStarts = responseList
	}

	//associate app updates with orgs
	for index, org := range orgs {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.update&q=organization_guid:"+org.guid, &returnStruct)
		if err != nil {
			bailWith("error associating app updates with orgs: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app updates with orgs %s", err)
		}
		orgs[index].associatedAppStarts = responseList
	}

	//associate space creates with orgs
	for index, org := range orgs {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.space.create&q=organization_guid:"+org.guid, &returnStruct)
		if err != nil {
			bailWith("error associating space creates with orgs: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating space creates with orgs %s", err)
		}
		orgs[index].associatedAppStarts = responseList
	}

	//get all apps based on org
	for index, org := range orgs {
		var returnStruct appsAPIResponse
		err := client.cfAppsAPIRequest("/v2/apps?q=organization_guid:"+org.guid, &returnStruct)
		if err != nil {
			bailWith("error associating apps with orgs: %s", err)
		}
		responseList, err := client.cfAppsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating apps with orgs: %s", err)
		}
		orgs[index].apps = responseList
	}

	//some app stuff for later?
	// for index, org := range orgs {
	// 	for index, app := range orgs[index].apps {
	// 		jsonResponse, err := cfAppsAPIRequest(client, "/v2/service_bindings?q=app_guid:"+orgs[index].apps[index].apps.guid)

	// 	}
	// }

	//get all service bindings based on apps by org

	fmt.Printf("\n\n\n\n\n")
	fmt.Println("orgs after data processing:", orgs)
	fmt.Printf("\n\n\n\n\n")

	//grab all the spaces
	spaces, err := client.getSpaces()
	if err != nil {
		bailWith("error getting spaces: %s", err)
	}

	fmt.Printf("\n\n\n\n\n")
	fmt.Printf("spaces before data processing\n")
	fmt.Println(spaces)
	fmt.Printf("\n\n\n\n\n")

	//associate app starts with spaces
	for index, space := range spaces {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.start&q=space_guid:"+space.guid, &returnStruct)
		if err != nil {
			bailWith("error associating app starts with spaces: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app starts with spaces %s", err)
		}
		spaces[index].associatedAppStarts = responseList
	}

	//associate app creates with spaces
	for index, space := range spaces {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.create&q=space_guid:"+space.guid, &returnStruct)
		if err != nil {
			bailWith("error associating app creates with spaces: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app creates with spaces %s", err)
		}
		spaces[index].associatedAppStarts = responseList
	}

	//associate app updates with spaces
	for index, space := range spaces {
		var returnStruct eventsAPIResponse
		err := client.cfEventsAPIRequest("/v2/events?q=type:audit.app.update&q=space_guid:"+space.guid, &returnStruct)
		if err != nil {
			bailWith("error associating app updates with spaces: %s", err)
		}
		responseList, err := client.cfEventsResourcesFromResponse(returnStruct)
		if err != nil {
			bailWith("error associating app updates with spaces %s", err)
		}
		spaces[index].associatedAppStarts = responseList
	}

	fmt.Printf("\n\n\n\n\n")
	fmt.Printf("spaces after data processing\n")
	fmt.Println(spaces)
	fmt.Printf("\n\n\n\n\n")

	//get all apps based on spaces

	// get all service bindings based on apps by space

	// fmt.Println(spaces
	// for {
	// 	serve()
	// }
}

func bailWith(f string, a ...interface{}) {
	ansi.Fprintf(os.Stderr, fmt.Sprintf("@R{%s}\n", f), a...)
	os.Exit(1)
}
