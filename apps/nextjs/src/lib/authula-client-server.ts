import { cookies } from "next/headers";

import { createClient } from "authula";
import {
  CSRFPlugin,
  EmailPasswordPlugin,
  OAuth2Plugin,
  MagicLinkPlugin,
} from "authula/plugins";

import { ENV_CONFIG } from "@/constants/env-config";

export const authulaClientServer = createClient({
  url: ENV_CONFIG.authula.url,
  plugins: [
    new EmailPasswordPlugin(),
    new OAuth2Plugin(),
    new CSRFPlugin({
      cookieName: "authula_csrf_token",
      headerName: "x-authula-csrf-token",
    }),
    new MagicLinkPlugin(),
  ],
  cookies,
});
