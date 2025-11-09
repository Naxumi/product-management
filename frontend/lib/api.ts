import {
  Product,
  CreateProductRequest,
  UpdateProductRequest,
  ListProductFilter,
  ListProductResponse,
  ApiResponse,
  ApiError,
  ValidationError,
} from "./types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          "Content-Type": "application/json",
          ...options?.headers,
        },
      });

      if (!response.ok) {
        const errorData: ApiError = await response.json();
        
        // Check if it's a validation error
        if (errorData.error?.code === "VALIDATION_ERROR" && errorData.error.details) {
          throw new ValidationError(
            errorData.error.message || "Validation failed",
            errorData.error.details
          );
        }
        
        // Generic error
        throw new Error(errorData.error?.message || "API request failed");
      }

      return await response.json();
    } catch (error) {
      console.error("API Error:", error);
      throw error;
    }
  }

  // Product APIs
  async getProducts(filter: Partial<ListProductFilter> = {}): Promise<ListProductResponse> {
    const params = new URLSearchParams();
    
    if (filter.name) params.append("name", filter.name);
    if (filter.sku) params.append("sku", filter.sku);
    if (filter.category) params.append("category", filter.category);
    if (filter.status) params.append("status", filter.status);
    if (filter.min_price) params.append("min_price", filter.min_price.toString());
    if (filter.max_price) params.append("max_price", filter.max_price.toString());
    if (filter.page) params.append("page", filter.page.toString());
    if (filter.limit) params.append("limit", filter.limit.toString());
    if (filter.sort_by) params.append("sort_by", filter.sort_by);
    if (filter.sort_order) params.append("sort_order", filter.sort_order);

    const queryString = params.toString();
    const endpoint = queryString ? `/product?${queryString}` : "/product";
    
    return this.request<ListProductResponse>(endpoint);
  }

  async getProduct(id: number): Promise<ApiResponse<Product>> {
    return this.request<ApiResponse<Product>>(`/product/${id}`);
  }

  async getProductBySKU(sku: string): Promise<ApiResponse<Product>> {
    return this.request<ApiResponse<Product>>(`/product/sku/${sku}`);
  }

  async createProduct(data: CreateProductRequest): Promise<ApiResponse<Product>> {
    return this.request<ApiResponse<Product>>("/product", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateProduct(data: UpdateProductRequest): Promise<ApiResponse<null>> {
    return this.request<ApiResponse<null>>("/product", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteProduct(id: number): Promise<ApiResponse<null>> {
    return this.request<ApiResponse<null>>(`/product/${id}`, {
      method: "DELETE",
    });
  }

  async uploadProductImage(id: number, file: File): Promise<ApiResponse<null>> {
    const formData = new FormData();
    formData.append("image", file);

    const url = `${this.baseURL}/product/${id}/image`;
    
    try {
      const response = await fetch(url, {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        const errorData: ApiError = await response.json();
        
        // Check if it's a validation error
        if (errorData.error?.code === "VALIDATION_ERROR" && errorData.error.details) {
          throw new ValidationError(
            errorData.error.message || "Validation failed",
            errorData.error.details
          );
        }
        
        throw new Error(errorData.error?.message || "Image upload failed");
      }

      return await response.json();
    } catch (error) {
      console.error("Upload Error:", error);
      throw error;
    }
  }

  async deleteProductImage(id: number): Promise<ApiResponse<null>> {
    return this.request<ApiResponse<null>>(`/product/${id}/image`, {
      method: "DELETE",
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
