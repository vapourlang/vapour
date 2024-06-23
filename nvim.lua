vim.lsp.start({
	name = "doctor",
	cmd = { "./doctor" },
	root_dir = vim.fn.getcwd(),
})
