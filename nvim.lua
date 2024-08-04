vim.lsp.start({
	name = "vapour",
	cmd = { "vapour -lsp" },
	root_dir = vim.fn.getcwd(),
})
