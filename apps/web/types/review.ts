import { UserResponse } from "./user";

export interface ReviewResponse {
  id: number;
  content: string;
  rating: number;
  gym_id: number;
  user_id: number;
  created_at: string;
  updated_at: string;

  user?: UserResponse | null;
}

export interface ReviewListResponse {
  reviews: ReviewResponse[];
  next_cursor: string | null;
}
