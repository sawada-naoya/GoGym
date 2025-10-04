export interface LoginResponse {
  user: { id: string | number; name: string; email: string };
  access_token: string;
  expires_in?: number;
}
