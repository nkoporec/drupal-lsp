# drupal-lsp

Drupal LSP (Drupal Language Server) is a language server implementation compatible with [Language Server Protocol](https://github.com/microsoft/language-server-protocol).

## Features

- [x] Service auto-completion
- [x] Service diagnostics
- [ ] Routes auto-completion
- [ ] Route diagnostics
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
