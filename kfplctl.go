package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	sdk "github.com/TenSt/kubeflow-pipelines-sdk"
)

var client sdk.KfPipelineClient

func init() {
	client = sdk.GetClient("http://188.40.161.51:8888")
}

func main() {
	switch os.Args[1] {
	case "get":
		get()
	case "delete":
		delete()
	case "create":
		create()
	case "upload":
		uploadPipeline()
	}
}

func get() {
	switch os.Args[2] {
	case "pipelines":
		pls := client.GetAllPipelines()
		for _, p := range pls.Pipelines {
			fmt.Println(p)
		}
	case "pipeline":
		subCommand := flag.NewFlagSet("pipeline", flag.ExitOnError)
		namePtr := subCommand.String("name", "default", "Name of the pipeline.")
		subCommand.Parse(os.Args[3:])
		if *namePtr != "default" {
			pls := client.GetAllPipelines()
			for _, p := range pls.Pipelines {
				if p.Name == *namePtr {
					fmt.Println(p.ID)
				}
			}
		} else {
			p := client.GetPipeline(os.Args[3])
			if p.ID != "" {
				fmt.Println(p)
			}
		}
	case "experiments":
		es := client.GetAllExperiments()
		fmt.Println(es)
	case "experiment":
		subCommand := flag.NewFlagSet("experiment", flag.ExitOnError)
		namePtr := subCommand.String("name", "default", "Name of the pipeline.")
		subCommand.Parse(os.Args[3:])
		if *namePtr != "default" {
			es := client.GetAllExperiments()
			for _, e := range es.Experiments {
				if e.Name == *namePtr {
					fmt.Println(e.ID)
				}
			}
		} else {
			e := client.GetExperiment(os.Args[3])
			if e.ID != "" {
				fmt.Println(e)
			}
		}
	case "runs":
		r := client.GetAllRuns()
		fmt.Println(r)
	case "run":
		r := client.GetRun(os.Args[3])
		if os.Args[4] == "status" {
			fmt.Println(r.Run.Status)
		} else {
			fmt.Println(r)
		}
	}
}

func delete() {
	switch os.Args[2] {
	case "pipeline":
		err := client.DeletePipeline(os.Args[3])
		if err != nil {
			fmt.Println(err)
		}
	case "experiment":
		err := client.DeleteExperiment(os.Args[3])
		if err != nil {
			fmt.Println(err)
		}
	}
}

func create() {
	switch os.Args[2] {
	case "experiment":
		subCommand := flag.NewFlagSet("experiment", flag.ExitOnError)
		descPtr := subCommand.String("desc", "", "Description.")
		subCommand.Parse(os.Args[4:])
		e := client.CreateExperiment(os.Args[3], *descPtr)
		fmt.Println(e.ID)
	case "run":
		subCommand := flag.NewFlagSet("run", flag.ExitOnError)
		fileParamsPtr := subCommand.String("parameters", "", "Filename.")
		descPtr := subCommand.String("desc", "", "Description.")
		plIDPtr := subCommand.String("pipeline-id", "", "Pipeline ID.")
		eIDPtr := subCommand.String("experiment-id", "", "Experiment ID.")
		subCommand.Parse(os.Args[4:])
		jsonFile, err := os.Open(*fileParamsPtr)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
		}

		var params *[]sdk.Parameter
		err = json.Unmarshal(byteValue, &params)
		if err != nil {
			fmt.Println(err)
		}

		r := sdk.Run{
			Name:        os.Args[3],
			Description: *descPtr,
			PipelineSpec: sdk.PipelineSpec{
				PipelineID: *plIDPtr,
				Parameters: *params,
			},
			ResourceReferences: []sdk.ResourceReference{
				{
					Key: sdk.ResourceKey{
						ID:   *eIDPtr,
						Type: "EXPERIMENT",
					},
					Relationship: "OWNER",
				},
			},
		}
		// fmt.Println(r)
		rDetail := client.CreateRun(r)
		fmt.Println(rDetail.Run.ID)
	}
}

func uploadPipeline() {
	subCommand := flag.NewFlagSet("upload", flag.ExitOnError)
	plPtr := subCommand.String("pipeline", "", "Filename.")
	subCommand.Parse(os.Args[3:])
	p := client.UploadPipeline(*plPtr, os.Args[2])
	fmt.Println(p.ID)
}
