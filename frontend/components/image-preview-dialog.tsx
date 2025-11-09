"use client";

import { X } from "lucide-react";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { VisuallyHidden } from "@radix-ui/react-visually-hidden";

interface ImagePreviewDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  imageUrl: string | null;
  productName?: string;
}

export function ImagePreviewDialog({
  open,
  onOpenChange,
  imageUrl,
  productName,
}: ImagePreviewDialogProps) {
  if (!imageUrl) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-[90vw] max-h-[90vh] p-0 bg-transparent border-0 shadow-none">
        <VisuallyHidden>
          <DialogTitle>{productName ? `${productName} - Image Preview` : "Product Image Preview"}</DialogTitle>
        </VisuallyHidden>
        <div className="relative flex items-center justify-center">
          <button
            onClick={() => onOpenChange(false)}
            className="absolute -top-12 right-0 z-10 rounded-full bg-white/90 p-2 text-gray-900 hover:bg-white transition-colors cursor-pointer shadow-lg"
          >
            <X className="h-5 w-5" />
            <span className="sr-only">Close</span>
          </button>
          <img
            src={imageUrl.startsWith('http') ? imageUrl : `http://localhost:8080/uploads/${imageUrl}`}
            alt={productName || "Product image"}
            className="w-full h-auto max-h-[85vh] object-contain rounded-lg shadow-2xl bg-white"
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}
