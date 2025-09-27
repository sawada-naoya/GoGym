let accessToken = "";
export const tokenStore = {
  get: () => accessToken,
  set: (t: string) => {
    accessToken = t;
  },
  clear: () => {
    accessToken = "";
  },
};
