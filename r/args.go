package r

import (
	"encoding/json"
	"fmt"

	"github.com/devOpifex/vapour/ast"
)

type arg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type args []arg

func GetFunctionArguments(pkg, operator, function string) *ast.FunctionLiteral {
	obj := &ast.FunctionLiteral{}
	function = pkg + operator + function

	output, err := Callr(
		fmt.Sprintf(`args <- as.list(args(%v))
    if(length(args) == 0){
		  cat("[]")
			return()
		}

    N <- length(args) - 1L
		args <- args[1:N]

    json_string <- c()
    for(i in seq_along(args)) {
			arg_string <- paste0(
				'{"name": "', names(args)[i], '", "value": "', deparse(args[i]), '"}'
			)
			json_string <- c(json_string, arg_string)
		}

		json_string <- paste0(json_string, collapse = ",")
		cat(paste0("[", json_string, "]"))`,
			function,
		),
	)

	if err != nil {
		return obj
	}

	var args args

	err = json.Unmarshal(output, &args)

	if err != nil {
		return obj
	}

	return obj
}
