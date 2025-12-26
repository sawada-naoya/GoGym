import { z } from "zod";
import i18n from "@/lib/i18n/client";

export const signUpSchema = z
  .object({
    name: z.string().min(1, i18n.t("auth.validation.nameRequired")),
    email: z
      .string()
      .min(1, i18n.t("auth.validation.emailRequired"))
      .regex(
        /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/,
        i18n.t("auth.validation.emailInvalid"),
      ),
    password: z
      .string()
      .min(8, i18n.t("auth.validation.passwordMinLength"))
      .regex(/[A-Z]/, i18n.t("auth.validation.passwordUppercase"))
      .regex(/[a-z]/, i18n.t("auth.validation.passwordLowercase"))
      .regex(/[0-9]/, i18n.t("auth.validation.passwordNumber")),
    confirmPassword: z
      .string()
      .min(1, i18n.t("auth.validation.confirmPasswordRequired")),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: i18n.t("auth.validation.passwordMismatch"),
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
