package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gzuidhof/tygo/config"
	"github.com/gzuidhof/tygo/tygo"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "tygo",
		Short: "Tool for generating Typescript from Go types",
		Long:  `Tygo generates Typescript interfaces and constants from Go files by parsing them.`,
	}

	rootCmd.PersistentFlags().
		String("config", "tygo.yaml", "config file to load (default is tygo.yaml in the current folder)")
	rootCmd.Version = FullVersion()
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode (prints debug messages)")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate and write to disk",
		Run:   generate,
	})

	cmd := &cobra.Command{
		Use:   "gendir",
		Short: "generate and write to disk no package required specify files only",
		Run:   GenDir,
	}
	cmd.Flags().BoolP("recursive", "r", false, "go inside all the dirs")
	cmd.Flags().StringP("output", "o", "types.d.ts", "the output ts file")
	rootCmd.AddCommand(cmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generate(cmd *cobra.Command, args []string) {
	cfgFilepath, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatal(err)
	}
	tygoConfig := config.ReadFromFilepath(cfgFilepath)
	t := tygo.New(&tygoConfig)

	err = t.Generate()
	if err != nil {
		log.Fatalf("Tygo failed: %v", err)
	}
}

func GenDir(cmd *cobra.Command, args []string) {
	fmt.Println("1")
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("1")
	if len(args) == 0 {
		log.Fatalf("please specify files")
	}

	fmt.Println("1")
	r, err := cmd.Flags().GetBool("recursive")
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	fmt.Println("1")
	files := []string{}
	if r {
		files, err = getFilesRec(args...)
		if err != nil {
			log.Fatalf("err: %v\n", err)
		}
	} else {
		files = args
	}

	fmt.Println("1")
	s, err := mergeGoFiles(files...)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	if s == "" {
		log.Fatalf("err: %v\n", "no data")
	}

	fmt.Println("2")
	data, err := tygo.ConvertGoToTypescript(s, tygo.PackageConfig{
		Indent:       "    ",
		TypeMappings: map[string]string{},
	})
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	err = os.WriteFile(outputPath, []byte(data), 0666)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
}
