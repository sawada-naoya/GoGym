import { z } from "zod";
import i18n from "@/lib/i18n/client";

export const loginSchema = z.object({
  email: z
    .string()
    .min(1, i18n.t("auth.validation.emailRequired"))
    .regex(
      /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
      i18n.t("auth.validation.emailInvalid"),
    ),
  password: z.string().min(1, i18n.t("auth.validation.passwordRequired")),
});

export type LoginForm = z.infer<typeof loginSchema>;

// APIに送るペイロード
export type LoginPayload = LoginForm;
