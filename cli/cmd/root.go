/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

var kernelID, portRange, responseAddr,
	publicKey, sparkContextInitMode,
	kernelTemplateName, kernelTemplateNamespace string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubeflow-launcher",
	Short: "Launch kernels",
	Long:  `Launch kernels in the jupyter enterprise gateway`,
	Run: func(cmd *cobra.Command, args []string) {
		if kernelTemplateName == "" || kernelTemplateNamespace == "" {
			panic(fmt.Errorf("Failed to get the template's name or namespace"))
		}

		if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
			panic(err)
		}

		cfg, err := config.GetConfig()
		if err != nil {
			panic(err)
		}

		cli, err := client.New(cfg, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			panic(err)
		}

		kt := &v1alpha1.JupyterKernelTemplate{}
		if err := cli.Get(context.TODO(), client.ObjectKey{
			Namespace: kernelTemplateNamespace,
			Name:      kernelTemplateName,
		}, kt); err != nil {
			panic(err)
		}

		tpl := kt.Spec.Template

		// Set image from the kernel spec.
		image := os.Getenv("KERNEL_IMAGE")
		if image != "" && len(tpl.Template.Spec.Containers) != 0 {
			tpl.Template.Spec.Containers[0].Image = image
		}

		pod := &v1.Pod{
			ObjectMeta: tpl.ObjectMeta,
			Spec:       tpl.Template.Spec,
		}

		pod.Name = os.Getenv("KERNEL_POD_NAME")
		pod.Namespace = os.Getenv("KERNEL_POD_NAMESPACE")
		if pod.Labels == nil {
			pod.Labels = make(map[string]string)
		}
		pod.Labels["kernel_id"] = os.Getenv("KERNEL_ID")

		if err := cli.Create(context.TODO(), pod); err != nil {
			panic(err)
		}

		log.Println(kernelID, portRange,
			responseAddr, publicKey, sparkContextInitMode)
		log.Println("done")

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVar(&kernelID,
		"RemoteProcessProxy.kernel-id", "", "kernel id")
	rootCmd.Flags().StringVar(&portRange,
		"RemoteProcessProxy.port-range", "", "port range")
	rootCmd.Flags().StringVar(&responseAddr,
		"RemoteProcessProxy.response-address", "", "response address")
	rootCmd.Flags().StringVar(&publicKey,
		"RemoteProcessProxy.public-key", "", "public key")
	rootCmd.Flags().StringVar(&sparkContextInitMode,
		"RemoteProcessProxy.spark-context-initialization-mode",
		"", "spark context init mode")

	rootCmd.Flags().StringVar(&kernelTemplateName,
		"kernel-template-name", "", "kernel template CRD name")
	rootCmd.Flags().StringVar(&kernelTemplateNamespace,
		"kernel-template-namespace", "", "kernel template CRD namesapce")
}
