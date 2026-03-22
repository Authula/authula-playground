const ENV_CONFIG = {
  authula: {
    url: import.meta.env.VITE_AUTHULA_URL as string,
  },
  baseUrl: import.meta.env.VITE_BASE_URL as string,
} as const;

export default ENV_CONFIG;
