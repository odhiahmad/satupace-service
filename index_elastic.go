package main

import (
	"log"

	es "github.com/odhiahmad/kasirku-service/config"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/spf13/cobra"
)

func RunCreateIndexCmd(cmd *cobra.Command, args []string) {
	exists, err := es.Client.Indices.Exists([]string{"products"})
	if err == nil && exists.StatusCode == 200 {
		log.Println("ℹ️ Index sudah ada, tidak membuat ulang.")
	} else {
		err := helper.CreateElasticProductIndex()
		if err != nil {
			log.Fatalf("❌ Gagal membuat index Elasticsearch: %v", err)
		}
	}
}
