import { z } from "zod";

export const signUpSchema = z
  .object({
    name: z.string().min(1, "名前は必須です"),
    email: z
      .string()
      .min(1, "メールアドレスは必須です")
      .regex(/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/, "有効なメールアドレスを入力してください"),
    password: z.string().min(8, "パスワードは8文字以上で入力してください").regex(/[A-Z]/, "パスワードには大文字を1文字以上含めてください").regex(/[a-z]/, "パスワードには小文字を1文字以上含めてください").regex(/[0-9]/, "パスワードには数字を1文字以上含めてください"),
    confirmPassword: z.string().min(1, "パスワード（確認）は必須です"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "パスワードが一致しません",
    path: ["confirmPassword"],
  });

export type SignUpForm = z.infer<typeof signUpSchema>;

// APIに送るペイロードは confirmPassword を含めない
export const signUpPayloadSchema = signUpSchema.pick({
  name: true,
  email: true,
  password: true,
});
export type SignUpPayload = z.infer<typeof signUpPayloadSchema>;
