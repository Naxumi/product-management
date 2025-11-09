"use client";

import { useState, useEffect, useCallback } from "react";
import { Product, ListProductFilter, ValidationError } from "@/lib/types";
import { apiClient } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { MoreVertical, Plus, Search, PackageOpen, Info, Edit, Trash2, X } from "lucide-react";
import { ProductFormDialog } from "@/components/product-form-dialog";
import { DeleteConfirmDialog } from "@/components/delete-confirm-dialog";
import { ImagePreviewDialog } from "@/components/image-preview-dialog";
import { toast } from "sonner";

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(0);
  const [formOpen, setFormOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [imagePreviewOpen, setImagePreviewOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [previewImageUrl, setPreviewImageUrl] = useState<string | null>(null);

  // Filter states
  const [filters, setFilters] = useState<Partial<ListProductFilter>>({});
  const [searchName, setSearchName] = useState("");
  const [searchSKU, setSearchSKU] = useState("");
  const [filterCategory, setFilterCategory] = useState<string>("");
  const [filterStatus, setFilterStatus] = useState<string>("");

  const loadProducts = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.getProducts(filters);
      setProducts(response.data.products || []);
      setTotal(response.data.total_count);
      setCurrentPage(response.data.page);
      setPageSize(response.data.limit);
      setTotalPages(response.data.total_pages);
    } catch (err) {
      let errorMessage = "Failed to load products";
      
      if (err instanceof ValidationError) {
        const fieldErrors = Object.entries(err.details)
          .map(([field, message]) => `${field}: ${message}`)
          .join(", ");
        errorMessage = fieldErrors;
      } else if (err instanceof Error) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);
      toast.error("Gagal Memuat Produk", {
        description: errorMessage,
      });
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  const handleSearch = () => {
    setFilters({
      ...filters,
      name: searchName || undefined,
      sku: searchSKU || undefined,
      category: filterCategory || undefined,
      status: filterStatus && filterStatus !== "all" ? (filterStatus as "Active" | "Inactive") : undefined,
      page: 1,
    });
    setCurrentPage(1);
  };

  const handleResetFilters = () => {
    setSearchName("");
    setSearchSKU("");
    setFilterCategory("");
    setFilterStatus("");
    setFilters({});
    setCurrentPage(1);
  };

  const handleCreate = () => {
    setSelectedProduct(null);
    setFormOpen(true);
  };

  const handleEdit = (product: Product) => {
    setSelectedProduct(product);
    setFormOpen(true);
  };

  const handleDelete = (product: Product) => {
    setSelectedProduct(product);
    setDeleteOpen(true);
  };

  const handleImageClick = (imageUrl: string | null | undefined, productName: string) => {
    if (imageUrl) {
      setPreviewImageUrl(imageUrl);
      setImagePreviewOpen(true);
      // Don't set selectedProduct here, just keep the current state
    }
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setSelectedProduct(null);
    loadProducts();
    toast.success(selectedProduct ? "Product updated successfully" : "Product created successfully");
  };

  const handleDeleteSuccess = () => {
    setDeleteOpen(false);
    setSelectedProduct(null);
    loadProducts();
    toast.success("Product deleted successfully");
  };

  const handlePageChange = (page: number) => {
    setFilters({ ...filters, page });
  };

  if (loading && products.length === 0) {
    return (
      <div className="container mx-auto py-10">
        <div className="flex items-center justify-center h-64">
          <p className="text-muted-foreground">Loading products...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6 px-4 max-w-7xl">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-2xl font-bold">Products</CardTitle>
            <Button onClick={handleCreate}>
              <Plus className="mr-2 h-4 w-4" />
              Add Product
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {error && (
            <Alert variant="destructive" className="flex justify-between items-start">
              <div className="flex-1">
                <AlertDescription>{error}</AlertDescription>
              </div>
              <button 
                className="cursor-pointer ml-2 mt-0.5" 
                onClick={() => setError(null)}
              >
                <X className="h-4 w-4" />
                <span className="sr-only">Close</span>
              </button>
            </Alert>
          )}

          {/* Filters */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-3">
            <Input
              placeholder="Search by name..."
              value={searchName}
              onChange={(e) => setSearchName(e.target.value)}
            />
            <Input
              placeholder="Search by SKU..."
              value={searchSKU}
              onChange={(e) => setSearchSKU(e.target.value)}
            />
            <Input
              placeholder="Filter by category..."
              value={filterCategory}
              onChange={(e) => setFilterCategory(e.target.value)}
            />
            <Select value={filterStatus} onValueChange={setFilterStatus}>
              <SelectTrigger>
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="Active">Active</SelectItem>
                <SelectItem value="Inactive">Inactive</SelectItem>
              </SelectContent>
            </Select>
            <div className="flex gap-2">
              <Button onClick={handleSearch} className="flex-1">
                <Search className="mr-2 h-4 w-4" />
                Search
              </Button>
              <Button onClick={handleResetFilters} variant="outline">
                Reset
              </Button>
            </div>
          </div>

          {/* Table */}
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-16">Image</TableHead>
                  <TableHead className="w-32">SKU</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead className="max-w-xs">Description</TableHead>
                  <TableHead className="w-28">Category</TableHead>
                  <TableHead className="w-32">Price</TableHead>
                  <TableHead className="w-20">Stock</TableHead>
                  <TableHead className="w-24">Status</TableHead>
                  <TableHead className="w-20 text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {products.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={9} className="h-64">
                    <div className="flex flex-col items-center justify-center text-center">
                      <PackageOpen className="h-12 w-12 text-muted-foreground mb-4" />
                      <h3 className="text-lg font-semibold mb-2">No products found</h3>
                      <p className="text-sm text-muted-foreground mb-4">
                        {Object.keys(filters).length > 0
                          ? "Try adjusting your search or filter to find what you're looking for."
                          : "Get started by creating your first product."}
                      </p>
                      {Object.keys(filters).length > 0 ? (
                        <Button onClick={handleResetFilters} variant="outline">
                          Clear filters
                        </Button>
                      ) : (
                        <Button onClick={handleCreate}>
                          <Plus className="mr-2 h-4 w-4" />
                          Add Product
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                  products.map((product) => (
                    <TableRow key={product.id}>
                      <TableCell className="py-2">
                        {product.image_url ? (
                          <img
                            src={product.image_url.startsWith('http') ? product.image_url : `http://localhost:8080/uploads/${product.image_url}`}
                            alt={product.name}
                            className="h-13 w-13 object-cover rounded cursor-pointer hover:opacity-80 transition-opacity"
                            onClick={() => handleImageClick(product.image_url, product.name)}
                          />
                        ) : (
                          <div className="h-13 w-13 bg-gray-200 rounded flex items-center justify-center text-xs text-gray-500">
                            No image
                          </div>
                        )}
                      </TableCell>
                      <TableCell className="font-mono text-xs py-2">
                        {product.sku.length > 50 ? `${product.sku.substring(0, 50)}...` : product.sku}
                      </TableCell>
                      <TableCell className="font-medium py-2">
                        {product.name.length > 50 ? `${product.name.substring(0, 50)}...` : product.name}
                      </TableCell>
                      <TableCell className="py-2 text-sm text-muted-foreground">
                        {product.description 
                          ? (product.description.length > 50 ? `${product.description.substring(0, 50)}...` : product.description)
                          : "-"}
                      </TableCell>
                      <TableCell className="py-2">
                        {product.category.length > 50 ? `${product.category.substring(0, 50)}...` : product.category}
                      </TableCell>
                      <TableCell className="py-2">Rp {Number(product.price).toLocaleString()}</TableCell>
                      <TableCell className="py-2">{product.stock}</TableCell>
                      <TableCell className="py-2">
                        <Badge variant={product.status === "Active" ? "default" : "secondary"}>
                          {product.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right py-2">
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm">
                              <MoreVertical className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() => {
                                setDetailOpen(true);
                                setSelectedProduct(product);
                              }}
                            >
                              <Info className="mr-2 h-4 w-4" />
                              Detail
                            </DropdownMenuItem>
                            <DropdownMenuItem onClick={() => handleEdit(product)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => handleDelete(product)}
                              className="text-red-600"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {total > 0 && (
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <p className="text-sm text-muted-foreground">
                  Showing {(currentPage - 1) * pageSize + 1} to{" "}
                  {Math.min(currentPage * pageSize, total)} of {total} products
                </p>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">Items per page:</span>
                  <Select
                    value={pageSize.toString()}
                    onValueChange={(value) => {
                      setFilters({ ...filters, limit: parseInt(value), page: 1 });
                    }}
                  >
                    <SelectTrigger className="w-[70px] h-8">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="5">5</SelectItem>
                      <SelectItem value="10">10</SelectItem>
                      <SelectItem value="20">20</SelectItem>
                      <SelectItem value="50">50</SelectItem>
                      <SelectItem value="100">100</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              {totalPages > 1 && (
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                <div className="flex items-center gap-1">
                  {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                    let pageNum;
                    if (totalPages <= 5) {
                      pageNum = i + 1;
                    } else if (currentPage <= 3) {
                      pageNum = i + 1;
                    } else if (currentPage >= totalPages - 2) {
                      pageNum = totalPages - 4 + i;
                    } else {
                      pageNum = currentPage - 2 + i;
                    }
                    return (
                      <Button
                        key={pageNum}
                        variant={currentPage === pageNum ? "default" : "outline"}
                        size="sm"
                        onClick={() => handlePageChange(pageNum)}
                      >
                        {pageNum}
                      </Button>
                    );
                  })}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      <ProductFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        product={selectedProduct}
        onSuccess={handleFormSuccess}
      />

      <DeleteConfirmDialog
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
        product={selectedProduct}
        onSuccess={handleDeleteSuccess}
      />

      <Sheet open={detailOpen} onOpenChange={setDetailOpen}>
        <SheetContent className="w-full sm:max-w-md overflow-y-auto overflow-x-hidden">
          <SheetHeader>
            <SheetTitle>Product Details</SheetTitle>
            <SheetDescription>
              Complete information about this product
            </SheetDescription>
          </SheetHeader>
          {selectedProduct && (
            <div className="mt-6 space-y-6 max-w-full overflow-hidden">
              {/* Image */}
              {selectedProduct.image_url && (
                <div className="space-y-2">
                  <h4 className="text-sm font-medium text-muted-foreground">Image</h4>
                  <img
                    src={selectedProduct.image_url.startsWith('http') ? selectedProduct.image_url : `http://localhost:8080/uploads/${selectedProduct.image_url}`}
                    alt={selectedProduct.name}
                    className="w-full h-48 object-cover rounded-lg border cursor-pointer hover:opacity-80 transition-opacity"
                    onClick={() => handleImageClick(selectedProduct.image_url, selectedProduct.name)}
                  />
                  <p className="text-xs text-muted-foreground text-center">Click image to view larger</p>
                </div>
              )}
              
              {/* ID */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">ID</h4>
                <p className="text-base break-all">{selectedProduct.id}</p>
              </div>

              {/* SKU */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">SKU</h4>
                <p className="text-base font-mono break-all max-w-full overflow-wrap-anywhere">
                  {selectedProduct.sku}
                </p>
              </div>

              {/* Name */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Name</h4>
                <p className="text-base font-medium break-all max-w-full overflow-wrap-anywhere">
                  {selectedProduct.name}
                </p>
              </div>

              {/* Description */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Description</h4>
                <p className="text-base text-justify wrap-break-word max-w-full overflow-wrap-anywhere">
                  {selectedProduct.description || "-"}
                </p>
              </div>

              {/* Category */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Category</h4>
                <p className="text-base break-all max-w-full">{selectedProduct.category}</p>
              </div>

              {/* Price */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Price</h4>
                <p className="text-base font-semibold">
                  Rp {Number(selectedProduct.price).toLocaleString('id-ID')}
                </p>
              </div>

              {/* Stock */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Stock</h4>
                <p className="text-base">{selectedProduct.stock} units</p>
              </div>

              {/* Status */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Status</h4>
                <Badge variant={selectedProduct.status === "Active" ? "default" : "secondary"}>
                  {selectedProduct.status}
                </Badge>
              </div>

              {/* Created At */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Created At</h4>
                <p className="text-base">
                  {new Date(selectedProduct.created_at).toLocaleString('id-ID', {
                    dateStyle: 'long',
                    timeStyle: 'short'
                  })}
                </p>
              </div>

              {/* Updated At */}
              <div className="space-y-1">
                <h4 className="text-sm font-medium text-muted-foreground">Updated At</h4>
                <p className="text-base">
                  {new Date(selectedProduct.updated_at).toLocaleString('id-ID', {
                    dateStyle: 'long',
                    timeStyle: 'short'
                  })}
                </p>
              </div>

              {/* Actions */}
              <div className="flex gap-2 pt-4 border-t">
                <Button onClick={() => {
                  setDetailOpen(false);
                  handleEdit(selectedProduct);
                }} className="flex-1">
                  Edit Product
                </Button>
                <Button 
                  onClick={() => {
                    setDetailOpen(false);
                    handleDelete(selectedProduct);
                  }} 
                  variant="destructive"
                  className="flex-1"
                >
                  Delete Product
                </Button>
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>

      <ImagePreviewDialog
        open={imagePreviewOpen}
        onOpenChange={setImagePreviewOpen}
        imageUrl={previewImageUrl}
        productName={selectedProduct?.name}
      />
    </div>
  );
}
