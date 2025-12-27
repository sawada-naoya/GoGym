// app/api/auth/[...nextauth]/route.ts
import { handlers } from "@/features/auth/nextauth/auth";

export const { GET, POST } = handlers;
