import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: any) {
  console.info('Activating Vapour');
  const serverOptions: ServerOptions = {
      command: "vapour",
      args: ["-lsp"],
  };

  let clientOptions: LanguageClientOptions = {
      documentSelector: [{ scheme: 'file', language: 'vapour' }],
      synchronize: {
          configurationSection: 'vapour.lsp',
      },
  };

  const client: LanguageClient = new LanguageClient(
      "vp",
      "Vapour Language Server", 
      serverOptions,
      clientOptions
  );

  console.info('Initialising Vapour');

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
