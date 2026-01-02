"use server";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export async function sendContactMessage(email: string, message: string) {
  try {
    if (!API_BASE) {
      throw new Error("API URL not configured");
    }

    const res = await fetch(`${API_BASE}/api/v1/contact`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, message }),
      cache: "no-store",
    });

    if (!res.ok) {
      throw new Error("Failed to send contact message");
    }

    return { success: true };
  } catch (error) {
    console.error("Contact action error:", error);
    return { success: false, error: "Failed to send message" };
  }
}
