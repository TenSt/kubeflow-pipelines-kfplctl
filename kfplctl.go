package main

import (
	"flag"
	"fmt"
	"os"

	sdk "github.com/TenSt/kubeflow-pipelines-sdk"
)

var client sdk.KfPipelineClient

func init() {
	client = sdk.GetClient("http://188.40.161.99:8888")
}

func main() {
	// textPtr := flag.String("text", "default", "Text to parse.")
	// metricPtr := flag.String("metric", "chars", "Metric {chars|words|lines};.")
	// uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")
	flag.Parse()

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
		// fmt.Println(p)
		for _, p := range pls.Pipelines {
			fmt.Println(p.ID)
		}
	case "pipeline":
		p := client.GetPipeline(os.Args[3])
		fmt.Println(p)
	case "experiments":
		e := client.GetAllExperiments()
		fmt.Println(e)
	case "experiment":
		e := client.GetExperiment(os.Args[3])
		fmt.Println(e)
	case "runs":
		r := client.GetAllRuns()
		fmt.Println(r)
	case "run":
		r := client.GetRun(os.Args[3])
		fmt.Println(r)
	}
}

func delete() {
	switch os.Args[2] {
	case "pipeline":
		p := client.DeletePipeline(os.Args[3])
		fmt.Println(p)
	case "experiment":
		e := client.DeleteExperiment(os.Args[3])
		fmt.Println(e)
	}
}

func create() {
	switch os.Args[2] {
	case "experiment":
		e := client.CreateExperiment("", "")
		fmt.Println(e.ID)
	case "run":
		r := sdk.Run{}
		rDetail := client.CreateRun(r)
		fmt.Println(rDetail.Run.ID)
	}
}

func uploadPipeline() {
	p := client.UploadPipeline("", "")
	fmt.Println(p.ID)
}
