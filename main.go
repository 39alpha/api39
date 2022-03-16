package main

import (
	"flag"
	"fmt"
	"github.com/39alpha/api39/api39"
	"github.com/39alpha/api39/api39/donate"
	"github.com/39alpha/api39/api39/site"
	"github.com/39alpha/api39/api39/survey"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"log"
	"os"
)

const apikeylen = 32

var (
	port       = 3964
	genconf    = false
	configpath = ""
)

func init() {
	flag.IntVar(&port, "port", port, "port on which the server will listen")
	flag.BoolVar(&genconf, "genconf", genconf, "generate and print a configuration file to STDOUT and exit")
	flag.StringVar(&configpath, "config", configpath, "path to configuration file (required)")
}

func main() {
	flag.Parse()

	if genconf {
		err := api39.GenerateConfig(apikeylen)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if configpath == "" {
			fmt.Fprintf(os.Stderr, "Error: -config flag is required\n\n")
			flag.Usage()
			os.Exit(1)
		}

		app := iris.New()

		app.UseRouter(recover.New())

		if withConfig, err := api39.NewWithConfig(configpath); err != nil {
			log.Fatal(err)
		} else {
			app.UseGlobal(withConfig)
		}

		v0 := app.Party("/api/v0", api39.RecordBody)
		{
			v0.Post("/site/update", api39.VerifyGithubSignature, api39.ParseBody, site.Update)
			v0.Get("/ipfs/addr", api39.Addr)
			donations := v0.Party("/donate")
			{
				donations.Post("/checkout", api39.ParseBody, donate.Checkout)
			}
			surveys := v0.Party("/surveys", survey.OpenDatabase)
			{
				surveys.Get("/", survey.ListSurveys)
				surveys.Get("/{id}", survey.GetSurvey)
				surveys.Post("/{id}/responses", api39.ParseBody, survey.AddSurveyResponses)
			}
		}

		app.Listen(fmt.Sprintf(":%d", port))
	}
}
