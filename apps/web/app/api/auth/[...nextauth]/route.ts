// app/api/auth/[...nextauth]/route.ts
import { handlers } from "@/app/api/auth/[...nextauth]/authOptions";

export const { GET, POST } = handlers;
