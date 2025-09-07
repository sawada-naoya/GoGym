// API base URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'

// Request options interface
interface RequestOptions {
  query?: Record<string, any>
  headers?: Record<string, string>
  cache?: RequestCache
  next?: { revalidate?: number; tags?: string[] }
}

// Helper function to build query parameters
const buildQueryParams = (params: Record<string, any>): string => {
  const searchParams = new URLSearchParams()
  
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.append(key, value.toString())
    }
  })
  
  return searchParams.toString()
}

// Base fetch wrapper
const request = async <T = any>(
  endpoint: string,
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
  options: RequestOptions & { body?: any } = {}
): Promise<T> => {
  // Build URL with query parameters
  let url = `${API_BASE_URL}${endpoint}`
  if (options.query && Object.keys(options.query).length > 0) {
    const queryString = buildQueryParams(options.query)
    url = `${url}?${queryString}`
  }
  
  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  const config: RequestInit = {
    method,
    headers: { ...defaultHeaders, ...options.headers },
    cache: options.cache || 'no-store',
    ...options.next && { next: options.next },
  }

  // Add body for non-GET requests
  if (options.body && method !== 'GET') {
    config.body = JSON.stringify(options.body)
  }

  try {
    const response = await fetch(url, config)
    
    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`)
    }

    // Handle empty responses
    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      return await response.json()
    } else {
      return {} as T
    }
  } catch (error) {
    console.error(`API Request failed: ${method} ${url}`, error)
    throw error
  }
}

// Direct HTTP method functions
export const GET = <T = any>(endpoint: string, options?: RequestOptions): Promise<T> => 
  request<T>(endpoint, 'GET', options)

export const POST = <T = any>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> => 
  request<T>(endpoint, 'POST', { ...options, body })

export const PUT = <T = any>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> => 
  request<T>(endpoint, 'PUT', { ...options, body })

export const PATCH = <T = any>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> => 
  request<T>(endpoint, 'PATCH', { ...options, body })

export const DELETE = <T = any>(endpoint: string, options?: RequestOptions): Promise<T> => 
  request<T>(endpoint, 'DELETE', options)