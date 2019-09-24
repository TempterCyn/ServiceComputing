package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	flag "github.com/spf13/pflag"
)

type selpg_args struct{
	start_page int
	end_page int
	in_filename string
	page_len int
	page_type bool
	print_dest string
}



func process_args(sa *selpg_args){
	flag.IntVarP(&sa.start_page,"start","s",-1,"start page(>1)")
	flag.IntVarP(&sa.end_page,"end","e",-1,"end page(>=start_page)")
	flag.IntVarP(&sa.page_len,"len","l",72,"page len")
	flag.StringVarP(&sa.print_dest,"dest","d","","print dest")
	flag.BoolVarP(&sa.page_type, "type", "f", false, "page type")
	flag.Usage = func(){
		fmt.Fprintf(os.Stderr,"USAGE:: \n-s start_page -e end_page [ -f | -l lines_per_page ]" + " [ -d dest ][ in_filename ]\n")
	}
	flag.Parse()

	if len(os.Args)<3 {
		fmt.Fprintf(os.Stderr,"\nnot enough arguments\n")
		flag.Usage()
		os.Exit(0)
	}
	if (sa.start_page == -1) || (sa.end_page == -1) {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage and endPage can't be empty! Please check your command!\n")
		flag.Usage()
		os.Exit(0)
	} 
	if (sa.start_page <= 0) || (sa.end_page <= 0) {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage and endPage can't be negative! Please check your command!\n")
		flag.Usage()
		os.Exit(0)
	} 
	if sa.start_page > sa.end_page {
		fmt.Fprintf(os.Stderr, "\n[Error]The startPage can't be bigger than the endPage! Please check your command!\n")
		flag.Usage()
		os.Exit(0)
	}
	if len(flag.Args()) == 1{
		_, err:=os.Stat(flag.Args()[0])
		if err!=nil && os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr,"\ninput file \"%s\" does not exist\n",flag.Args()[0])
			os.Exit(0)
		}
		sa.in_filename = flag.Args()[0]
	}
	if (sa.page_type == true) && (sa.page_len != 72) {
		fmt.Fprintf(os.Stderr, "\n[Error]The command -l and -f are exclusive, you can't use them together!\n")
		flag.Usage()
		os.Exit(0)
	} 
	if sa.page_len <= 0 {
		fmt.Fprintf(os.Stderr, "\n[Error]The pageLen can't be less than 1 ! Please check your command!\n")
		flag.Usage()
		os.Exit(0)
	}
}

func process_input(args *selpg_args){
	var fin *os.File
	if args.in_filename == "" {
		fin = os.Stdin
	} else {
		var err error
		fin, err = os.Open(args.in_filename )
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n[Error]%s:", args.in_filename)
			os.Exit(0)
		}
	}
	line_count := 0
	page_count := 1
	buf := bufio.NewReader(fin)

	var cmd *exec.Cmd
	var fout io.WriteCloser
	if len(args.print_dest) > 0 {
		cmd := exec.Command("lp", "-d"+args.print_dest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fout, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n[Error]%s:", "Input pipe open\n")
			os.Exit(0)
		}
		defer fout.Close()
	}
	
	for true {
		var line string
		var err error
		if  args.page_type {
			line, err = buf.ReadString('\f')
			page_count++
		} else {
			line, err = buf.ReadString('\n')
			line_count++
			if line_count > args.page_len {
				page_count++
				line_count = 1
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n[Error]%s:", "Input pipe open\n","file read in\n")
			os.Exit(0)
		}
		if (page_count >= args.start_page) && (page_count <= args.end_page) {
			var outputErr error
			if len(args.print_dest) == 0 {
				_, outputErr = fmt.Fprintf(os.Stdout, "%s", line)
			} else {
				_, outputErr = fout.Write([]byte(line))
				if outputErr != nil {
					fmt.Fprintf(os.Stderr, "\n[Error]%s:", "pipe input")
					os.Exit(0)
				}
			}
			if outputErr != nil {
				fmt.Fprintf(os.Stderr, "\n[Error]%s:", "Error happend when output the pages.")
				os.Exit(0)
			}
		}
	}

	if len(args.print_dest) > 0 {
		fout.Close()
		errStart := cmd.Run()
		if errStart != nil {
			fmt.Fprintf(os.Stderr, "CMD Run")
			os.Exit(0)
		}
	}

	if page_count < args.start_page {
		fmt.Fprintf(os.Stderr, "\n[Error]: startPage (%d) greater than total pages (%d), no output written\n", args.start_page, page_count)
		os.Exit(0)
	} else if page_count < args.start_page {
		fmt.Fprintf(os.Stderr, "\n[Error]: endPage (%d) greater than total pages (%d), less output than expected\n",args. end_page, page_count)
		os.Exit(0)
	}
}

func main(){
	var args selpg_args
	process_args(&args)
	process_input(&args)
}

