//What I am trying to achieve here : a simple tag-based AWS ECS admin tool
//Just focused on what I need

//This tool will be considered feature-full when it can :
// - List exiting instances on any region, their name, Owner tag and their status
// - Search instances by ID or Tags
// - Start and Stop command AGAINST ALL INSTANCES MATCHING THE SEARCH FILTER
// - Add tags or retags command AGAINST ALL INSTANCES MATCHING THE SEARCH FILTER

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//ran once
func main() {
	parseInput()
	svc := EC2_init(region)

	searchInstances(svc)
	displayResults()
	commandInstances(svc)

}

// CLI arguments
// CLAWS expect any of the following option as arguments
// Region: -r REGION      : could be eu-west-1 or ireland format
// Tags:   -t KEY:VALUE   : Tag-based search, wildcard can be used
// Command: -c start / -c stop /-c tag Owner:Batman : add or edit existing tag
func parseInput() {
	//process the arguments
	flag.Parse()
	command_param = flag.Args()

	if tmp, exist := regions[*regionTag]; exist {
		region = tmp
	} else {
		region = *regionTag
	}
	parsed_arg := strings.Split(*Tag, ":")
	keyTag, valueTag = parsed_arg[0], parsed_arg[1]
	if *instanceId != "" && keyTag != "" {
		instance_mode = "true"
	}
}

func searchInstances(svc *ec2.EC2) {
	var inputFilter *ec2.DescribeInstancesInput
	//define inputs (filters)
	//The type DescribeInstancesInput obviously contains the inputs (here filter)
	//for the DescribeInstances below
	switch instance_mode {
	case "false":
		inputFilter = addTagFilter(keyTag, valueTag, inputFilter)
	case "true":
		inputFilter = addInstanceIDFilter(*instanceId, inputFilter)
	}
	//call
	resp, err := svc.DescribeInstances(inputFilter)
	if err != nil {
		panic(err)
	}
	indexResult(resp)
}

//As each command ran has different config/syntax, below is the parsing.
//For each instance found in the tag-based search
//Build AWS object (ionput/output)
//Create resources
func commandInstances(svc *ec2.EC2) {
	//loop through all instances
	for _, iid := range iids {
		// start commamd
		if *command == "start" {
			fmt.Println("starting- ", iid)
			input := &ec2.StartInstancesInput{
				InstanceIds: []*string{
					aws.String(iid),
				},
				DryRun: aws.Bool(false),
			}
			_, err := svc.StartInstances(input)
			if err != nil {
				fmt.Println(err)
			}
			//stop command
		} else if *command == "stop" {
			fmt.Println("stopping - ", iid)
			input := &ec2.StopInstancesInput{
				InstanceIds: []*string{
					aws.String(iid),
				},
				DryRun: aws.Bool(false),
			}
			_, err := svc.StopInstances(input)
			if err != nil {
				fmt.Println(err)
			}
			//retag command
		} else if *command == "tag" {
			parseNewTags(svc, iid)
		}
	}
}
