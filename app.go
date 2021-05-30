package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/andersonlira/goutils/ft"
	"github.com/andersonlira/goutils/io"
)

func main(){
	args := os.Args[1:]
	matrix := ft.Matrix{Width: 150, Separator: '-'}
	if len(args) != 2 && len(args) != 3 {
		matrix.Line()
		matrix.Println("Command must have 2 or 3 arguments. First argument local pt.json path and second argument production pt.json path.") 
		matrix.Println("Third argument if equals to 'Y' or 'y' it will ignore prompt")
		matrix.Line()
		os.Exit(1)
	}
	fileWeak := args[0]
	jsonWeak, err := io.ReadFile(fileWeak)

	if err != nil {
		panic(fmt.Sprintf("Local file was wrong: %s", fileWeak))
	}

	jsonStrong, err := io.ReadFile(args[1])

	if err != nil {
		panic(fmt.Sprintf("Production file was wrong: %s", args[1]))
	}

	result, _ := ft.MergeJson(jsonWeak,jsonStrong)

	matrix.Line()
	matrix.Println(fileWeak)
	matrix.Line()
	
	matrix.Println("")
	matrix.Println("Missing keys")
	matrix.Println("============")
	if len(result.MissingKeys) == 0 {
		matrix.PrettyP("There is not missing keys",ft.BLUE)
	}
	for k, v := range result.MissingKeys {
		matrix.PrettyP(fmt.Sprintf(`Key "%s" appear only in production with value "%s"`,k,v),ft.RED)
	}
	matrix.Line()


	matrix.Println("")
	matrix.Println("New keys")
	matrix.Println("========")
	if len(result.NewKeys) == 0 {
		matrix.PrettyP("There is not new keys",ft.BLUE)
	}
	for k, v := range result.NewKeys {
		matrix.PrettyP(fmt.Sprintf(`"%s"="%s"`,k,v),ft.GREEN)
	}
	matrix.Line()

	matrix.Println("")
	matrix.Println("DIFFS")
	matrix.Println("=====")
	for _, diff := range result.DiffKeys {
		matrix.PrettyP(diff.Key,ft.BLUE)
		matrix.PrettyP(fmt.Sprintf("   Local: %s",diff.WeakValue),ft.RED)
		matrix.PrettyP(fmt.Sprintf("   Prod: %s",diff.StrongValue),ft.GREEN)
	}
	if len(result.DiffKeys) == 0 {
		matrix.PrettyP("No differences",ft.BLUE)
	}else {
		ignorePrompt := "n"
		if len(args) == 3 {
			ignorePrompt = args[2]
		}
		if ignorePrompt == "y" || ignorePrompt == "Y" {
			mergeFile(fileWeak, result)
		}
		matrix.PrettyP(fmt.Sprintf("There are %v modifications. Apply now on local file? (Y or y)",len(result.DiffKeys)),ft.YELLOW)
		reader := bufio.NewReader(os.Stdin)
		char, _, _ := reader.ReadRune()
		switch char {
		case 'Y':
			mergeFile(fileWeak,result)
		  	break
		case 'y':
		  	mergeFile(fileWeak,result)
		  	break
		}

	}
	matrix.Line()


}

func mergeFile(file string,result ft.MergeStatus){
	io.WriteFile(file,result.FinalJSON())
}