vim.lsp.start({
	name = "vapour",
	cmd = { "./vapour -lsp=true" },
	root_dir = vim.fn.getcwd(),
})
