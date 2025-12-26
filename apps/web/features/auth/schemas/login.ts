import { z } from "zod";

export const loginSchema = z.object({
  email: z
    .string()
    .min(1, "メールアドレスは必須です")
    .regex(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/, "有効なメールアドレスを入力してください"),
  password: z.string().min(1, "パスワードは必須です"),
});

export type LoginForm = z.infer<typeof loginSchema>;

// APIに送るペイロード
export type LoginPayload = LoginForm;
