"use client";

import { useState, useRef } from "react";
import { apiClient } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Upload, X } from "lucide-react";
import { toast } from "sonner";

interface ImageUploadProps {
  productId: number;
  currentImageUrl?: string | null;
  onSuccess: () => void;
}

export function ImageUpload({
  productId,
  currentImageUrl,
  onSuccess,
}: ImageUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [preview, setPreview] = useState<string | null>(currentImageUrl || null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file type - only allow specific image formats
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];
    const allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif'];
    const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
    
    if (!allowedTypes.includes(file.type) && !allowedExtensions.includes(fileExtension)) {
      const errorMsg = "Invalid image format. Only JPG, JPEG, PNG, and GIF are allowed";
      setError(errorMsg);
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
      setError(errorMsg);
      toast.error("Image Too Large", {
        description: errorMsg,
      });
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
      return;
    }

    setSelectedFile(file);
    setError(null);

    // Create preview
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreview(reader.result as string);
    };
    reader.readAsDataURL(file);
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    setUploading(true);
    setError(null);

    try {
      await apiClient.uploadProductImage(productId, selectedFile);
      setSelectedFile(null);
      toast.success("Image uploaded successfully");
      onSuccess();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : "Failed to upload image";
      setError(errorMsg);
      toast.error("Upload Failed", {
        description: errorMsg,
      });
      setPreview(currentImageUrl || null);
    } finally {
      setUploading(false);
    }
  };

  const handleRemovePreview = () => {
    setPreview(currentImageUrl || null);
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleDeleteImage = async () => {
    if (!currentImageUrl) return;

    setUploading(true);
    setError(null);

    try {
      await apiClient.deleteProductImage(productId);
      setPreview(null);
      setSelectedFile(null);
      toast.success("Image deleted successfully");
      onSuccess();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : "Failed to delete image";
      setError(errorMsg);
      toast.error("Delete Failed", {
        description: errorMsg,
      });
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-4">
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <div className="flex items-start gap-4">
        {preview ? (
          <div className="relative">
            <img
              src={preview}
              alt="Product preview"
              className="h-32 w-32 object-cover rounded-lg border"
            />
            <Button
              type="button"
              variant="destructive"
              size="icon"
              className="absolute -top-2 -right-2 h-6 w-6"
              onClick={handleRemovePreview}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        ) : (
          <div className="h-32 w-32 border-2 border-dashed rounded-lg flex items-center justify-center bg-gray-50 overflow-hidden">
            <div className="text-center px-3">
              <Upload className="h-7 w-7 text-gray-400 mx-auto mb-1" />
              <p className="text-[9px] text-gray-500 leading-none whitespace-nowrap">No image</p>
            </div>
          </div>
        )}

        <div className="flex-1 space-y-2">
          <input
            ref={fileInputRef}
            type="file"
            accept=".jpg,.jpeg,.png,.gif,image/jpeg,image/png,image/gif"
            onChange={handleFileSelect}
            className="hidden"
            disabled={uploading}
          />
          <div className="flex gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
            >
              Choose Image
            </Button>
            {selectedFile && (
              <Button
                type="button"
                onClick={handleUpload}
                disabled={uploading}
              >
                {uploading ? "Uploading..." : "Upload"}
              </Button>
            )}
            {currentImageUrl && !selectedFile && (
              <Button
                type="button"
                variant="destructive"
                onClick={handleDeleteImage}
                disabled={uploading}
              >
                {uploading ? "Deleting..." : "Delete Image"}
              </Button>
            )}
          </div>
          <p className="text-xs text-muted-foreground">
            Supported formats: JPG, PNG, GIF. Max size: 5MB
          </p>
        </div>
      </div>
    </div>
  );
}
