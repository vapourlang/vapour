write_config <- function() {
  config <- list(
    types = .globals$types,
    yields = .globals$yields
  )

  print(config)
}
