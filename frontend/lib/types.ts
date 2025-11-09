export type ProductStatus = "Active" | "Inactive";

export interface Product {
  id: number;
  sku: string;
  name: string;
  description?: string | null;
  price: number;
  stock: number;
  category: string;
  status: ProductStatus;
  image_url?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateProductRequest {
  sku: string;
  name: string;
  description?: string;
  price: number;
  stock: number;
  category: string;
  status: ProductStatus;
}

export interface UpdateProductRequest {
  id: number;
  sku?: string;
  name?: string;
  description?: string;
  price?: number;
  stock?: number;
  category?: string;
  status?: ProductStatus;
}

export interface ListProductFilter {
  name?: string;
  sku?: string;
  category?: string;
  status?: ProductStatus;
  min_price?: number;
  max_price?: number;
  page: number;
  limit: number;
  sort_by?: string;
  sort_order?: string;
}

export interface ListProductResponse {
  success: boolean;
  data: {
    products: Product[];
    total_count: number;
    page: number;
    limit: number;
    total_pages: number;
  };
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface ApiError {
  success: false;
  error: {
    code?: string;
    message: string;
    details?: Record<string, string>;
  };
}

export class ValidationError extends Error {
  public details: Record<string, string>;
  
  constructor(message: string, details: Record<string, string>) {
    super(message);
    this.name = 'ValidationError';
    this.details = details;
  }
}
