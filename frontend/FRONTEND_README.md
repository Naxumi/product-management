# Product Management Frontend

A Next.js frontend application for the Product Management System with full CRUD operations.

## Features

- ✅ Product listing with pagination
- ✅ Advanced filtering (by name, SKU, category, status)
- ✅ Create new products
- ✅ Update existing products
- ✅ Delete products with confirmation
- ✅ Image upload for products
- ✅ Responsive UI with shadcn/ui components
- ✅ TypeScript for type safety

## Tech Stack

- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **shadcn/ui** - UI components
- **Tailwind CSS** - Styling
- **Lucide Icons** - Icons

## Prerequisites

- Node.js 18+ installed
- Backend API running on `http://localhost:8080`

## Getting Started

1. Install dependencies:
```bash
npm install
```

2. Create environment file:
```bash
cp .env.local.example .env.local
```

3. Update `.env.local` with your API URL if different:
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

4. Run the development server:
```bash
npm run dev
```

5. Open [http://localhost:3000](http://localhost:3000) in your browser

## Project Structure

```
frontend/
├── app/
│   ├── layout.tsx           # Root layout
│   ├── page.tsx             # Home page
│   └── products/
│       └── page.tsx         # Products CRUD page
├── components/
│   ├── ui/                  # shadcn/ui components
│   ├── product-form-dialog.tsx    # Create/Edit product dialog
│   ├── delete-confirm-dialog.tsx  # Delete confirmation dialog
│   └── image-upload.tsx           # Image upload component
├── lib/
│   ├── api.ts               # API client
│   ├── types.ts             # TypeScript types
│   └── utils.ts             # Utility functions
└── public/
    └── uploads/             # Product images (if served locally)
```

## Available Pages

- `/` - Home page with link to products
- `/products` - Product management interface

## API Endpoints Used

All endpoints are prefixed with `NEXT_PUBLIC_API_URL`:

- `GET /product` - List products with filters
- `GET /product/:id` - Get product by ID
- `GET /product/sku/:sku` - Get product by SKU
- `POST /product` - Create new product
- `PUT /product` - Update product
- `DELETE /product/:id` - Delete product
- `POST /product/:id/image` - Upload product image

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API base URL | `http://localhost:8080/api/v1` |

## Features Detail

### Product List
- Table view with product details
- Search by name or SKU
- Filter by category and status
- Pagination controls
- Actions: Edit, Delete

### Create/Edit Product
- Form with validation
- All product fields:
  - SKU (required, unique)
  - Name (required)
  - Description (optional)
  - Price (required, decimal)
  - Stock (required, integer)
  - Category (required)
  - Status (Active/Inactive)
- Image upload (only for existing products)

### Delete Product
- Confirmation dialog
- Shows product details before deletion

### Image Upload
- Drag and drop or click to upload
- Image preview
- Max file size: 5MB
- Supported formats: JPG, PNG, GIF

## Development

### Running in Development Mode
```bash
npm run dev
```

### Building for Production
```bash
npm run build
npm start
```

### Linting
```bash
npm run lint
```

## Troubleshooting

### CORS Errors
Make sure the backend has CORS enabled for `http://localhost:3000`:

```go
c.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"http://localhost:3000"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Content-Type"},
    AllowCredentials: false,
    MaxAge:           300,
}))
```

### API Connection Issues
- Check that backend is running on `http://localhost:8080`
- Verify `NEXT_PUBLIC_API_URL` in `.env.local`
- Check browser console for error messages

### Images Not Loading
- Ensure backend serves static files from `/uploads/*`
- Check that images are uploaded to correct directory
- Verify image URLs in database

## License

MIT
