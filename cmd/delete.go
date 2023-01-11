/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		database := "toto"

		curlcmd := fmt.Sprintf("curl http://qserv-repl-ctl:8080/ingest/database/%v -X DELETE  -H \"Content-Type: application/json\" -d \"{\\\"auth_key\\\":\\\"\\\"}\"", database)
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

		error := responseJson["error"].(string)
		error_ext := responseJson["error_ext"].(map[string]any)
		success := responseJson["success"].(float64)
		warning := responseJson["warning"].(string)

		log.Info().Str("error", error).Str("error_ext", fmt.Sprintf("%v", error_ext)).Float64("success", success).Str("warning", warning).Msg("JSON Response")

	},
}

func init() {
	databaseCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
