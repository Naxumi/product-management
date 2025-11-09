"use client";

import { useState } from "react";
import { Product, ValidationError } from "@/lib/types";
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
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertTriangle } from "lucide-react";
import { toast } from "sonner";

interface DeleteConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  product: Product | null;
  onSuccess: () => void;
}

export function DeleteConfirmDialog({
  open,
  onOpenChange,
  product,
  onSuccess,
}: DeleteConfirmDialogProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleDelete = async () => {
    if (!product) return;

    setLoading(true);
    setError(null);

    try {
      await apiClient.deleteProduct(product.id);
      onSuccess();
      toast.success("Produk Dihapus", {
        description: `${product.name} telah dihapus`,
      });
    } catch (err) {
      let errorMessage = "Failed to delete product";
      
      if (err instanceof ValidationError) {
        const fieldErrors = Object.entries(err.details)
          .map(([field, message]) => `${field}: ${message}`)
          .join(", ");
        errorMessage = fieldErrors;
      } else if (err instanceof Error) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);
      toast.error("Gagal Menghapus Produk", {
        description: errorMessage,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-red-600">
            <AlertTriangle className="h-5 w-5" />
            Delete Product
          </DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete the product.
          </DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {product && (
          <div className="bg-gray-50 p-4 rounded-lg space-y-3 max-w-full overflow-hidden">
            <div className="space-y-1">
              <div className="text-sm font-medium text-muted-foreground">SKU:</div>
              <div className="text-sm font-mono break-all max-w-full overflow-wrap-anywhere">
                {product.sku}
              </div>
            </div>
            <div className="space-y-1">
              <div className="text-sm font-medium text-muted-foreground">Name:</div>
              <div className="text-sm break-all max-w-full overflow-wrap-anywhere">
                {product.name}
              </div>
            </div>
            <div className="space-y-1">
              <div className="text-sm font-medium text-muted-foreground">Category:</div>
              <div className="text-sm break-all max-w-full">{product.category}</div>
            </div>
          </div>
        )}

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={loading}
          >
            {loading ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
