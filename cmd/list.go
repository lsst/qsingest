/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		w := new(tabwriter.Writer)
		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		podName := "qserv-repl-ctl-0"
		// pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %s in namespace %s: %v\n",
		// 		pod, namespace, statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		// }

		clientset, config := setKubeClient()

		curlcmd := "curl http://qserv-repl-ctl:8080/replication/config -X GET  -H \"Content-Type: application/json\" -d \"{\\\"auth_key\\\":\\\"$PASSWORD\\\"}\""
		log.Info().Str("Pod", podName).Str("Cmd", curlcmd).Msg("Launch command inside container")
		outbuf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		err := ExecCmd(clientset, config, podName, namespace, curlcmd, outbuf, errbuf)
		if err != nil {
			panic(err.Error())
		}
		response := outbuf.Bytes()
		if !json.Valid([]byte(response)) {
			// handle the error here
			log.Fatal().Msg("invalid JSON string")
		}
		var responseJson map[string]any
		json.Unmarshal(response, &responseJson)
		//fmt.Fprint(w, string(data))
		log.Debug().Bytes("Error message", errbuf.Bytes()).Msg("Error message")

		configdb := responseJson["config"].(map[string]any)

		error := responseJson["error"].(string)
		error_ext := responseJson["error_ext"].(map[string]any)
		success := responseJson["success"].(float64)
		warning := responseJson["warning"].(string)

		log.Info().Str("error", error).Str("error_ext", fmt.Sprintf("%v", error_ext)).Float64("success", success).Str("warning", warning).Msg("JSON Response")

		//fmt.Printf("%v\n", configdb["databases"])
		if configdb["databases"] != nil {
			databases := configdb["databases"].([]any)
			w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
			fmt.Fprintln(w, "DATABASE\tFAMILY\tIS_PUBLISHED")
			for _, value := range databases {
				db := value.(map[string]any)
				dbname := db["database"]
				dbfamily := db["family_name"]
				is_published := db["is_published"]
				fmt.Fprintf(w, "%v\t%v\t%v\n", dbname, dbfamily, is_published)
			}
			w.Flush()
		} else {
			fmt.Printf("No database registered\n")
		}
	},
}

func init() {
	databaseCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
