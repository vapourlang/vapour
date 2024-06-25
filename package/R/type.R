#' Type
#'
#' Add type to the roxygen2 documentation.
#'
#' @importFrom roxygen2 roclet roxy_tag_warning block_get_tags roclet_output
#' @importFrom roxygen2 roclet_process roxy_tag_parse rd_section roxy_tag_rd
#'
#' @import roxygen2
#'
#' @export
roclet_type <- function() {
  roclet("type")
}

#' @export
roxy_tag_parse.roxy_tag_type <- function(x) {
  parts <- strsplit(x$raw, ":")[[1]]

  if(length(parts) != 2){
    roxy_tag_warning("Invalid @type tag, expects <param>: <type> | <type>")
    return()
  }

  parts <- gsub("\\n|\\t", "", parts)
  types <- strsplit(parts[2], "\\|")[[1]]

  x$val <- list(
    list(
      arg = parts[1] |> trimws(),
      types = types |> trimws()
    )
  )

  x
}

#' @export
roxy_tag_rd.roxy_tag_type <- function(x, base_path, env) {
  rd_section("type", x$val)
}

#' @export
format.rd_section_type <- function(x, ...) {
  types <- ""
  for (val in x$value) {
    t <- paste0(val$types, collapse = ", or ")
    type <- paste0("  \\item{", val$arg, "}{", t, "}\n")
    types <- paste0(types, type)
  }

  paste0(
    "\\section{Types}{\n",
    "\\itemize{\n",
    types,
    "}\n",
    "}\n"
  )
}

#' @export
roclet_process.roclet_type <- function(x, blocks, env, base_path) {
  results <- list()

  for (block in blocks) {
    tags <- block_get_tags(block, "type")
    for(tag in tags){
      t <- list(
        value = tag$val,
        type = "type",
        file = tag$file
      )
      results <- c(results, tag$val)
    }
  }

  results
}

#' @export
roclet_output.roclet_type <- function(x, results, base_path, ...) {
  .globals$types <- results
  invisible(NULL)
}
