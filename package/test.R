devtools::document()
devtools::load_all()
library(roxygen2)

text <- "
  #' Title
  #'
  #' @type x: integer | float
  #' @type y: numeric
  #'
  #' @yield numeric
  #'
  #' @export
  f <- function(x, y) {
    # ...
  }
"

topic <- roc_proc_text(rd_roclet(), text)[[1]]
topic$get_section("type")
