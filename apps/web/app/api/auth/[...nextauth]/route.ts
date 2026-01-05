// app/api/auth/[...nextauth]/route.ts
import { handlers } from "@/features/auth/auth";

export const { GET, POST } = handlers;
