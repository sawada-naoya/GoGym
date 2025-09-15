export interface Location {
  latitude: number;
  longitude: number;
}

export interface Tag {
  id: number;
  name: string;
}
export interface Gym {
  id: number;
  name: string;
  description: string | null;
  location: Location;
  address: string;
  city: string | null;
  prefecture: string | null;
  postal_code: string | null;
  is_active: boolean;
  tags: Tag[];
  average_rating: number | null;
  review_count: number;
  created_at: string;
  updated_at: string;
}

export interface SearchGymResponse {
  gyms?: Gym[];
  next_cursor: string | null;
  has_more: boolean;
}

export interface SearchGymParams {
  q?: string;
  lat?: number;
  lon?: number;
  radius_m?: number;
  cursor?: string;
  limit?: number;
}
