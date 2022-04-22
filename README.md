# drupal-lsp

Drupal LSP (Drupal Language Server) is a language server implementation compatible with [Language Server Protocol](https://github.com/microsoft/language-server-protocol).



![ezgif-5-924ec3e574](https://user-images.githubusercontent.com/35064680/159663609-8d8df566-d81f-4db1-8a3c-f67e4726fa7f.gif)




## Features

- [x] Service auto-completion
- [x] Service diagnostics
- [x] Service go-to definition
- [ ] Routes auto-completion
- [ ] Routes diagnostics
- [ ] Routes go-to definition
- [ ] Hooks

### Installation

To build and install the standalone drupal-lsp run

```bash
git clone https://github.com/nkoporec/drupal-lsp
go install
```

### Configuration for [neovim builtin LSP](https://neovim.io/doc/user/lsp.html) with [nvim-lspconfig](https://github.com/neovim/nvim-lspconfig)

init.vim

```vim
lua <<EOF
local lspconfig = require 'lspconfig'
local configs = require 'lspconfig.configs'
if not configs.drupal then
	configs.drupal = {
		default_config = {
    		cmd = {'drupal-lsp'},
    		filetypes = { 'php'},
			root_dir = function(fname)
			  return lspconfig.util.root_pattern('composer.json', '.git')(fname)
			end
		};
	}
   end
lspconfig.drupal.setup{autostart = true }
EOF

```

## License

MIT © [nkoporec](https://github.com/nkoporec) 
