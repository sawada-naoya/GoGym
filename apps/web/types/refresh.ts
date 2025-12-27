export type RefreshResponse = {
  user: { id: string; name: string; email: string };
  access_token: string;
  refresh_token: string;
  expires_in: number;
};
