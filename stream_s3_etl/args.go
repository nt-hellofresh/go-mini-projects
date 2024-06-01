package main

import (
	"flag"
	"log"
)

type CommandLineArgs struct {
	awsProfileName *string
	bucket         *string
	prefix         *string
}

func getArgs() CommandLineArgs {
	awsProfileName := flag.String("profile", "default", "The name of the profile in the shared credentials file")
	bucket := flag.String("bucket", "", "The name of the bucket")
	prefix := flag.String("prefix", "", "The prefix of the object")

	flag.Parse()

	return CommandLineArgs{
		awsProfileName: awsProfileName,
		bucket:         bucket,
		prefix:         prefix,
	}
}

func (args CommandLineArgs) validate() {
	if *args.awsProfileName == "" || *args.bucket == "" || *args.prefix == "" {
		log.Fatal("aws profile, bucket and prefix command line arguments are required")
	}
}
