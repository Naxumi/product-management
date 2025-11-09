"use client";

import { useState, useEffect, useRef } from "react";
import { Product, CreateProductRequest, UpdateProductRequest, ValidationError } from "@/lib/types";
import { apiClient } from "@/lib/api";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { toast } from "sonner";

interface ProductFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: Product | null;
  onSuccess: () => void;
}

export function ProductFormDialog({
  open,
  onOpenChange,
  product,
  onSuccess,
}: ProductFormDialogProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [imageDeleted, setImageDeleted] = useState(false); // Track if user deleted the image
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [formData, setFormData] = useState({
    sku: "",
    name: "",
    description: "",
    price: "",
    stock: "",
    category: "",
    status: "Active" as "Active" | "Inactive",
  });

  // Format number to IDR
  const formatToIDR = (value: string): string => {
    // Remove non-digit characters
    const numbers = value.replace(/\D/g, '');
    if (!numbers) return '';
    
    // Format with thousand separator
    const formatted = parseInt(numbers).toLocaleString('id-ID');
    return `Rp ${formatted}`;
  };

  // Parse IDR formatted string to number
  const parseIDR = (value: string): string => {
    return value.replace(/\D/g, '');
  };

  useEffect(() => {
    if (product) {
      setFormData({
        sku: product.sku,
        name: product.name,
        description: product.description || "",
        price: product.price.toString(),
        stock: product.stock.toString(),
        category: product.category,
        status: product.status,
      });
      setImagePreview(product.image_url || null);
      setImageDeleted(false); // Reset delete flag
    } else {
      setFormData({
        sku: "",
        name: "",
        description: "",
        price: "",
        stock: "",
        category: "",
        status: "Active",
      });
      setSelectedImage(null);
      setImagePreview(null);
      setImageDeleted(false);
    }
    setError(null);
  }, [product, open]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // Validate price limit before submission
    const priceValue = parseFloat(parseIDR(formData.price));
    const maxPrice = 99999999.99;
    
    if (priceValue > maxPrice) {
      toast.error("Harga Terlalu Tinggi", {
        description: `Harga maksimal adalah Rp ${maxPrice.toLocaleString('id-ID')}`,
      });
      setLoading(false);
      return;
    }

    // Validate image format before submission if there's a selected image
    if (selectedImage) {
      const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];
      const allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif'];
      const fileExtension = selectedImage.name.toLowerCase().substring(selectedImage.name.lastIndexOf('.'));
      
      const isValidType = allowedTypes.includes(selectedImage.type);
      const isValidExtension = allowedExtensions.includes(fileExtension);
      
      if (!isValidType || !isValidExtension) {
        const errorMsg = "Invalid image format. Only JPG, JPEG, PNG, and GIF are allowed";
        toast.error("Invalid Image Format", {
          description: errorMsg,
        });
        setLoading(false);
        return;
      }

      // Validate file size (max 5MB)
      const maxSize = 5 * 1024 * 1024; // 5MB in bytes
      if (selectedImage.size > maxSize) {
        const errorMsg = "Image file size exceeds maximum limit of 5MB";
        toast.error("Image Too Large", {
          description: errorMsg,
        });
        setLoading(false);
        return;
      }
    }

    try {
      if (product) {
        // Update existing product
        const updateData: UpdateProductRequest = {
          id: product.id,
        };

        // Only include changed fields
        if (formData.sku !== product.sku) updateData.sku = formData.sku;
        if (formData.name !== product.name) updateData.name = formData.name;
        if (formData.description !== (product.description || ""))
          updateData.description = formData.description;
        if (formData.price !== product.price.toString())
          updateData.price = parseFloat(parseIDR(formData.price));
        if (formData.stock !== product.stock.toString())
          updateData.stock = parseInt(formData.stock);
        if (formData.category !== product.category)
          updateData.category = formData.category;
        if (formData.status !== product.status)
          updateData.status = formData.status;

        await apiClient.updateProduct(updateData);

        // Handle image operations
        if (imageDeleted && product.image_url) {
          // User deleted the existing image
          await apiClient.deleteProductImage(product.id);
        } else if (selectedImage) {
          // User uploaded a new image
          await apiClient.uploadProductImage(product.id, selectedImage);
        }
      } else {
        // Create new product
        const createData: CreateProductRequest = {
          sku: formData.sku,
          name: formData.name,
          description: formData.description || undefined,
          price: parseFloat(parseIDR(formData.price)),
          stock: parseInt(formData.stock),
          category: formData.category,
          status: formData.status,
        };

        const response = await apiClient.createProduct(createData);
        
        // Upload image if selected and product created successfully
        if (selectedImage && response.data.id) {
          await apiClient.uploadProductImage(response.data.id, selectedImage);
        }
      }

      onSuccess();
    } catch (err) {
      // Handle validation errors with detailed field messages
      if (err instanceof ValidationError) {
        const fieldErrors = Object.entries(err.details)
          .map(([field, message]) => `${field}: ${message}`)
          .join(", ");
        
        toast.error("Validasi Gagal", {
          description: fieldErrors,
        });
        setError(err.message);
      } else if (err instanceof Error) {
        toast.error("Error", {
          description: err.message,
        });
        setError(err.message);
      } else {
        toast.error("Error", {
          description: "Failed to save product",
        });
        setError("Failed to save product");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handlePriceChange = (value: string) => {
    // Only allow digits and format to IDR
    const numbers = value.replace(/\D/g, '');
    
    // Check max price limit (NUMERIC(10,2) allows max 99,999,999.99)
    const maxPrice = 99999999.99;
    const numericValue = parseFloat(numbers) || 0;
    
    if (numericValue > maxPrice) {
      toast.error("Harga Terlalu Tinggi", {
        description: `Harga maksimal adalah Rp ${maxPrice.toLocaleString('id-ID')}`,
      });
      return;
    }
    
    setFormData((prev) => ({ ...prev, price: numbers }));
  };

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type - only allow specific image formats
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];
    const allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif'];
    const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
    
    const isValidType = allowedTypes.includes(file.type);
    const isValidExtension = allowedExtensions.includes(fileExtension);
    
    if (!isValidType || !isValidExtension) {
      const errorMsg = "Invalid image format. Only JPG, JPEG, PNG, and GIF are allowed";
      toast.error("Invalid Image Format", {
        description: errorMsg,
      });
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
      return;
    }

    // Validate file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB in bytes
    if (file.size > maxSize) {
      const errorMsg = "Image file size exceeds maximum limit of 5MB";
      toast.error("Image Too Large", {
        description: errorMsg,
      });
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
      return;
    }

    setSelectedImage(file);
    setError(null);

    // Create preview
    const reader = new FileReader();
    reader.onloadend = () => {
      setImagePreview(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  const handleRemoveImage = () => {
    setSelectedImage(null);
    setImagePreview(null); // Always set to null to show "No image"
    setImageDeleted(true); // Mark that user deleted the image
    // Reset file input
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{product ? "Edit Product" : "Add New Product"}</DialogTitle>
          <DialogDescription>
            {product
              ? "Update the product details below."
              : "Fill in the details to create a new product."}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="sku">SKU *</Label>
              <Input
                id="sku"
                value={formData.sku}
                onChange={(e) => handleChange("sku", e.target.value)}
                required
                placeholder="e.g., PROD-001"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input
                id="name"
                value={formData.name}
                onChange={(e) => handleChange("name", e.target.value)}
                required
                placeholder="Product name"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              value={formData.description}
              onChange={(e) => handleChange("description", e.target.value)}
              placeholder="Product description (optional)"
              rows={3}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="price">Price *</Label>
              <Input
                id="price"
                type="text"
                value={formData.price ? formatToIDR(formData.price) : ''}
                onChange={(e) => handlePriceChange(e.target.value)}
                required
                placeholder="Rp 0"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="stock">Stock *</Label>
              <Input
                id="stock"
                type="number"
                min="0"
                value={formData.stock}
                onChange={(e) => handleChange("stock", e.target.value)}
                required
                placeholder="0"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="category">Category *</Label>
              <Input
                id="category"
                value={formData.category}
                onChange={(e) => handleChange("category", e.target.value)}
                required
                placeholder="e.g., Electronics"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="status">Status *</Label>
              <Select
                value={formData.status}
                onValueChange={(value) =>
                  handleChange("status", value as "Active" | "Inactive")
                }
              >
                <SelectTrigger id="status">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Active">Active</SelectItem>
                  <SelectItem value="Inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="image">Product Image</Label>
            <div className="space-y-4">
              <div className="flex items-start gap-4">
                {imagePreview ? (
                  <div className="relative">
                    <img
                      src={imagePreview}
                      alt="Product preview"
                      className="h-32 w-32 object-cover rounded-lg border"
                    />
                    <Button
                      type="button"
                      variant="destructive"
                      size="icon"
                      className="absolute -top-2 -right-2 h-6 w-6"
                      onClick={handleRemoveImage}
                    >
                      ×
                    </Button>
                  </div>
                ) : (
                  <div className="h-32 w-32 border-2 border-dashed rounded-lg flex items-center justify-center bg-gray-50 overflow-hidden">
                    <div className="text-center px-3">
                      <div className="text-3xl mb-1">📷</div>
                      <p className="text-[9px] text-gray-500 leading-none whitespace-nowrap">No image</p>
                    </div>
                  </div>
                )}

                <div className="flex-1 space-y-2">
                  <Input
                    ref={fileInputRef}
                    id="image"
                    type="file"
                    accept=".jpg,.jpeg,.png,.gif,image/jpeg,image/png,image/gif"
                    onChange={handleImageSelect}
                    className="cursor-pointer"
                  />
                  <p className="text-xs text-muted-foreground">
                    {product ? "Upload a new image to replace the current one." : "Optional."} Supported formats: JPG, PNG, GIF. Max size: 5MB
                  </p>
                </div>
              </div>
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? "Saving..." : product ? "Update" : "Create"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
