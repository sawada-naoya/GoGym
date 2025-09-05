#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Generating TypeScript client from OpenAPI spec...${NC}"

# Check if openapi-typescript is installed
if ! command -v openapi-typescript &> /dev/null; then
    echo -e "${RED}❌ openapi-typescript is not installed${NC}"
    echo "Install it with: npm install -g openapi-typescript"
    exit 1
fi

# Create output directory
OUTPUT_DIR="../../apps/web/types/generated"
mkdir -p "$OUTPUT_DIR"

# Generate TypeScript types
echo -e "${YELLOW}📋 Generating TypeScript types...${NC}"
openapi-typescript openapi.yaml -o "$OUTPUT_DIR/api.ts"

# Create client helper
echo -e "${YELLOW}🔌 Creating API client helper...${NC}"
cat > "$OUTPUT_DIR/client.ts" << 'EOF'
/**
 * Generated API Client for GoGym
 * 
 * This file provides a typed client for the GoGym API.
 * Generated from OpenAPI specification.
 */

import type { paths } from './api';

// Extract operation types from the paths
export type SearchGymsParams = paths['/gyms/search']['get']['parameters']['query'];
export type SearchGymsResponse = paths['/gyms/search']['get']['responses'][200]['content']['application/json'];

export type GetGymParams = paths['/gyms/{id}']['get']['parameters']['path'];
export type GetGymResponse = paths['/gyms/{id}']['get']['responses'][200]['content']['application/json'];

export type SignupRequest = paths['/auth/signup']['post']['requestBody']['content']['application/json'];
export type SignupResponse = paths['/auth/signup']['post']['responses'][201]['content']['application/json'];

export type LoginRequest = paths['/auth/login']['post']['requestBody']['content']['application/json'];
export type LoginResponse = paths['/auth/login']['post']['responses'][200]['content']['application/json'];

export type PresignRequest = paths['/photos/presign']['post']['requestBody']['content']['application/json'];
export type PresignResponse = paths['/photos/presign']['post']['responses'][200]['content']['application/json'];

export type ErrorResponse = paths['/gyms/search']['get']['responses'][400]['content']['application/json'];

// API Client Configuration
export interface ApiConfig {
  baseURL: string;
  accessToken?: string;
}

// Default configuration
export const defaultApiConfig: ApiConfig = {
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
};

// API Client Class
export class ApiClient {
  private config: ApiConfig;

  constructor(config: Partial<ApiConfig> = {}) {
    this.config = { ...defaultApiConfig, ...config };
  }

  // Update access token
  setAccessToken(token: string) {
    this.config.accessToken = token;
  }

  // Remove access token
  clearAccessToken() {
    delete this.config.accessToken;
  }

  // Get headers for requests
  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.config.accessToken) {
      headers.Authorization = `Bearer ${this.config.accessToken}`;
    }

    return headers;
  }

  // Generic request method
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.config.baseURL}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.getHeaders(),
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // API Methods
  async searchGyms(params: SearchGymsParams = {}): Promise<SearchGymsResponse> {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value));
      }
    });

    const query = searchParams.toString();
    return this.request(`/gyms/search${query ? `?${query}` : ''}`);
  }

  async getGym(id: number): Promise<GetGymResponse> {
    return this.request(`/gyms/${id}`);
  }

  async signup(data: SignupRequest): Promise<SignupResponse> {
    return this.request('/auth/signup', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async presignPhoto(data: PresignRequest): Promise<PresignResponse> {
    return this.request('/photos/presign', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}

// Export default instance
export const apiClient = new ApiClient();
EOF

# Create index file for easy imports
echo -e "${YELLOW}📄 Creating index file...${NC}"
cat > "$OUTPUT_DIR/index.ts" << 'EOF'
// Export all types and client
export type * from './api';
export * from './client';
EOF

# Add header to the generated api.ts file
echo -e "${YELLOW}📝 Adding header to generated types...${NC}"
sed -i '1i\
/**\
 * Generated TypeScript definitions for GoGym API\
 * \
 * This file contains TypeScript definitions generated from the OpenAPI specification.\
 * DO NOT EDIT - This file is auto-generated.\
 */\
' "$OUTPUT_DIR/api.ts"

echo -e "${GREEN}✅ TypeScript client generation completed!${NC}"
echo -e "${GREEN}📁 Generated files in: $OUTPUT_DIR${NC}"
echo ""
echo "Generated files:"
ls -la "$OUTPUT_DIR"
echo ""
echo -e "${GREEN}💡 Usage example:${NC}"
echo "import { apiClient, SearchGymsParams } from '@/types/generated';"
echo ""
echo "const params: SearchGymsParams = { lat: 35.6762, lon: 139.6503, radius_m: 5000 };"
echo "const result = await apiClient.searchGyms(params);"