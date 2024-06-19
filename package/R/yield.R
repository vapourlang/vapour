#' Yield
#'
#' Add yield to the roxygen2 documentation.
#'
#' @importFrom roxygen2 roclet roxy_tag_warning block_get_tags roclet_output
#' @importFrom roxygen2 roclet_process roxy_tag_parse rd_section roxy_tag_rd
#'
#' @import roxygen2
#' 
#' @export
roclet_yield <- function() {
  roclet("yield")
}

#' @export
roxy_tag_parse.roxy_tag_yield <- function(x) {
  raw <- gsub("\\n|\\t", "", x$raw)
  yields <- strsplit(raw, "\\|")[[1]]

  x$val <- list(
    yield = yields |> trimws()
  )

  x
}

#' @export
roxy_tag_rd.roxy_tag_yield <- function(x, base_path, env) {
  rd_section("yield", x$val)
}

#' @export
format.rd_section_yield <- function(x, ...) {
  yield <- paste0(x$value$yield, collapse = ", or ")
  paste0(
    "\\yield{", yield, "}\n"
  )
}

#' @export
roclet_process.roclet_yield <- function(x, blocks, env, base_path) {
  results <- list()
  
  for (block in blocks) {
    tags <- block_get_tags(block, "yield")
    class(tags) <- "list"
    results <- append(results, list(x))
  }
  
  results
}

#' @export
roclet_output.roclet_yield <- function(x, results, base_path, ...) {
  .globals$yields <- results
  invisible(NULL)
}
