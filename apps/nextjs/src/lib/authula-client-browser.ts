import { createClient } from "authula";
import {
  CorePlugin,
  EmailPasswordPlugin,
  OAuth2Plugin,
  CSRFPlugin,
  MagicLinkPlugin,
} from "authula/plugins";

import { ENV_CONFIG } from "@/constants/env-config";

export const authulaClientBrowser = createClient({
  url: ENV_CONFIG.authula.url,
  plugins: [
    new CorePlugin(),
    new EmailPasswordPlugin(),
    new OAuth2Plugin(),
    new CSRFPlugin({
      cookieName: "authula_csrf_token",
      headerName: "x-authula-csrf-token",
    }),
    new MagicLinkPlugin(),
  ],
});
