import { ProductVariantWeightUnit } from 'shopify-api-node';

export type Product = {
  product_id: number;
  variant_id: number;
  vendor: string;
  weight: string;
  handle: string;
  weight_unit?: ProductVariantWeightUnit;
} & {
  [specification: string]: string;
}